package main

import (
	"archive/zip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func main() {
	destination := "./"
	archive, err := zip.OpenReader("dynamic_loading.zip")
	if err != nil {
		log.Fatal(err)
	}
	defer archive.Close()

	for _, f := range archive.File {
		filePath := filepath.Join(destination, f.Name)
		fmt.Println("Extracting file", filePath)

		// TODO: confirm this would work with S3
		if f.FileInfo().IsDir() {
			fmt.Println("Creating directory")
			os.MkdirAll(filePath, os.ModePerm)
			continue
		}
		if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
			log.Fatal(err)
		}

		destFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			log.Fatal(err)
		}
		defer destFile.Close()

		archiveFile, err := f.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer archiveFile.Close()

		if _, err := io.Copy(destFile, archiveFile); err != nil {
			log.Fatal(err)
		}
	}
}
