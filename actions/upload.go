package actions

import (
	"builder/integrations"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func UploadArchive(projectName string, tag string, bucket string) error {
	aws := &integrations.AwsService{Region: "eu-west-2"} // FIXME: set region by env.
	aws.Init()

	targetBackupFile := fmt.Sprintf("archive/%s.tar.gz", tag)
	archivePath := fmt.Sprintf("artifacts/%s/%s", projectName, targetBackupFile)

	err := aws.UploadFile(archivePath, targetBackupFile, bucket)

	if err != nil {
		return err
	}

	return nil
}

// Upload uploads all artifact contents to a specified bucket.
func Upload(projectName string, tag string, bucket string) error {
	aws := &integrations.AwsService{Region: "eu-west-2"} // FIXME: set region by env.
	aws.Init()

	files := make(map[string]string)
	artifactPath := fmt.Sprintf("artifacts/%s/%s", projectName, tag)

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

	for targetPath, filePath := range files {
		if err := aws.UploadFile(filePath, targetPath, bucket); err != nil {
			panic(err)
		}
	}

	return nil
}
