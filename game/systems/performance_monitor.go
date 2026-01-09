package systems

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// PerformanceMonitor tracks game performance metrics
type PerformanceMonitor struct {
	// FPS tracking
	frameCount    int
	fps           float64
	lastFPSUpdate time.Time

	// Frame time tracking
	frameTimes   []time.Duration
	frameTimeIdx int
	avgFrameTime time.Duration
	maxFrameTime time.Duration
	minFrameTime time.Duration

	// Memory tracking
	lastMemStats     runtime.MemStats
	lastMemStatsTime time.Time
	allocRate        float64 // Bytes per second

	// Entity counts
	entityCounts map[string]int

	// Update counts for profiling
	updateCount   int
	collisionTime time.Duration
	renderTime    time.Duration

	mu sync.RWMutex
}

// NewPerformanceMonitor creates a new performance monitor
func NewPerformanceMonitor() *PerformanceMonitor {
	return &PerformanceMonitor{
		lastFPSUpdate:    time.Now(),
		frameTimes:       make([]time.Duration, 60), // Track last 60 frames
		minFrameTime:     time.Hour,                 // Start with large value
		entityCounts:     make(map[string]int),
		lastMemStatsTime: time.Now(),
	}
}

// RecordFrame records metrics for a frame
func (pm *PerformanceMonitor) RecordFrame(frameTime time.Duration) {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	pm.frameCount++
	pm.frameTimes[pm.frameTimeIdx] = frameTime
	pm.frameTimeIdx = (pm.frameTimeIdx + 1) % len(pm.frameTimes)

	// Update min/max
	if frameTime > pm.maxFrameTime {
		pm.maxFrameTime = frameTime
	}
	if frameTime < pm.minFrameTime && frameTime > 0 {
		pm.minFrameTime = frameTime
	}

	// Calculate average frame time
	var total time.Duration
	count := 0
	for _, ft := range pm.frameTimes {
		if ft > 0 {
			total += ft
			count++
		}
	}
	if count > 0 {
		pm.avgFrameTime = total / time.Duration(count)
	}

	// Update FPS every second
	now := time.Now()
	elapsed := now.Sub(pm.lastFPSUpdate)
	if elapsed >= time.Second {
		pm.fps = float64(pm.frameCount) / elapsed.Seconds()
		pm.frameCount = 0
		pm.lastFPSUpdate = now
	}
}

// RecordCollisionTime records time spent in collision detection
func (pm *PerformanceMonitor) RecordCollisionTime(duration time.Duration) {
	pm.mu.Lock()
	pm.collisionTime = duration
	pm.mu.Unlock()
}

// RecordRenderTime records time spent in rendering
func (pm *PerformanceMonitor) RecordRenderTime(duration time.Duration) {
	pm.mu.Lock()
	pm.renderTime = duration
	pm.mu.Unlock()
}

// UpdateEntityCount updates the count for a specific entity type
func (pm *PerformanceMonitor) UpdateEntityCount(entityType string, count int) {
	pm.mu.Lock()
	pm.entityCounts[entityType] = count
	pm.mu.Unlock()
}

// UpdateMemoryStats updates memory statistics
func (pm *PerformanceMonitor) UpdateMemoryStats() {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(pm.lastMemStatsTime)

	if elapsed >= time.Second {
		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)

		// Calculate allocation rate
		allocDiff := memStats.TotalAlloc - pm.lastMemStats.TotalAlloc
		pm.allocRate = float64(allocDiff) / elapsed.Seconds()

		pm.lastMemStats = memStats
		pm.lastMemStatsTime = now
	}
}

// GetFPS returns current FPS
func (pm *PerformanceMonitor) GetFPS() float64 {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.fps
}

// GetAverageFrameTime returns average frame time
func (pm *PerformanceMonitor) GetAverageFrameTime() time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.avgFrameTime
}

// GetFrameTimeStats returns min, max, and average frame times
func (pm *PerformanceMonitor) GetFrameTimeStats() (min, max, avg time.Duration) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.minFrameTime, pm.maxFrameTime, pm.avgFrameTime
}

// GetMemoryStats returns current memory statistics
func (pm *PerformanceMonitor) GetMemoryStats() (alloc, totalAlloc, sys uint64, allocRate float64) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.lastMemStats.Alloc, pm.lastMemStats.TotalAlloc, pm.lastMemStats.Sys, pm.allocRate
}

// GetEntityCounts returns entity counts
func (pm *PerformanceMonitor) GetEntityCounts() map[string]int {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	counts := make(map[string]int)
	for k, v := range pm.entityCounts {
		counts[k] = v
	}
	return counts
}

// GetCollisionTime returns time spent in collision detection
func (pm *PerformanceMonitor) GetCollisionTime() time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.collisionTime
}

// GetRenderTime returns time spent in rendering
func (pm *PerformanceMonitor) GetRenderTime() time.Duration {
	pm.mu.RLock()
	defer pm.mu.RUnlock()
	return pm.renderTime
}

// GetSummary returns a formatted summary of all metrics
func (pm *PerformanceMonitor) GetSummary() string {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	min, max, avg := pm.minFrameTime, pm.maxFrameTime, pm.avgFrameTime
	alloc := pm.lastMemStats.Alloc
	totalAlloc := pm.lastMemStats.TotalAlloc
	sys := pm.lastMemStats.Sys
	numGC := pm.lastMemStats.NumGC

	summary := fmt.Sprintf(`Performance Summary:
  FPS: %.1f
  Frame Time: avg=%v, min=%v, max=%v
  Memory: alloc=%s, total=%s, sys=%s
  GC Count: %d
  Alloc Rate: %.2f MB/s
  Collision Time: %v
  Render Time: %v
  Entity Counts:`,
		pm.fps,
		avg, min, max,
		formatBytes(alloc), formatBytes(totalAlloc), formatBytes(sys),
		numGC,
		pm.allocRate/1024/1024,
		pm.collisionTime,
		pm.renderTime,
	)

	for entityType, count := range pm.entityCounts {
		summary += fmt.Sprintf("\n    %s: %d", entityType, count)
	}

	return summary
}

// ResetMinMax resets min/max statistics
func (pm *PerformanceMonitor) ResetMinMax() {
	pm.mu.Lock()
	defer pm.mu.Unlock()
	pm.minFrameTime = time.Hour
	pm.maxFrameTime = 0
}

// formatBytes formats byte count as human-readable string
func formatBytes(bytes uint64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := uint64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}
