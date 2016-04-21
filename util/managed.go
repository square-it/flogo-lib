package util

import "fmt"

// Managed is an interface that is implemented by an object that needs to be
// managed via start/stop
type Managed interface {

	// Start starts the managed object
	Start() error

	// Stop stops the manged object
	Stop()
}

// StopManaged stops a "Managed" object
func StopManaged(managed Managed) error {

	defer func() error {
		if r := recover(); r != nil {
			err, ok := r.(error)
			if !ok {
				err = fmt.Errorf("%v", r)
			}

			return err
		}

		return nil
	}()

	managed.Stop()

	return nil
}
