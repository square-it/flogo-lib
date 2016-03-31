package process

// Provider is the interface that describes an object
// that can provide process definitions from a URI
type Provider interface {

	// GetProcess retrieves the process definition for the specified URI
	GetProcess(processURI string) *Definition
}
