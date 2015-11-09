package main

import "time"

// Visit describes an access to a web resource.
type Visit struct {
	Timestamp time.Time `json:"timestamp"` // Time the resource was accessed.
	Identity  string    `json:"identity"`  // Identity of the client.
	Uri       Resource  `json:"uri"`       // IDentifier of the resource.
}

// Resource describes the location of a resource.
type Resource struct {
	Protocol string  `json:"protocol"` // Protocol of the resource.
	Hostname string  `json:"hostname"` // Hostname of the resource.
	Port     *string `json:"port"`     // Port of the resource, optional.
	Path     *string `json:"path"`     // Path of the resource, optional.
}
