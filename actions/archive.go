package actions

import (
	"archive/tar"
	"builder/integrations"
	"builder/utils"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Archive archives the selected build.
func Archive(projectName string, tag string) error {
	aws := &integrations.AwsService{Region: "eu-west-2"} // FIXME: set region by env.
	aws.Init()
	files := make(map[string]string)
	artifactPath := fmt.Sprintf("artifacts/%s/%s", projectName, tag)
	backupLocation := fmt.Sprintf("artifacts/%s/archive", projectName)
	destinationTar := fmt.Sprintf("%s/%s.tar.gz", backupLocation, tag)

	err := filepath.Walk(artifactPath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if file, err := os.Stat(path); os.IsNotExist(err) {
				fmt.Println(err)
			} else {
				if !file.IsDir() {
					key := strings.Replace(path, artifactPath, "build", 1)
					files[key] = path

					fmt.Println(key)
				}
			}

			return nil
		})

	if err != nil {
		return err
	}

	tarfile, err := os.Create(destinationTar)
	defer tarfile.Close()

	var fileWriter io.WriteCloser = tarfile

	tarfileWriter := tar.NewWriter(fileWriter)

	defer tarfileWriter.Close()

	for _, filePath := range files {
		if err != nil {
			utils.Print(err.Error(), utils.ERROR)

			continue
		}

		log.Printf("adding %s to archive\n", filePath)

		file, err := os.Open(filePath)

		if err != nil {
			utils.Print(err.Error(), utils.ERROR)

			continue
		}

		defer file.Close()

		stats, err := file.Stat()

		if err != nil {
			log.Println(err)

			continue
		}

		header := new(tar.Header)
		header.Name = file.Name()
		header.Size = stats.Size()
		header.Mode = int64(stats.Mode())
		header.ModTime = stats.ModTime()

		err = tarfileWriter.WriteHeader(header)

		if err != nil {
			log.Println(err)
		}

		_, err = io.Copy(tarfileWriter, file)

		if err != nil {
			log.Println(err)
		}
	}

	return nil
}
