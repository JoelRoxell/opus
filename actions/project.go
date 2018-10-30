package actions

import (
	"builder/models"
	"builder/repositories"
	"errors"
)

// CreateProject creates a new store for project `name`.
func CreateProject(name string, remote string, session *repositories.MongoDBDataStore) error {
	db := session.Copy()
	defer db.Close()

	session.CreateProject(name, remote)

	return nil
}

func GetProject(name string, session *repositories.MongoDBDataStore) (models.Project, error) {
	db := session.Copy()
	defer db.Close()

	return session.GetProject(name)
}

// DeleteProject removed project by name.
func DeleteProject(name string) error {
	return errors.New("Not implemented")
}
