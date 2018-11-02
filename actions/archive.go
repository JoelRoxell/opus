package actions

import (
	"fmt"
	"os"

	"github.com/mholt/archiver"
)

// Archive archives the selected build.
func Archive(projectName string, tag string) error {
	archivePath := fmt.Sprintf("artifacts/%s/archive", projectName)
	target := fmt.Sprintf("%s/%s.tar.gz", archivePath, tag)

	err := os.MkdirAll(archivePath, os.ModePerm)

	if err != nil {
		return err
	}

	err = archiver.TarGz.Make(target, []string{fmt.Sprintf("artifacts/%s/%s", projectName, tag)})

	if err != nil {
		return err
	}

	return nil
}

// OpenArchive un-compresses the selected file to destination.
func OpenArchive(projectName string, tag string, output string) error {
	source := fmt.Sprintf("artifacts/%s/archive/%s.tar.gz", projectName, tag)

	err := archiver.TarGz.Open(source, output)

	if err != nil {
		return err
	}

	return nil
}
