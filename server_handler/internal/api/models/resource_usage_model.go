package models

// ResourceUsage represents the system's resource usage.
type ResourceUsage struct {
	CPUUsage        float64 `json:"cpu_usage"`
	MemoryUsed      uint64  `json:"memory_used"`
	MemoryTotal     uint64  `json:"memory_total"`
	MemoryUsedPercent float64 `json:"memory_used_percent"`
	DiskUsed        uint64  `json:"disk_used"`
	DiskTotal       uint64  `json:"disk_total"`
	DiskUsedPercent float64 `json:"disk_used_percent"`
}
