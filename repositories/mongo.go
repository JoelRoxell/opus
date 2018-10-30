package repositories

import (
	"builder/models"
	"log"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDBDataStore holds a mongo-db session.
type MongoDBDataStore struct {
	*mgo.Session
}

const projectCollection = "projects"

// NewMongoDBConnection creates a new session to the DB.
func NewMongoDBConnection(url string) (*MongoDBDataStore, error) {
	session, err := mgo.Dial(url)

	if err != nil {
		log.Fatalln(err)
	}

	return &MongoDBDataStore{Session: session}, nil
}

// GetProjects fetches all projects that have been added.
func (m *MongoDBDataStore) GetProjects() ([]models.Project, error) {
	session := m.Copy()

	defer session.Close()

	collection := session.DB("front-supervisor").C(projectCollection)

	var res []models.Project

	err := collection.Find(nil).All(&res)

	if err != nil {
		return nil, err
	}

	return res, nil
}

// CreateProject creates a new project in db.
func (m *MongoDBDataStore) CreateProject(name string, source string) error {
	session := m.Copy()
	defer session.Close()

	collection := session.DB("front-supervisor").C(projectCollection)

	project := &models.Project{Name: name, Source: source, Builds: []models.Build{}}

	err := collection.Insert(project)

	if err != nil {
		return err
	}

	return nil
}

// GetProject fetches project by name.
func (m *MongoDBDataStore) GetProject(name string) (models.Project, error) {
	session := m.Copy()
	defer session.Close()

	collection := session.DB("front-supervisor").C(projectCollection)

	var res models.Project

	err := collection.Find(bson.M{"name": name}).One(&res)

	return res, err
}

// CreateBuild adds a build to the db.
func (m *MongoDBDataStore) CreateBuild(name string, build *models.Build) error {
	session := m.Copy()

	defer session.Close()

	projects := session.DB("front-supervisor").C(projectCollection)

	err := projects.Update(bson.M{"name": name}, bson.M{"$push": bson.M{"builds": build}})

	if err != nil {
		return err
	}

	return nil
}

// UpdateBuild adds a build to the db.
func (m *MongoDBDataStore) UpdateBuild(buildID string, status string) error {
	session := m.Copy()

	defer session.Close()

	projects := session.DB("front-supervisor").C(projectCollection)

	// db.projects.update({"builds.id": "7b6b1cc2-b222-4175-adc4-d069eb8f6f67"}, {$set: { "builds.$.status": "TEST" } })

	err := projects.Update(bson.M{"builds.id": buildID}, bson.M{"$set": bson.M{"builds.$.status": status}})

	if err != nil {
		return err
	}

	return nil
}
