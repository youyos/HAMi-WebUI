package biz

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
)

type Node struct {
	Name                    string
	IP                      string
	IsSchedulable           bool
	IsReady                 bool
	Uid                     string
	OSImage                 string
	OperatingSystem         string
	KernelVersion           string
	ContainerRuntimeVersion string
	KubeletVersion          string
	KubeProxyVersion        string
	Architecture            string
	CreationTimestamp       string
	Devices                 []*DeviceInfo
	CPUCores                int64 // CPU 核数
	GPUCount                int64 // 显卡数量
	TotalMemory             int64 // 总内存（字节）
	AvailableMemory         int64 // 可用内存（字节）
	DiskTotal               int64 // 磁盘总大小（字节）
	StorageNum              int64
}

type DeviceInfo struct {
	Index    int
	Id       string
	AliasId  string
	Count    int32
	Devmem   int32
	Devcore  int32
	Type     string
	Numa     int
	Mode     string
	Health   bool
	NodeName string
	NodeUid  string
	Provider string
	Driver   string
}

type DeviceTotal struct {
	VgpuCount int32
	Cores     int32
	Memory    int32
}

type NodeRepo interface {
	ListAll(context.Context) ([]*Node, error)
	ListAllV2(context.Context) ([]*Node, error)
	GetNode(context.Context, string) (*Node, error)
	ListAllDevices(context.Context) ([]*DeviceInfo, error)
	FindDeviceByAliasId(string) (*DeviceInfo, error)
	EnableNode(context.Context, string) error
	DisableNode(context.Context, string) error
}

type NodeUsecase struct {
	repo NodeRepo
	log  *log.Helper
}

func NewNodeUsecase(repo NodeRepo, logger log.Logger) *NodeUsecase {
	return &NodeUsecase{repo: repo, log: log.NewHelper(logger)}
}

func (uc *NodeUsecase) ListAllNodes(ctx context.Context) ([]*Node, error) {
	return uc.repo.ListAll(ctx)
}

func (uc *NodeUsecase) ListAllNodesV2(ctx context.Context) ([]*Node, error) {
	return uc.repo.ListAllV2(ctx)
}

func (uc *NodeUsecase) GetNode(ctx context.Context, nodeName string) (*Node, error) {
	return uc.repo.GetNode(ctx, nodeName)
}

func (uc *NodeUsecase) ListAllDevices(ctx context.Context) ([]*DeviceInfo, error) {
	return uc.repo.ListAllDevices(ctx)
}

func (uc *NodeUsecase) FindDeviceByAliasId(aliasId string) (*DeviceInfo, error) {
	return uc.repo.FindDeviceByAliasId(aliasId)
}

func (uc *NodeUsecase) EnableNode(ctx context.Context, nodeName string) error {
	return uc.repo.EnableNode(ctx, nodeName)
}

func (uc *NodeUsecase) DisableNode(ctx context.Context, nodeName string) error {
	return uc.repo.DisableNode(ctx, nodeName)
}
