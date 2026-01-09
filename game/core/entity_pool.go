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
type EntityPool[T Poolable] struct {
	pool    []T
	factory func() T
	mutex   sync.Mutex

	// Statistics
	totalCreated  int
	totalReused   int
	currentActive int
	maxActive     int
}

// NewEntityPool creates a new entity pool with a factory function
func NewEntityPool[T Poolable](factory func() T, initialCapacity int) *EntityPool[T] {
	pool := &EntityPool[T]{
		pool:    make([]T, 0, initialCapacity),
		factory: factory,
	}

	// Pre-allocate initial capacity
	for i := 0; i < initialCapacity; i++ {
		entity := factory()
		entity.SetActive(false)
		pool.pool = append(pool.pool, entity)
	}
	pool.totalCreated = initialCapacity

	return pool
}

// Get retrieves an entity from the pool or creates a new one
func (p *EntityPool[T]) Get() T {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// Try to find an inactive entity in the pool
	for i := range p.pool {
		if !p.pool[i].IsActive() {
			entity := p.pool[i]
			entity.Reset()
			entity.SetActive(true)
			p.totalReused++
			p.currentActive++
			if p.currentActive > p.maxActive {
				p.maxActive = p.currentActive
			}
			return entity
		}
	}

	// No inactive entities, create a new one
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

// Return returns an entity to the pool
func (p *EntityPool[T]) Return(entity T) {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if entity.IsActive() {
		entity.SetActive(false)
		p.currentActive--
	}
}

// ReturnAll returns all active entities to the pool
func (p *EntityPool[T]) ReturnAll() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for i := range p.pool {
		if p.pool[i].IsActive() {
			p.pool[i].SetActive(false)
		}
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
	p.currentActive = 0
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
