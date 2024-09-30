package service

import "fmt"

type Service interface {
	Initialize(config map[string]interface{}) error
	Start() error
	Stop() error
	Status() string
}

type Manager struct {
	services map[string]Service
}

// RegisterService Register a new service
func (sm *Manager) RegisterService(name string, service Service) {
	sm.services[name] = service
}

// EnableService Enable and start a service
func (sm *Manager) EnableService(name string) error {
	service, exists := sm.services[name]
	if !exists {
		return fmt.Errorf("service %s not found", name)
	}
	if err := service.Initialize(nil); err != nil {
		return err
	}
	return service.Start()
}

// DisableService Disable and stop a service
func (sm *Manager) DisableService(name string) error {
	service, exists := sm.services[name]
	if !exists {
		return fmt.Errorf("service %s not found", name)
	}
	return service.Stop()
}

// GetStatus Get service status
func (sm *Manager) GetStatus(name string) string {
	service, exists := sm.services[name]
	if !exists {
		return "Unknown service"
	}
	return service.Status()
}
