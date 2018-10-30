package actions

import (
	"archive/tar"
	"builder/repositories"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
)

// Builder struct for managing builds.
type Builder struct {
	ID    string
	Store *repositories.MongoDBDataStore
}

// CreateImage creates a new build and tries to build a container for the specified project.
func CreateImage(tag string) error {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))

	if err != nil {
		return err
	}

	files := make(map[string]string)

	err = filepath.Walk("./tmp",
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if file, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Println(err)
			} else {
				if file.IsDir() {
					log.Println(path + " is a dir, skipping.")
				} else {
					files[strings.Replace(path, "tmp/", "", 1)] = path

					log.Println(file.Name(), path, info.Size())
				}
			}

			return nil
		})

	if err != nil {
		return err
	}

	tarBuffer := new(bytes.Buffer)
	tarWriter := tar.NewWriter(tarBuffer)

	defer tarWriter.Close()

	for fileName, path := range files {
		log.Printf("\x1b[0;93madding file to tar-context file: %s at: %s\x1b[0m", fileName, path)
		file, err := os.Open(path)

		if err != nil {
			return err
		}

		fileBytes, err := ioutil.ReadAll(file)

		if err != nil {
			return err
		}

		err = tarWriter.WriteHeader(&tar.Header{
			Name: fileName,
			Size: int64(len(fileBytes))})

		if err != nil {
			return err
		}

		_, err = tarWriter.Write(fileBytes)

		if err != nil {
			return err
		}
	}

	tarContextReader := bytes.NewReader(tarBuffer.Bytes())

	log.Println("\x1b[0;32mfinished building tar-context\x1b[0m")

	res, err := cli.ImageBuild(context.Background(), tarContextReader, types.ImageBuildOptions{
		Dockerfile: "Dockerfile",
		Tags:       []string{tag}})

	if err != nil {
		log.Println(err, " :unable to build docker image")
	} else {
		defer res.Body.Close()

		_, err = io.Copy(os.Stdout, res.Body)

		if err != nil {
			return errors.New(err.Error() + " :unable to read image build response")
		}
	}

	log.Println("\x1b[0;32mSuccessfully built image\x1b[0m")

	return nil
}
