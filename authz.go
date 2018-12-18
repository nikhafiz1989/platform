package platform

import (
	"errors"
	"fmt"
)

var (
	// ErrAuthorizerNotSupported notes that the provided authorizer is not supported for the action you are trying to perform.
	ErrAuthorizerNotSupported = errors.New("your authorizer is not supported, please use *platform.Authorization as authorizer")
	// ErrInvalidPermissionResource notes that the provided resource is invalid
	ErrInvalidPermissionResource = errors.New("unknown resource for permission")
	// ErrInvalidPermissionAction notes that the provided action is invalid
	ErrInvalidPermissionAction = errors.New("unknown action for permission")
)

// Authorizer will authorize a permission.
type Authorizer interface {
	// Allowed returns true is the associated permission is allowed by the authorizer
	Allowed(p Permission) bool

	// ID returns an identifier used for auditing.
	Identifier() ID

	// GetUserID returns the user id.
	GetUserID() ID

	// Kind metadata for auditing.
	Kind() string
}

func allowed(p Permission, ps []Permission) bool {
	for _, perm := range ps {
		if perm.Action == p.Action && perm.Resource == p.Resource {
			return true
		}
	}
	return false
}

// Action is an enum defining all possible resource operations
type Action string

const (
	// ReadAction is the action for reading.
	ReadAction Action = "read"
	// WriteAction is the action for writing.
	WriteAction Action = "write"
	// CreateAction is the action for creating new resources.
	CreateAction Action = "create"
	// DeleteAction is the action for deleting an existing resource.
	DeleteAction Action = "delete"
)

var actions = []Action{
	ReadAction,
	WriteAction,
	CreateAction,
	DeleteAction,
}

// Resource is an enum defining all resources that have a permission model in platform
type Resource string

const (
	// UsersResource gives permissions to one or more users.
	UsersResource = Resource("users")
	// OrgsResource gives permissions to one or more orgs.
	OrgsResource = Resource("orgs")
	// TasksResource gives permissions to one or more tasks.
	TasksResource = Resource("tasks")
	// BucketsResource gives permissions to one or more buckets.
	BucketsResource = Resource("buckets")
	// DashboardsResource gives permissions to one or more dashboards.
	DashboardsResource = Resource("dashboards")
	// SourcesResource gives permissions to one or more sources.
	SourcesResource = Resource("sources")
)

var resources = []Resource{
	UsersResource,
	OrgsResource,
	TasksResource,
	BucketsResource,
	DashboardsResource,
	SourcesResource,
}

// Permission defines an action and a resource.
type Permission struct {
	Action   Action   `json:"action"`
	Resource Resource `json:"resource"`
	ID       *ID      `json:"id,omitempty"`
	Name     *string  `json:"name,omitempty"`
}

func (p Permission) String() string {
	str := fmt.Sprintf("%s:%s", p.Action, p.Resource)
	if p.Name != nil && *p.Name != "" {
		str += fmt.Sprintf(":%s", *p.Name)
	}

	return str
}

// Valid checks if there the resource and action provided is known
func (p *Permission) Valid() error {
	switch p.Resource {
	case UsersResource, OrgsResource, TasksResource, BucketsResource, DashboardsResource, SourcesResource:
	default:
		return ErrInvalidPermissionResource
	}

	switch p.Action {
	case ReadAction, WriteAction, DeleteAction, CreateAction:
	default:
		return ErrInvalidPermissionAction
	}

	if p.ID != nil && !(*p.ID).Valid() {
		return ErrInvalidID
	}

	return nil
}

// NewPermission returns a permission with provided arguments
func NewPermission(a Action, r Resource) (*Permission, error) {
	p := &Permission{
		Action:   a,
		Resource: r,
	}

	return p, p.Valid()
}

// NewPermissionAtID creates a permission with the provided arguments
func NewPermissionAtID(id ID, name string, a Action, r Resource) (*Permission, error) {
	p := &Permission{
		Action:   a,
		Resource: r,
		ID:       &id,
		Name:     &name,
	}

	return p, p.Valid()
}

// OperPermissions are the default permissions for those who setup the application
func OperPermissions() []Permission {
	ps := []Permission{}
	for _, r := range resources {
		for _, a := range actions {
			ps = append(ps, Permission{Action: a, Resource: r})
		}
	}

	return ps
}
