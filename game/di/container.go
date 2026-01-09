package di

import (
	"fmt"
	"sync"
)

// ServiceLifetime defines the lifetime of a service
type ServiceLifetime int

const (
	// Singleton services are created once and reused
	Singleton ServiceLifetime = iota
	// Transient services are created each time they're requested
	Transient
)

// ServiceDescriptor describes how to create a service
type ServiceDescriptor struct {
	Lifetime ServiceLifetime
	Factory  func(c *Container) (interface{}, error)
}

// Container is a simple dependency injection container
type Container struct {
	services  map[string]ServiceDescriptor
	instances map[string]interface{}
	mutex     sync.RWMutex
}

// NewContainer creates a new dependency injection container
func NewContainer() *Container {
	return &Container{
		services:  make(map[string]ServiceDescriptor),
		instances: make(map[string]interface{}),
	}
}

// Register registers a service with the container
func (c *Container) Register(name string, lifetime ServiceLifetime, factory func(*Container) (interface{}, error)) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.services[name] = ServiceDescriptor{
		Lifetime: lifetime,
		Factory:  factory,
	}
}

// RegisterSingleton registers a singleton service
func (c *Container) RegisterSingleton(name string, factory func(*Container) (interface{}, error)) {
	c.Register(name, Singleton, factory)
}

// RegisterTransient registers a transient service
func (c *Container) RegisterTransient(name string, factory func(*Container) (interface{}, error)) {
	c.Register(name, Transient, factory)
}

// RegisterInstance registers an existing instance as a singleton
func (c *Container) RegisterInstance(name string, instance interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.instances[name] = instance
	c.services[name] = ServiceDescriptor{
		Lifetime: Singleton,
		Factory: func(*Container) (interface{}, error) {
			return instance, nil
		},
	}
}

// Resolve resolves a service from the container
func (c *Container) Resolve(name string) (interface{}, error) {
	c.mutex.RLock()
	descriptor, exists := c.services[name]
	c.mutex.RUnlock()

	if !exists {
		return nil, fmt.Errorf("service '%s' not registered", name)
	}

	// Check if singleton instance already exists
	if descriptor.Lifetime == Singleton {
		c.mutex.RLock()
		instance, exists := c.instances[name]
		c.mutex.RUnlock()

		if exists {
			return instance, nil
		}
	}

	// Create new instance
	instance, err := descriptor.Factory(c)
	if err != nil {
		return nil, fmt.Errorf("failed to create service '%s': %w", name, err)
	}

	// Store singleton instance
	if descriptor.Lifetime == Singleton {
		c.mutex.Lock()
		c.instances[name] = instance
		c.mutex.Unlock()
	}

	return instance, nil
}

// MustResolve resolves a service or panics if it fails
func (c *Container) MustResolve(name string) interface{} {
	service, err := c.Resolve(name)
	if err != nil {
		panic(err)
	}
	return service
}

// Clear clears all instances (useful for testing)
func (c *Container) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.instances = make(map[string]interface{})
}

// Service name constants for type-safe resolution
const (
	ServiceSoundManager       = "SoundManager"
	ServiceSpriteManager      = "SpriteManager"
	ServiceInputHandler       = "InputHandler"
	ServiceCollisionManager   = "CollisionManager"
	ServiceEntityManager      = "EntityManager"
	ServiceCameraSystem       = "CameraSystem"
	ServiceLeaderboardManager = "LeaderboardManager"
	ServiceAchievementManager = "AchievementManager"
	ServiceProgressionManager = "ProgressionManager"
	ServiceSpawner            = "Spawner"
	ServiceHUD                = "HUD"
	ServiceStarfield          = "Starfield"
	ServiceMenu               = "Menu"
	ServiceInfoMenu           = "InfoMenu"
)
