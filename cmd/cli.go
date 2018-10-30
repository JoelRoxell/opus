package main

import (
	"bufio"
	"builder/actions"
	"builder/models"
	"builder/repositories"
	"log"
	"os"
	"strings"

	"github.com/google/uuid"
)

func main() {
	session, _ := repositories.NewMongoDBConnection("localhost")

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()
		res := strings.Split(line, " ")

		options := res[1:]

		switch res[0] {

		case "list":
			err := actions.List(&options, session)

			if err != nil {
				log.Println(err)
			}

		case "project":
			switch options[0] {
			case "add":
				actions.CreateProject(options[1], options[2], session)

			case "remove":
				log.Println("Not implemented")

			default:
				log.Println("project [add, remove]")
			}

		case "build":
			if len(options) < 2 {
				log.Println("build supports only: build <project-name> <tag>")
				continue
			}

			project, err := actions.GetProject(options[0], session)

			if err != nil {
				log.Println("project not found; " + err.Error())

				continue
			}

			tag := project.Name + ":" + options[1]
			buildID := uuid.New().String()

			err = session.CreateBuild(project.Name, &models.Build{Tag: tag, ID: buildID, Status: models.STARTED})

			if err != nil {
				log.Println(err)
			}

			if err := actions.CreateImage(tag); err != nil {
				log.Println(err, "\x1b[0;91mfailed to build container\x1b[0m")
			}

			session.UpdateBuild(buildID, models.COMPLETED)

		case "exit":
			os.Exit(0)

		case "ls":

		default:
			log.Print("Unsupported command")
		}
	}
}
