/*
Package memorymonitor provides functionality for monitoring memory usage of a Go application and uploading memory profiles to a designated storage based on the provided Writer implementation.

The memorymonitor package monitors the memory usage of a Go application and writes a pprof memory profile to a designated storage when the memory usage exceeds a specified limit.

The behavior of the package is controlled by the following components:
- A Writer interface is used for uploading the pprof memory profile. The package is designed to be storage-agnostic. The actual storage destination (such as local disk, S3, or any other location) is determined by the provided implementation of the Writer interface.
- A Monitor interface is used for controlling the monitoring process. The WithMemoryLimit and WithMonitorFreq methods are used to customize the memory limit and monitor frequency respectively. The StartMonitoring method starts the monitoring process.
- The memory limit (memoryLimit) is set to 5 MB by default, but it can be customized using the WithMemoryLimit method.
- The monitor frequency (monitorFreq) is set to 10 seconds by default, but it can be customized using the WithMonitorFreq method.
- The StartMonitoring method starts the memory monitoring process, periodically checking the memory usage and uploading a memory profile if the memory limit is exceeded.
- The checkAndWriteProfile method checks the memory usage, triggers a garbage collection (GC), and uploads a memory profile to the storage specified by the Writer if the memory limit is exceeded.
- The memory profile is written in pprof format and includes information about memory allocations and usage.
- The memory profile file is named using the current timestamp and a unique ID.
*/
package memorymonitor

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"
)

const (
	defaultMemoryLimit      = 5 * 1024 * 1024
	defaultMonitorFrequency = 10 * time.Second
)

type Writer interface {
	Write(fileName string, buffer bytes.Buffer) error
}

type Monitor interface {
	StartMonitoring()
	WithMemoryLimit(limit uint64) *memory
	WithMonitorFreq(freq time.Duration) *memory
}

type memory struct {
	// memoryLimit holds the memory limit in Bytes
	memoryLimit uint64
	// monitorFreq holds the monitor frequency
	monitorFreq time.Duration
	// writer holds the Writer to write the memory profile
	writer Writer
}

func NewMemoryMonitor(w Writer) Monitor {
	return &memory{
		memoryLimit: defaultMemoryLimit,
		monitorFreq: defaultMonitorFrequency,
		writer:      w,
	}
}

func (m *memory) WithMemoryLimit(limit uint64) *memory {
	m.memoryLimit = limit
	return m
}

func (m *memory) WithMonitorFreq(freq time.Duration) *memory {
	m.monitorFreq = freq
	return m
}

func (m *memory) StartMonitoring() {
	ticker := time.NewTicker(m.monitorFreq)
	defer ticker.Stop()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	for {
		select {
		case <-ticker.C:
			m.checkAndWriteProfile()
		case <-sigCh:
			return
		}
	}
}

func (m *memory) checkAndWriteProfile() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	if memStats.Alloc < m.memoryLimit {
		return
	}

	runtime.GC()
	var buf bytes.Buffer
	if err := pprof.WriteHeapProfile(&buf); err != nil {
		return
	}

	currentTime := time.Now()
	uniqueId := int(currentTime.Unix())
	fileName := fmt.Sprintf("%s_%d.pprof", currentTime.Format("20060102150405"), uniqueId)

	// Write this pprof to somewhere which its client will decide by passing interface which has write func
	if err := m.writer.Write(fileName, buf); err != nil {
	}

}
