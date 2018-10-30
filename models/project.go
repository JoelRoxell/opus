package models

// Project is a project, which holds a remote source which is used to fetch source code. The project also contains a history of tried builds.
type Project struct {
	Name   string
	Source string
	Builds []Build
}

// Build represents a single tried build.
type Build struct {
	ID     string
	Tag    string
	Status string
}

// Represents the different states a Build can have.
const (
	STARTED   = "STARTED"
	PENDING   = "PENDING"
	FAILED    = "FAILED"
	COMPLETED = "COMPLETED"
)
