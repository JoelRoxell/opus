package actions

import (
	"builder/repositories"
	"encoding/json"
	"fmt"
)

// List handles the different options for the list action.
func List(session *repositories.MongoDBDataStore) error {
	projects, _ := session.GetProjects()
	count := len(projects)

	fmt.Printf("Project count: %d \n", count)

	for i := range projects {
		project := projects[i]
		str, _ := json.MarshalIndent(project, "", "    ")

		fmt.Println(string(str))
	}

	return nil
}
