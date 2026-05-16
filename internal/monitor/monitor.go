package monitor

import (
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/mem"
)

// Stats holds the latest snapshot of system metrics
type Stats struct {
	CPUUsagePercent    float64   `json:"cpu_usage_percent"`
	MemoryTotalMB      uint64    `json:"memory_total_mb"`
	MemoryUsedMB       uint64    `json:"memory_used_mb"`
	MemoryFreeMB       uint64    `json:"memory_free_mb"`
	MemoryUsagePercent float64   `json:"memory_usage_percent"`
	DiskTotalGB        uint64    `json:"disk_total_gb"`
	DiskUsedGB         uint64    `json:"disk_used_gb"`
	DiskFreeGB         uint64    `json:"disk_free_gb"`
	DiskUsagePercent   float64   `json:"disk_usage_percent"`
	UpdatedAt          time.Time `json:"updated_at"`
}

// Monitor is responsible for collecting and storing system metrics
type Monitor struct {
	mu      sync.RWMutex
	stats   Stats
	stopCh  chan struct{}
}

// NewMonitor creates a new instance of Monitor
func NewMonitor() *Monitor {
	return &Monitor{
		stopCh: make(chan struct{}),
	}
}

// Start begins the background worker that polls metrics periodically
func (m *Monitor) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	
	// Initial fetch
	m.fetchAndStore()

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				m.fetchAndStore()
			case <-m.stopCh:
				return
			}
		}
	}()
}

// Stop halts the background worker
func (m *Monitor) Stop() {
	close(m.stopCh)
}

// GetStats returns a thread-safe copy of the latest metrics
func (m *Monitor) GetStats() Stats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.stats
}

func (m *Monitor) fetchAndStore() {
	var s Stats

	// CPU
	cpuPercents, err := cpu.Percent(0, false)
	if err == nil && len(cpuPercents) > 0 {
		s.CPUUsagePercent = cpuPercents[0]
	} else {
		log.Printf("Error fetching CPU metrics: %v", err)
	}

	// Memory
	vmStat, err := mem.VirtualMemory()
	if err == nil {
		s.MemoryTotalMB = vmStat.Total / 1024 / 1024
		s.MemoryUsedMB = vmStat.Used / 1024 / 1024
		s.MemoryFreeMB = vmStat.Available / 1024 / 1024
		s.MemoryUsagePercent = vmStat.UsedPercent
	} else {
		log.Printf("Error fetching Memory metrics: %v", err)
	}

	// Disk
	diskPath := "/"
	if runtime.GOOS == "windows" {
		diskPath = "C:"
	}
	diskStat, err := disk.Usage(diskPath)
	if err == nil {
		s.DiskTotalGB = diskStat.Total / 1024 / 1024 / 1024
		s.DiskUsedGB = diskStat.Used / 1024 / 1024 / 1024
		s.DiskFreeGB = diskStat.Free / 1024 / 1024 / 1024
		s.DiskUsagePercent = diskStat.UsedPercent
	} else {
		log.Printf("Error fetching Disk metrics: %v", err)
	}

	s.UpdatedAt = time.Now()

	m.mu.Lock()
	defer m.mu.Unlock()
	m.stats = s
}
