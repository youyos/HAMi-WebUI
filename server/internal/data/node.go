package data

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	k8stypes "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/informers"
	listerscorev1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/tools/cache"
	"strings"
	"sync"
	"time"
	"vgpu/internal/biz"
	"vgpu/internal/provider"
	"vgpu/internal/provider/ascend"
	"vgpu/internal/provider/hygon"
	"vgpu/internal/provider/mlu"
	"vgpu/internal/provider/nvidia"
)

type nodeRepo struct {
	data       *Data
	nodeNotify chan struct{}
	nodeLister listerscorev1.NodeLister
	nodes      map[k8stypes.UID]*biz.Node
	allNodes   []*biz.Node
	log        *log.Helper
	mutex      sync.RWMutex
	providers  []provider.Provider
}

// NewNodeRepo .
func NewNodeRepo(data *Data, nodeSelectors map[string]string, logger log.Logger) biz.NodeRepo {
	nodeRepo := &nodeRepo{
		data:       data,
		nodeNotify: make(chan struct{}, 1),
		nodes:      map[k8stypes.UID]*biz.Node{},
		log:        log.NewHelper(logger),
		providers: []provider.Provider{
			nvidia.NewNvidia(data.promCl, log.NewHelper(logger), nodeSelectors[biz.NvidiaGPUDevice]),
			mlu.NewCambricon(data.promCl, log.NewHelper(logger), nodeSelectors[biz.CambriconGPUDevice]),
			ascend.NewAscend(data.promCl, log.NewHelper(logger), nodeSelectors[biz.AscendGPUDevice]),
			hygon.NewHygon(data.promCl, log.NewHelper(logger), nodeSelectors[biz.HygonGPUDevice]),
		},
	}
	nodeRepo.init()
	return nodeRepo
}

func (r *nodeRepo) updateLocalNodes() {
	for {
		select {
		case <-r.nodeNotify:
		}
		n := map[k8stypes.UID]*biz.Node{}
		for _, p := range r.providers {
			selector, err := p.GetNodeDevicePluginLabels()
			if err != nil {
				r.log.Warnf("create selector use labels failed :%s", err)
				return
			}
			nodes, err := r.nodeLister.List(selector)
			if err != nil {
				r.log.Warnf("list node info error: %s", err)
				continue
			}
			for _, node := range nodes {
				bizNode := r.fetchNodeInfo(node)
				if _, ok := n[node.UID]; !ok {
					n[node.UID] = bizNode
				}

				devices, err := p.FetchDevices(node)
				if err != nil {
					r.log.Warnf("list devices info error: %s", err)
					continue
				}
				for _, device := range devices {
					bizNode.Devices = append(bizNode.Devices, &biz.DeviceInfo{
						Index:    int(device.Index),
						Id:       device.ID,
						AliasId:  device.AliasId,
						Count:    device.Count,
						Devmem:   device.Devmem,
						Devcore:  device.Devcore,
						Type:     device.Type,
						Numa:     device.Numa,
						Mode:     device.Mode,
						Health:   device.Health,
						NodeName: node.Name,
						NodeUid:  string(node.UID),
						Provider: p.GetProvider(),
						Driver:   device.Driver,
					})
				}
			}
		}
		r.nodes = n

		var all []*biz.Node
		allNodes, _ := r.nodeLister.List(labels.Everything())
		for _, node := range allNodes {
			bizNode := r.fetchNodeInfo(node)
			gpuNode := n[k8stypes.UID(bizNode.Uid)]
			if gpuNode != nil {
				bizNode.Devices = gpuNode.Devices
			}
			all = append(all, bizNode)
		}
		r.allNodes = all
	}
}

func (r *nodeRepo) init() {
	go r.updateLocalNodes()
	informerFactory := informers.NewSharedInformerFactoryWithOptions(r.data.k8sCl, time.Hour*1)
	r.nodeLister = informerFactory.Core().V1().Nodes().Lister()
	informer := informerFactory.Core().V1().Nodes().Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc:    r.onAddNode,
		UpdateFunc: r.onUpdateNode,
		DeleteFunc: r.onDeletedNode,
	})
	stopCh := make(chan struct{})
	informerFactory.Start(stopCh)
	informerFactory.WaitForCacheSync(stopCh)
}

func (r *nodeRepo) onAddNode(obj interface{}) {
	r.nodeNotify <- struct{}{}
}

func (r *nodeRepo) onUpdateNode(old interface{}, new interface{}) {
	r.nodeNotify <- struct{}{}
}

func (r *nodeRepo) onDeletedNode(obj interface{}) {
	r.nodeNotify <- struct{}{}
}

func (r *nodeRepo) fetchNodeInfo(node *corev1.Node) *biz.Node {
	//b, _ := json.MarshalIndent(node, "", "  ")
	//fmt.Println(string(b))
	n := &biz.Node{IsSchedulable: !node.Spec.Unschedulable}
	for _, addr := range node.Status.Addresses {
		if addr.Type == corev1.NodeInternalIP {
			n.IP = addr.Address
			break
		}
	}
	for _, cond := range node.Status.Conditions {
		if cond.Type == corev1.NodeReady && cond.Status == corev1.ConditionTrue {
			n.IsReady = true
		}
	}
	n.Uid = string(node.UID)
	n.Name = node.Name
	n.OSImage = node.Status.NodeInfo.OSImage
	n.OperatingSystem = strings.ToUpper(node.Status.NodeInfo.OperatingSystem[:1]) + strings.ToLower(node.Status.NodeInfo.OperatingSystem[1:])
	n.KernelVersion = node.Status.NodeInfo.KernelVersion
	n.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion
	n.KubeletVersion = node.Status.NodeInfo.KubeletVersion
	n.KubeProxyVersion = node.Status.NodeInfo.KubeProxyVersion
	n.Architecture = strings.ToUpper(node.Status.NodeInfo.Architecture)
	n.CreationTimestamp = node.CreationTimestamp.Format("2006-01-02 15:04:05")

	capacity := node.Status.Capacity
	allocatable := node.Status.Allocatable
	// CPU 核数
	if cpu, ok := capacity[corev1.ResourceCPU]; ok {
		n.CPUCores = cpu.Value()
	}

	// GPU 数量（nvidia.com/gpu）
	//if gpu, ok := capacity["nvidia.com/gpu"]; ok {
	//	n.GPUCount = gpu.Value()
	//}

	// 总内存
	if mem, ok := capacity[corev1.ResourceMemory]; ok {
		n.TotalMemory = mem.Value()
	}

	// 可用内存
	if mem, ok := allocatable[corev1.ResourceMemory]; ok {
		n.AvailableMemory = mem.Value()
	}

	// 磁盘总大小（临时存储）
	if disk, ok := capacity[corev1.ResourceEphemeralStorage]; ok {
		n.DiskTotal = disk.Value()
	}

	return n
}

func (r *nodeRepo) ListAll(context.Context) ([]*biz.Node, error) {
	var nodeList []*biz.Node
	for _, node := range r.nodes {
		nodeList = append(nodeList, node)
	}
	return nodeList, nil
}

func (r *nodeRepo) ListAllV2(context.Context) ([]*biz.Node, error) {
	var nodeList []*biz.Node
	for _, node := range r.allNodes {
		nodeList = append(nodeList, node)
	}
	return nodeList, nil
}

func (r *nodeRepo) GetNode(_ context.Context, uid string) (*biz.Node, error) {
	if _, ok := r.nodes[k8stypes.UID(uid)]; !ok {
		return nil, errors.New("node not found")
	}
	return r.nodes[k8stypes.UID(uid)], nil
}

func (r *nodeRepo) ListAllDevices(context.Context) ([]*biz.DeviceInfo, error) {
	var deviceList []*biz.DeviceInfo
	for _, node := range r.nodes {
		deviceList = append(deviceList, node.Devices...)
	}
	return deviceList, nil
}

func (r *nodeRepo) FindDeviceByAliasId(aliasId string) (*biz.DeviceInfo, error) {
	for _, node := range r.nodes {
		for _, d := range node.Devices {
			if d.AliasId == aliasId {
				return d, nil
			}
		}
	}
	return nil, errors.New(fmt.Sprintf("aliasID:%s device not found", aliasId))
}
