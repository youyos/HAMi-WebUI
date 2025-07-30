package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/prometheus/common/model"
	"log"
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
	log.Println("CreateResourcePool called", req)
	poolName := req.PoolName

	if database.ExistsResourcePoolByPoolName(poolName) {
		return &pb.BaseResponse{Code: -1, Message: "资源池：'" + poolName + "'已经存在"}, nil
	}

	poolId, err := database.InsertResourcePool(poolName)
	if err != nil {
		return &pb.BaseResponse{Code: -1, Message: poolName + "创建资源池失败"}, nil
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
		return &pb.BaseResponse{Code: -1, Message: poolName + "创建资源池失败"}, nil
	}

	log.Println("CreateResourcePool success", poolName, rows)
	return &pb.BaseResponse{Code: 1, Message: "成功"}, nil
}

func (s *ResourcePoolService) Update(ctx context.Context, req *pb.ResourcePoolUpdateRequest) (*pb.BaseResponse, error) {
	log.Println("UpdateResourcePool called", req)
	poolId := req.PoolId
	resourcePool, err := database.QueryResourcePoolById(poolId)
	if err != nil {
		return &pb.BaseResponse{Code: -1, Message: "更新资源池失败"}, nil
	}

	if resourcePool == nil {
		return &pb.BaseResponse{Code: -1, Message: "资源池不存在"}, nil
	}

	_, err = database.DeleteNodesByPoolId(poolId)
	if err != nil {
		return &pb.BaseResponse{Code: -1, Message: "更新资源池失败"}, nil
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
		return &pb.BaseResponse{Code: -1, Message: "更新资源池失败"}, nil
	}

	return &pb.BaseResponse{Code: 1, Message: "成功"}, nil
}

func (s *ResourcePoolService) Delete(ctx context.Context, req *pb.ResourcePoolDeleteRequest) (*pb.BaseResponse, error) {
	log.Println("DeleteResourcePool called", req)
	poolId := req.PoolId
	num, err := database.DeleteNodesByPoolId(poolId)
	if err != nil {
		return &pb.BaseResponse{Code: -1, Message: "删除资源池失败"}, nil
	}

	log.Println("DeleteNodes success", poolId, num)
	num, err = database.DeleteResourcePoolById(poolId)
	if err != nil {
		return &pb.BaseResponse{Code: -1, Message: "删除资源池失败"}, nil
	}

	log.Println("DeleteResourcePool success", poolId, num)
	return &pb.BaseResponse{Code: 1, Message: "成功"}, nil
}

func (s *ResourcePoolService) List(ctx context.Context, req *pb.ResourcePoolListRequest) (*pb.ResourcePoolListResponse, error) {
	log.Println("GetResourcePoolList", req)

	resourcePoolList, err := database.QueryResourcePoolListAll()
	if err != nil {
		return nil, errors.New("获取资源池列表失败")
	}

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
				return nil, errors.New("node not found: " + n.NodeName)
			}
			poolData.CpuCores = poolData.CpuCores + node.CPUCores
			poolData.GpuNum = poolData.GpuNum + node.GPUCount
			poolData.TotalMemory = poolData.TotalMemory + node.TotalMemory
			poolData.AvailableMemory = poolData.AvailableMemory + node.AvailableMemory
			poolData.DiskSize = poolData.DiskSize + node.DiskTotal

		}
		data = append(data, &poolData)
	}
	return &pb.ResourcePoolListResponse{Data: data}, nil
}

func (s *ResourcePoolService) GetDetail(ctx context.Context, req *pb.ResourcePoolDetailRequest) (*pb.ResourcePoolDetailResponse, error) {
	log.Println("GetResourcePoolDetail", req)

	gpuMemory := fmt.Sprintf(prometheus.GpuMemoryQuery, "k8s2")
	fmt.Println("gpuMemory", s.simpleQuery(ctx, gpuMemory))

	diskNumQuery := fmt.Sprintf(prometheus.NumberOfDiskQuery, "k8s1")

	fmt.Println("diskNum", s.simpleQuery(ctx, diskNumQuery))
	return &pb.ResourcePoolDetailResponse{Code: 1, Message: "成功11111"}, nil
}

func (s *ResourcePoolService) GetAvailableNodes(ctx context.Context, req *pb.AvailableNodesRequest) (*pb.AvailableNodesResponse, error) {
	log.Println("GetAvailableNodes", req)

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
		node.GPUCount = int64(len(node.Devices))
		b, _ := json.MarshalIndent(node, "", "  ")
		fmt.Println(string(b))
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
