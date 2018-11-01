package main

import (
	"builder/actions"
	"builder/models"
	"builder/repositories"
	"builder/utils"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"gopkg.in/src-d/go-git.v4"
)

func main() {
	session, _ := repositories.NewMongoDBConnection("localhost")

	add := flag.String("add", "", "Adds a new project to the systems context, add <project-name> <source>")
	remove := flag.String("delete", "", "removes project by name, remove <project-name>")
	trigger := flag.String("trigger", "", "trigger <project-name> <tag> <bucket>")
	status := flag.String("status", "", "status <project-filter-by-name>")

	flag.Parse()

	options := flag.Args()

	if len(*add) > 0 {
		fmt.Println(*add)
		fmt.Println(options)
		actions.CreateProject(*add, options[0], session)
	} else if len(*remove) > 0 {
		fmt.Println(*remove)
		fmt.Println(options)
		utils.Print("Not implemented", utils.INFO)
		return
	} else if len(*status) > 0 {
		err := actions.List(session)

		if err != nil {
			utils.Print(err.Error(), utils.ERROR)
		}
	} else if len(*trigger) > 0 {
		projectName := *trigger
		tag := options[0]
		bucket := options[1]

		// TODO: each project context should exist in an own folder.
		project, err := actions.GetProject(projectName, session)
		imageTag := project.Name + ":" + tag
		buildID := uuid.New().String()

		sourcePath := fmt.Sprintf("jobs/%s/%s", projectName, tag)
		targetPath := fmt.Sprintf("artifacts/%s/%s", projectName, tag)

		sourceAbsPath, err := filepath.Abs(sourcePath)
		targetAbsPath, err := filepath.Abs(targetPath)

		err = session.CreateBuild(
			project.Name, &models.Build{
				Tag:    imageTag,
				ID:     buildID,
				Status: models.STARTED,
			})

		if err := os.MkdirAll(sourceAbsPath, os.ModePerm); err != nil {
			utils.Print(err.Error(), utils.ERROR)

			return
		}

		if err := os.MkdirAll(targetPath, os.ModePerm); err != nil {
			utils.Print(err.Error(), utils.ERROR)

			return
		}

		if err != nil {
			log.Println(err)
			return
		}

		if err := os.RemoveAll(sourceAbsPath); err != nil {
			log.Println(err)
		}

		if err != nil {
			utils.Print("project not found; "+err.Error(), utils.WARNING)
			return
		}

		git.PlainClone(sourcePath, false, &git.CloneOptions{
			URL:      project.Source,
			Progress: os.Stdout,
		})

		if err != nil {
			utils.Print(err.Error(), utils.ERROR)
		}

		if err := actions.CreateImage(tag, sourcePath); err != nil {
			utils.Print("failed to build container: "+err.Error(), utils.ERROR)
		}

		session.UpdateBuild(buildID, models.COMPLETED)

		if err := actions.CreateAndStart(tag, targetAbsPath); err != nil {
			utils.Print(err.Error(), utils.ERROR)
		}

		if err := actions.Archive(projectName, tag); err != nil {
			utils.Print(err.Error(), utils.WARNING)
		}

		actions.Upload(projectName, tag, bucket)
	}

	return
}
