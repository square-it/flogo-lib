package resources

import "errors"

// ResourceManager interface
type ResourceManager interface {

	// LoadResources tells the manager to load the specified resource set
	LoadResources(config *ResourceSetConfig)

	// GetResource get the resource that corresponds to the specified id
	GetResource(id string) interface{}
}

var managers = make(map[string]ResourceManager)

// RegisterManager registers a resource manager for the specified type
func RegisterManager(resourceType string, manager ResourceManager) error {

	_, exists := managers[resourceType]

	if exists {
		return errors.New("ResourceManager already registered for type: " + resourceType)
	}

	managers[resourceType] = manager
	return nil
}

// GetManager gets the manager for the specified resource type
func GetManager(resourceType string) ResourceManager {
	return managers[resourceType]
}

// Get gets the specified resource, URI format is "res://{type}/{id}"
func Get(uri string) interface {} {

	return nil
}