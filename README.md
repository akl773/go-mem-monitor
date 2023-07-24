# Memory Monitor Package (go-mem-monitor)
Package memorymonitor provides functionality for monitoring memory usage of a Go application and uploading memory profiles to a designated storage based on the provided Writer implementation.

## Installation
To use this package in your Go project, you need to install it using go get:

```
go get github.com/akl773/go-mem-monitor
```

## Components
The behavior of the package is controlled by the following components:

* **Writer Interface**
  The Writer interface is used for uploading the pprof memory profile. The package does not impose any specific storage destination, allowing you to define your own implementation based on your requirements. Any location that satisfies the Writer interface can be used to store the memory profiles.

* **Monitor Interface**
  The Monitor interface is used for controlling the monitoring process. It allows you to customize the memory limit and monitor frequency. The available methods are as follows:

* ```StartMonitoring()```: Initiates the memory monitoring process, periodically checking the memory usage and uploading a memory profile if the memory limit is exceeded.
* ```WithMemoryLimit(limit uint64) *memory```: Sets a custom memory limit (in bytes) for triggering memory profile uploads.
* ```WithMonitorFreq(freq time.Duration) *memory```: Sets a custom monitor frequency for how often the memory usage is checked.

## Default Settings
The package comes with default settings:

* The memory limit (defaultMemoryLimit) is set to 5 MB (5 * 1024 * 1024 bytes) by default but can be customized using the WithMemoryLimit method.
* The monitor frequency (defaultMonitorFrequency) is set to 10 seconds by default but can be customized using the WithMonitorFreq method.

## Usage

	1. Import the package
	2. Create a custom implementation of the Writer interface to define where the memory profiles should be stored.
	3. Initialize the memory monitor using the NewMemoryMonitor function, passing your custom Writer implementation as an argument.
	4. Optionally, customize the memory limit and monitor frequency using the WithMemoryLimit and WithMonitorFreq methods.
	5. Start the memory monitoring process by calling the StartMonitoring method.

## Example

```
package main

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"runtime/pprof"
	"syscall"
	"time"

	"github.com/akl773/go-mem-monitor"
)

// CustomWriter is an example implementation of the Writer interface.
type CustomWriter struct {
	// Define your custom storage destination and authentication details here.
}

func (cw *CustomWriter) Write(fileName string, buffer bytes.Buffer) error {
	// Implement the logic to upload the memory profile to your custom storage.
	return nil
}

func main() {
	// Create an instance of your custom writer.
	customWriter := &CustomWriter{}

	// Initialize the memory monitor with your custom writer.
	monitor := memorymonitor.NewMemoryMonitor(customWriter)

	// Optionally, customize the memory limit and monitor frequency.
	// monitor.WithMemoryLimit(10 * 1024 * 1024) // 10 MB
	// monitor.WithMonitorFreq(5 * time.Second)

	// Start the memory monitoring process.
	monitor.StartMonitoring()
}

```

## Note

* The memory profile is written in pprof format and includes information about memory allocations and usage.
* The memory profile file is named using the current timestamp and a unique ID to avoid overwriting previous profiles.
* The memory monitoring process triggers a garbage collection (GC) before writing the memory profile to provide more accurate memory usage information.
  Feel free to use this package and customize it according to your specific needs. If you encounter any issues or have suggestions for improvements, please don't hesitate to contribute to the project. Happy coding!
