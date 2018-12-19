package platform

import "fmt"

// Status defines if a resource is active or inactive.
type Status string

const (
	// Active status means that the resource can be used.
	Active Status = "active"
	// Inactive status means that the resource cannot be used.
	Inactive Status = "inactive"
)

// Valid determines if a Status value matches the enum.
func (s Status) Valid() error {
	switch s {
	case Active, Inactive:
		return nil
	default:
		return fmt.Errorf("invalid status")
	}
}
