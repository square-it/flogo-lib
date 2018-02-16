package resources

import (
	"errors"
	"strings"
	"fmt"
)

// ResourceManager interface
type ResourceManager interface {
	// LoadResources tells the manager to load the specified resource set
	LoadResources(config *ResourceConfig)

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

// Load specified resource into its corresponding Resource Manager
func Load(config *ResourceConfig) {

}

// Get gets the specified resource, URI format is "res://{type}:{id}"
func Get(uri string) (interface{}, error) {


	return nil, nil
}

func GetType(uri string) (string, error) {

	if !strings.HasPrefix(uri, "res://") {

		return "", errors.New("Invalid resource uri: " + uri)
	}

	fmt.Println(uri[6:])

	idx := strings.Index(uri[6:], ":")

	if idx < 0 {
		return "", errors.New("Invalid resource uri: " + uri)
	}

	return uri[6:6+idx], nil
}
