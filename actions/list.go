package actions

import (
	"builder/repositories"
	"encoding/json"
	"errors"
	"fmt"
)

// List handles the different options for the list action.
func List(args *[]string, session *repositories.MongoDBDataStore) error {
	length := len(*args)

	if length != 1 {
		return errors.New("list command only supports 1 option; [all, progress, failed, completed]")
	}

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
