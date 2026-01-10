package core

import (
	"sync"
)

// Poolable interface for entities that can be pooled
type Poolable interface {
	Reset() // Reset entity to default state
	IsActive() bool
	SetActive(active bool)
}

// EntityPool is a generic object pool for reusable entities
// Uses a free list for O(1) Get() instead of linear search
type EntityPool[T Poolable] struct {
	pool     []T
	freeList []int // Stack of indices to inactive entities
	factory  func() T
	mutex    sync.Mutex

	// Statistics
	totalCreated  int
	totalReused   int
	currentActive int
	maxActive     int
}

// NewEntityPool creates a new entity pool with a factory function
func NewEntityPool[T Poolable](factory func() T, initialCapacity int) *EntityPool[T] {
	pool := &EntityPool[T]{
		pool:     make([]T, 0, initialCapacity),
		freeList: make([]int, 0, initialCapacity),
		factory:  factory,
	}

	// Pre-allocate initial capacity
	for i := 0; i < initialCapacity; i++ {
		entity := factory()
		entity.SetActive(false)
		pool.pool = append(pool.pool, entity)
		pool.freeList = append(pool.freeList, i) // All entities start as free
	}
	pool.totalCreated = initialCapacity

	return pool
}

// Get retrieves an entity from the pool or creates a new one
// Uses free list for O(1) retrieval instead of linear search
func (p *EntityPool[T]) Get() T {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Try to pop from free list (O(1))
	if len(p.freeList) > 0 {
		idx := p.freeList[len(p.freeList)-1]
		p.freeList = p.freeList[:len(p.freeList)-1]

		entity := p.pool[idx]
		entity.Reset()
		entity.SetActive(true)
		p.totalReused++
		p.currentActive++
		if p.currentActive > p.maxActive {
			p.maxActive = p.currentActive
		}
		return entity
	}

	// No free entities, create a new one and add to pool
	entity := p.factory()
	entity.SetActive(true)
	p.pool = append(p.pool, entity)
	p.totalCreated++
	p.currentActive++
	if p.currentActive > p.maxActive {
		p.maxActive = p.currentActive
	}

	return entity
}

// Return returns an entity to the pool by adding its index to the free list
func (p *EntityPool[T]) Return(entity T) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if entity.IsActive() {
		entity.SetActive(false)
		p.currentActive--

		// Find the entity's index and add to free list
		// Compare using any to work around generic comparison limitations
		for i := range p.pool {
			if any(p.pool[i]) == any(entity) {
				p.freeList = append(p.freeList, i)
				break
			}
		}
	}
}

// ReturnAll returns all active entities to the pool
func (p *EntityPool[T]) ReturnAll() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Clear free list and rebuild it with all indices
	p.freeList = p.freeList[:0]
	for i := range p.pool {
		if p.pool[i].IsActive() {
			p.pool[i].SetActive(false)
		}
		p.freeList = append(p.freeList, i)
	}
	p.currentActive = 0
}

// GetActive returns all active entities
func (p *EntityPool[T]) GetActive() []T {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	active := make([]T, 0, p.currentActive)
	for i := range p.pool {
		if p.pool[i].IsActive() {
			active = append(active, p.pool[i])
		}
	}
	return active
}

// GetStats returns pool statistics
func (p *EntityPool[T]) GetStats() PoolStats {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	return PoolStats{
		TotalCreated:  p.totalCreated,
		TotalReused:   p.totalReused,
		CurrentActive: p.currentActive,
		MaxActive:     p.maxActive,
		PoolSize:      len(p.pool),
		ReuseRate:     p.calculateReuseRate(),
	}
}

// calculateReuseRate returns the percentage of entities that were reused
func (p *EntityPool[T]) calculateReuseRate() float64 {
	total := p.totalCreated + p.totalReused
	if total == 0 {
		return 0
	}
	return float64(p.totalReused) / float64(total) * 100
}

// Clear clears the pool
func (p *EntityPool[T]) Clear() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.pool = p.pool[:0]
	p.freeList = p.freeList[:0]
	p.currentActive = 0
}

// TrimExcess removes excess inactive entities from pool to reduce memory usage.
// Keeps at least 'minCapacity' entities in pool. Only trims if pool size exceeds
// 2x the target size to avoid frequent reallocations.
func (p *EntityPool[T]) TrimExcess(minCapacity int) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Calculate target size: max of minCapacity or 2x current active
	targetSize := minCapacity
	if p.currentActive*2 > targetSize {
		targetSize = p.currentActive * 2
	}

	// Only trim if pool is significantly larger than needed
	if len(p.pool) <= targetSize {
		return
	}

	// Compact: keep active entities and fill remaining with inactive up to targetSize
	newPool := make([]T, 0, targetSize)
	newFreeList := make([]int, 0, targetSize)

	// First, add all active entities
	for i := range p.pool {
		if p.pool[i].IsActive() {
			newPool = append(newPool, p.pool[i])
		}
	}

	// Then fill with inactive up to targetSize, tracking their new indices
	for i := range p.pool {
		if len(newPool) >= targetSize {
			break
		}
		if !p.pool[i].IsActive() {
			newFreeList = append(newFreeList, len(newPool))
			newPool = append(newPool, p.pool[i])
		}
	}

	p.pool = newPool
	p.freeList = newFreeList
}

// ResetStats resets the pool statistics (useful after game restart)
func (p *EntityPool[T]) ResetStats() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	p.totalCreated = len(p.pool) // Count current pool as "created"
	p.totalReused = 0
	p.maxActive = p.currentActive
}

// PoolStats contains statistics about pool usage
type PoolStats struct {
	TotalCreated  int     // Total entities created
	TotalReused   int     // Total entities reused from pool
	CurrentActive int     // Current active entities
	MaxActive     int     // Maximum concurrent active entities
	PoolSize      int     // Total pool size
	ReuseRate     float64 // Reuse rate percentage
}

// SimplePool is a non-generic pool for backwards compatibility
type SimplePool struct {
	items []interface{}
	new   func() interface{}
	reset func(interface{})
	mutex sync.Mutex
}

// NewSimplePool creates a simple object pool
func NewSimplePool(new func() interface{}, reset func(interface{})) *SimplePool {
	return &SimplePool{
		items: make([]interface{}, 0, 32),
		new:   new,
		reset: reset,
	}
}

// Get retrieves an item from the pool
func (p *SimplePool) Get() interface{} {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if len(p.items) == 0 {
		return p.new()
	}

	item := p.items[len(p.items)-1]
	p.items = p.items[:len(p.items)-1]
	return item
}

// Put returns an item to the pool
func (p *SimplePool) Put(item interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.reset != nil {
		p.reset(item)
	}
	p.items = append(p.items, item)
}
