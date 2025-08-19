package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/gookit/goutil/arrutil"
	"github.com/prometheus/common/model"
	"sort"
	pb "vgpu/api/v1"
	"vgpu/internal/biz"
	"vgpu/internal/database"
	"vgpu/internal/prometheus"
)

type ResourcePoolService struct {
	pb.UnimplementedCardServer

	uc      *biz.NodeUsecase
	pod     *biz.PodUseCase
	summary *biz.SummaryUseCase
	ms      *MonitorService
}

func NewResourcePoolService(uc *biz.NodeUsecase, pod *biz.PodUseCase, summary *biz.SummaryUseCase, ms *MonitorService) *ResourcePoolService {
	return &ResourcePoolService{uc: uc, pod: pod, summary: summary, ms: ms}
}

func (s *ResourcePoolService) Create(ctx context.Context, req *pb.ResourcePoolCreateRequest) (*pb.BaseResponse, error) {
	log.Info("CreateResourcePool called", req)
	poolName := req.PoolName

	if database.ExistsResourcePoolByPoolName(poolName) {
		return &pb.BaseResponse{Code: 500, Message: "资源池：'" + poolName + "'已经存在"}, nil
	}

	poolId, err := database.InsertResourcePool(poolName)
	if err != nil {
		return &pb.BaseResponse{Code: 500, Message: poolName + "创建资源池失败"}, nil
	}

	nodes := make([]*database.NodeInfo, 0, len(req.Nodes))
	for _, node := range req.Nodes {
		nodes = append(nodes, &database.NodeInfo{
			Name: node.NodeName,
			IP:   node.NodeIp,
		})
	}

	rows, err := database.InsertNodes(poolId, nodes)
	if err != nil {
		return &pb.BaseResponse{Code: 500, Message: poolName + "创建资源池失败"}, nil
	}

	log.Info("CreateResourcePool success", poolName, rows)
	return &pb.BaseResponse{Code: 200, Message: "成功"}, nil
}

func (s *ResourcePoolService) Update(ctx context.Context, req *pb.ResourcePoolUpdateRequest) (*pb.BaseResponse, error) {
	log.Info("UpdateResourcePool called", req)
	poolId := req.PoolId
	resourcePool, err := database.QueryResourcePoolById(poolId)
	if err != nil {
		return &pb.BaseResponse{Code: 500, Message: "更新资源池失败"}, nil
	}

	if resourcePool == nil {
		return &pb.BaseResponse{Code: 500, Message: "资源池不存在"}, nil
	}

	_, err = database.DeleteNodesByPoolId(poolId)
	if err != nil {
		return &pb.BaseResponse{Code: 500, Message: "更新资源池失败"}, nil
	}

	nodes := make([]*database.NodeInfo, 0, len(req.Nodes))
	for _, node := range req.Nodes {
		nodes = append(nodes, &database.NodeInfo{
			Name: node.NodeName,
			IP:   node.NodeIp,
		})
	}
	_, err = database.InsertNodes(poolId, nodes)
	_, err = database.UpdateResourcePool(poolId, req.PoolName)
	if err != nil {
		return &pb.BaseResponse{Code: 500, Message: "更新资源池失败"}, nil
	}

	return &pb.BaseResponse{Code: 200, Message: "成功"}, nil
}

func (s *ResourcePoolService) Delete(ctx context.Context, req *pb.ResourcePoolDeleteRequest) (*pb.BaseResponse, error) {
	log.Info("DeleteResourcePool called", req)
	poolId := req.PoolId
	num, err := database.DeleteNodesByPoolId(poolId)
	if err != nil {
		return &pb.BaseResponse{Code: 500, Message: "删除资源池失败"}, nil
	}

	log.Infof("DeleteNodes success, poolId: %d, 影响行数: %d", poolId, num)
	num, err = database.DeleteResourcePoolById(poolId)
	if err != nil {
		return &pb.BaseResponse{Code: 500, Message: "删除资源池失败"}, nil
	}

	log.Infof("DeleteResourcePool success poolId: %d, 影响行数: %d", poolId, num)
	return &pb.BaseResponse{Code: 200, Message: "成功"}, nil
}

func (s *ResourcePoolService) RemoveNode(ctx context.Context, req *pb.RemoveNodeRequest) (*pb.BaseResponse, error) {
	log.Info("RemoveNode called", req)
	nodeId := req.NodeId
	num, err := database.DeleteNodeById(nodeId)
	if err != nil {
		return &pb.BaseResponse{Code: 500, Message: "移除节点失败"}, nil
	}
	log.Infof("RemoveNode success poolId: %d, 影响行数: %d", nodeId, num)
	return &pb.BaseResponse{Code: 200, Message: "成功"}, nil
}

func (s *ResourcePoolService) List(ctx context.Context, req *pb.ResourcePoolListRequest) (*pb.ResourcePoolListResponse, error) {
	log.Info("GetResourcePoolList", req)

	resourcePoolList, err := database.QueryResourcePoolListAll()
	if err != nil {
		return nil, errors.New("获取资源池列表失败")
	}

	// Sort the resourcePoolList to put ID 1 first
	sort.Slice(resourcePoolList, func(i, j int) bool {
		if resourcePoolList[i].Id == 1 {
			return true
		}
		if resourcePoolList[j].Id == 1 {
			return false
		}
		return resourcePoolList[i].CreateTime.After(resourcePoolList[j].CreateTime)
	})

	var data []*pb.ResourcePoolListData
	k8sNodes := s.getK8sNodes(ctx)
	for _, resourcePool := range resourcePoolList {
		var poolData pb.ResourcePoolListData
		poolData.PoolId = resourcePool.Id
		poolData.PoolName = resourcePool.PoolName

		dbNodes, _ := database.QueryNodesByPoolId(resourcePool.Id)
		poolData.NodeNum = int64(len(dbNodes))

		for _, n := range dbNodes {
			node := k8sNodes[n.NodeName]
			if node == nil {
				continue
			}
			poolData.CpuCores = poolData.CpuCores + node.CPUCores
			poolData.GpuNum = poolData.GpuNum + node.GPUCount
			poolData.TotalMemory = poolData.TotalMemory + node.TotalMemory
			poolData.AvailableMemory = poolData.AvailableMemory + node.AvailableMemory
			poolData.DiskSize = poolData.DiskSize + node.DiskTotal
			poolData.NodeList = append(poolData.NodeList, &pb.Nodes{
				NodeIp:   n.NodeIp,
				NodeName: n.NodeName,
			})
		}
		data = append(data, &poolData)
	}

	listData := data[0]
	linkUrl, _ := database.Get("big_model_resource_pool_link_url")
	listData.LinkUrl = linkUrl.(string)
	return &pb.ResourcePoolListResponse{Data: data}, nil
}

func (s *ResourcePoolService) GetDetail(ctx context.Context, req *pb.ResourcePoolDetailRequest) (*pb.ResourcePoolDetailResponse, error) {
	log.Info("GetResourcePoolDetail", req)
	poolNodes, err := database.QueryNodesByPoolId(req.PoolId)
	if err != nil {
		return nil, err
	}
	if len(poolNodes) == 0 {
		return &pb.ResourcePoolDetailResponse{}, nil
	}
	log.Info("GetResourcePoolDetail success", poolNodes)
	var res = &pb.ResourcePoolDetailResponse{List: []*pb.PoolNodeReply{}}
	nodes, err := s.uc.ListAllNodesV2(ctx)

	for _, poolNode := range poolNodes {
		node := s.filterNode(poolNode.NodeIp, nodes)
		if node == nil {
			continue
		}
		nodeReply, err := s.buildNodeReply(ctx, node)
		nodeReply.NodeId = poolNode.Id
		if err != nil {
			return nil, err
		}
		res.List = append(res.List, nodeReply)
	}

	sort.SliceStable(res.List, func(i, j int) bool {
		return res.List[i].Name < res.List[j].Name
	})

	return res, nil
}

func (s *ResourcePoolService) GetAvailableNodes(ctx context.Context, req *pb.AvailableNodesRequest) (*pb.AvailableNodesResponse, error) {
	log.Info("GetAvailableNodes", req)

	var data []*pb.AvailableNodesInfo
	k8sNodes := s.getK8sNodes(ctx)
	for _, node := range k8sNodes {
		var info pb.AvailableNodesInfo
		if node.Devices != nil {
			gpuMemoryQuery := fmt.Sprintf(prometheus.GpuMemoryQuery, node.Name)
			info.GpuMemory = s.simpleQuery(ctx, gpuMemoryQuery) * 1024 * 1024
			info.GpuNum = int64(len(node.Devices))
		} else {
			info.GpuMemory = 0
			info.GpuNum = 0
		}
		info.NodeName = node.Name
		info.NodeIp = node.IP
		info.TotalMemory = node.TotalMemory
		info.CpuCores = node.CPUCores
		info.DiskSize = node.DiskTotal
		data = append(data, &info)
	}
	return &pb.AvailableNodesResponse{Data: data}, nil
}

func (s *ResourcePoolService) getK8sNodes(ctx context.Context) map[string]*biz.Node {
	nodes, _ := s.uc.ListAllNodesV2(ctx)
	m := make(map[string]*biz.Node)
	for _, node := range nodes {
		if node.Lables["gpu"] != "on" {
			continue
		}
		node.GPUCount = int64(len(node.Devices))
		m[node.Name] = node
	}
	return m
}

func (s *ResourcePoolService) simpleQuery(ctx context.Context, query string) int64 {
	value, err := s.ms.promClient.Query(ctx, query)
	if err != nil {
		return 0
	}

	vector, ok := value.(model.Vector)
	if !ok || len(vector) == 0 {
		return 0
	}

	// 提取第一个样本的值（24576）
	sampleValue := int64(vector[0].Value)
	return sampleValue
}

func (s *ResourcePoolService) filterNode(nodeIp string, nodes []*biz.Node) *biz.Node {
	for _, node := range nodes {
		if node.IP == nodeIp {
			return node
		}
	}

	return nil
}

func (s *ResourcePoolService) buildNodeReply(ctx context.Context, node *biz.Node) (*pb.PoolNodeReply, error) {
	nodeReply := &pb.PoolNodeReply{
		Name:                    node.Name,
		Uid:                     node.Uid,
		Ip:                      node.IP,
		IsSchedulable:           node.IsSchedulable,
		IsReady:                 node.IsReady,
		OsImage:                 node.OSImage,
		OperatingSystem:         node.OperatingSystem,
		KernelVersion:           node.KernelVersion,
		ContainerRuntimeVersion: node.ContainerRuntimeVersion,
		KubeletVersion:          node.KubeletVersion,
		KubeProxyVersion:        node.KubeProxyVersion,
		Architecture:            node.Architecture,
		CreationTimestamp:       node.CreationTimestamp,
		CpuCores:                node.CPUCores,
		TotalMemory:             node.TotalMemory,
		DiskSize:                node.DiskTotal,
	}

	for _, device := range node.Devices {
		nodeReply.Type = append(nodeReply.Type, device.Type)
		nodeReply.VgpuTotal += device.Count
		nodeReply.CoreTotal += device.Devcore
		nodeReply.MemoryTotal += device.Devmem
		vGPU, core, memory, err := s.pod.StatisticsByDeviceId(ctx, device.AliasId)
		if err == nil {
			nodeReply.VgpuUsed += vGPU
			nodeReply.CoreUsed += core
			nodeReply.MemoryUsed += memory
		}
	}

	nodeReply.Type = arrutil.Unique(nodeReply.Type)
	nodeReply.CardCnt = int32(len(node.Devices))

	return nodeReply, nil
}
