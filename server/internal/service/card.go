package service

import (
	"context"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"slices"
	"sort"
	"strings"
	pb "vgpu/api/v1"
	"vgpu/internal/biz"
	"vgpu/internal/database"
)

type CardService struct {
	pb.UnimplementedCardServer

	node *biz.NodeUsecase
	pod  *biz.PodUseCase
	ms   *MonitorService
}

func NewCardService(node *biz.NodeUsecase, pod *biz.PodUseCase, ms *MonitorService) *CardService {
	return &CardService{node: node, pod: pod, ms: ms}
}

func (s *CardService) GetAllGPUs(ctx context.Context, req *pb.GetAllGpusReq) (*pb.GPUsReply, error) {
	filters := req.Filters
	deviceInfos, err := s.node.ListAllDevices(ctx)
	if err != nil {
		return nil, err
	}
	var res = &pb.GPUsReply{List: []*pb.GPUReply{}}
	for _, device := range deviceInfos {
		gpu := &pb.GPUReply{}
		//nodeName := strings.Trim(filters.NodeName, " ")
		//if nodeName != "" && nodeName != device.NodeName {
		//	continue
		//}
		//deviceType := strings.Trim(filters.Type, " ")
		//if deviceType != "" && deviceType != device.Type {
		//	continue
		//}
		//deviceUid := strings.Trim(filters.Uid, " ")
		//if deviceUid != "" && deviceUid != device.Id {
		//	continue
		//}

		nodeNames := strings.Trim(filters.NodeName, " ")
		if nodeNames != "" {
			names := strings.Split(nodeNames, "|")
			log.Info("GetAllGPUs names: ", names)
			if !slices.Contains(names, device.NodeName) {
				continue
			}
		}

		deviceTypes := strings.Trim(filters.Type, " ")
		if deviceTypes != "" {
			types := strings.Split(deviceTypes, "|")
			log.Info("GetAllGPUs types: ", types)
			if !slices.Contains(types, device.Type) {
				continue
			}
		}
		deviceUids := strings.Trim(filters.Uid, " ")
		if deviceUids != "" {
			uids := strings.Split(deviceUids, "|")
			log.Info("GetAllGPUs uids: ", uids)
			if !slices.Contains(uids, device.NodeUid) {
				continue
			}
		}
		gpu.Uuid = device.Id
		gpu.NodeName = device.NodeName
		gpu.Type = device.Type
		gpu.VgpuTotal = device.Count
		gpu.CoreTotal = device.Devcore
		gpu.MemoryTotal = device.Devmem
		gpu.NodeUid = device.NodeUid
		gpu.Health = device.Health
		gpu.Mode = device.Mode
		resourcePoolNames, err := database.QueryResourceNamesByNodeName(device.NodeName)
		if err != nil {
			return nil, err
		}
		gpu.ResourcePools = resourcePoolNames

		vGPU, core, memory, err := s.pod.StatisticsByDeviceId(ctx, device.AliasId)
		if err == nil {
			gpu.VgpuUsed = vGPU
			gpu.CoreUsed = core
			gpu.MemoryUsed = memory
		}
		resp, err := s.ms.QueryInstant(ctx, &pb.QueryInstantRequest{Query: fmt.Sprintf("avg(hami_core_size{deviceuuid=~\"%s\"})", device.Id)})
		if err == nil && len(resp.Data) > 0 {
			gpu.CoreTotal = int32(resp.Data[0].Value)
		}
		resp, err = s.ms.QueryInstant(ctx, &pb.QueryInstantRequest{Query: fmt.Sprintf("avg(hami_memory_size{deviceuuid=~\"%s\"})", device.Id)})
		if err == nil && len(resp.Data) > 0 {
			gpu.MemoryTotal = int32(resp.Data[0].Value)
		}
		res.List = append(res.List, gpu)
	}

	sort.SliceStable(res.List, func(i, j int) bool {
		return res.List[i].Uuid < res.List[j].Uuid
	})
	return res, nil
}

func (s *CardService) GetAllGPUTypes(ctx context.Context, req *pb.GetAllGpusReq) (*pb.GPUsReply, error) {
	deviceInfos, err := s.node.ListAllDevices(ctx)
	if err != nil {
		return nil, err
	}

	var res = &pb.GPUsReply{List: []*pb.GPUReply{}}
	seenTypes := make(map[string]struct{})

	filters := req.Filters
	provider := strings.Trim(filters.Provider, " ")
	for _, device := range deviceInfos {
		if provider != "" && provider != device.Provider {
			continue
		}

		if _, exists := seenTypes[device.Type]; !exists {
			seenTypes[device.Type] = struct{}{}
			gpu := &pb.GPUReply{}
			gpu.Type = device.Type
			res.List = append(res.List, gpu)
		}
	}

	return res, nil
}

func (s *CardService) GetGPU(ctx context.Context, req *pb.GetGpuReq) (*pb.GPUReply, error) {
	devices, err := s.node.ListAllDevices(ctx)
	if err != nil {
		return nil, err
	}
	gpu := &pb.GPUReply{}
	for _, device := range devices {
		deviceUid := strings.Trim(req.Uid, " ")
		if deviceUid == "" || deviceUid != device.Id {
			continue
		}
		gpu.Uuid = device.Id
		gpu.NodeName = device.NodeName
		gpu.Type = device.Type
		gpu.VgpuTotal = device.Count
		gpu.CoreTotal = device.Devcore
		gpu.MemoryTotal = device.Devmem
		gpu.NodeUid = device.NodeUid
		gpu.Health = device.Health
		gpu.Mode = device.Mode

		vGPU, core, memory, err := s.pod.StatisticsByDeviceId(ctx, device.AliasId)
		if err == nil {
			gpu.VgpuUsed = vGPU
			gpu.CoreUsed = core
			gpu.MemoryUsed = memory
		}
		return gpu, nil
	}
	return gpu, nil
}
