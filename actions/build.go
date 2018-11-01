package actions

import (
	"archive/tar"
	"builder/repositories"
	"builder/utils"
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
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/client"
)

// Builder struct for managing builds.
type Builder struct {
	ID    string
	Store *repositories.MongoDBDataStore
}

// CreateImage creates a new build and tries to build a container for the specified project.
func CreateImage(tag string, sourcePath string) error {
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))

	if err != nil {
		return err
	}

	files := make(map[string]string)

	err = filepath.Walk(sourcePath,
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
					key := strings.Replace(path, sourcePath+"/", "", 1)
					files[key] = path

					log.Println(path, sourcePath)

					log.Println(file.Name(), key, info.Size())
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

	utils.Print("finished building tar-context", utils.NORMAL)

	res, err := cli.ImageBuild(
		context.Background(),
		tarContextReader,
		types.ImageBuildOptions{
			Dockerfile: "Dockerfile",
			Tags:       []string{tag},
		},
	)

	if err != nil {
		return err
	}

	defer res.Body.Close()

	_, err = io.Copy(os.Stdout, res.Body)

	if err != nil {
		return errors.New(err.Error() + " :unable to read image build response")
	}

	log.Println("\x1b[0;32msuccessfully built image\x1b[0m")

	return nil
}

// CreateAndStart runs a built image and collects the produced artifacts.
func CreateAndStart(image string, targetAbsPath string) error {
	// FIXME: Read cli version from ENV.
	cli, err := client.NewClientWithOpts(client.WithVersion("1.37"))

	if err != nil {
		return err
	}

	utils.Print("artifacts will be built to "+targetAbsPath, utils.INFO)

	res, err := cli.ContainerCreate(
		context.Background(),
		&container.Config{
			Image: image,
			// TODO: Override cmd via project yaml and/or environment variables.
			// Cmd:   []string{"touch", "/opt/app/build/bin.dat"},
			Tty: true,
		},
		&container.HostConfig{
			Mounts: []mount.Mount{
				{
					Type:   mount.TypeBind,
					Source: targetAbsPath,
					Target: "/usr/src/app/build",
				},
			},
		},
		nil,
		"",
	)
	ctx := context.Background()
	err = cli.ContainerStart(
		ctx,
		res.ID,
		types.ContainerStartOptions{},
	)
	statusCh, errCh := cli.ContainerWait(
		ctx,
		res.ID,
		container.WaitConditionNotRunning,
	)

	// Block and wait for container to finish the internal build.
	select {
	case err := <-errCh:
		return err
	case <-statusCh:
	}

	utils.Print(fmt.Sprintf("%s finished successfully", res.ID), utils.NORMAL)

	return nil
}
