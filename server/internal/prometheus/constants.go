package prometheus

const (
	// GpuMemoryQuery 查询gpu显存大小
	GpuMemoryQuery = "hami_memory_size{node=\"%s\"}"
	// NumberOfDiskQuery 查询磁盘数量
	NumberOfDiskQuery = "count(\n  node_disk_info{\n    instance=\"%s\", \n    device=~\"sd[a-z]+|nvme[0-9]+n[0-9]+|vd[a-z]+\"\n  }\n)"
)
