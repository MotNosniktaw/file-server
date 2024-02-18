package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
)

const DATA_DIR = "./data"
const MANIFEST_PATH = "./manifest.txt"

type File struct {
	Name string
	Hash string
}

func main() {
	// Create the data directory and files
	createData()

	// Generate the file manifest
	generateFileManifest()
}

func createData() {
	fmt.Println("Creating data")
	createFile("data/top.txt")
	createFile("data/middle/middle.txt")
	createFile("data/middle/bottom/bottom.txt")

}

func createFile(name string) error {
	fmt.Println("Creating file")

	os.MkdirAll(path.Dir(name), 0770)

	file, err := os.Create(name)
	if err != nil {
		log.Fatalln(err)
	}

	bytes := make([]byte, 1024)
	rand.Read(bytes)
	file.Write(bytes)

	file.Close()
	return nil
}

func generateFileManifest() {
	fmt.Println("Generating file manifest")

	// Walk the data directory
	// For each file, generate a hash and add it to the manifest
	// For each directory, add it to the manifest
	files := checkDirectory(DATA_DIR)

	// Write the manifest to a file
	manifestFile, err := os.Create(MANIFEST_PATH)
	if err != nil {
		log.Fatalln(err)
	}
	defer manifestFile.Close()

	for _, file := range files {
		manifestFile.WriteString(file.Name + " " + file.Hash + "\n")
	}
}

func checkDirectory(dir string) []File {
	files := []File{}

	// For each file, check if it is a directory or a file
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
		} else {
			// Generate a hash and add the file to the manifest
			hash := generateHash(path)
			files = append(files, File{path, hash})
		}
		return nil
	})

	return files
}

func generateHash(filePath string) string {
	hasher := sha256.New()
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	// Copy the file into the hasher
	if _, err := file.WriteTo(hasher); err != nil {
		log.Fatalln(err)
	}

	// Generate the hash
	hash := fmt.Sprintf("%x", hasher.Sum(nil))

	return hash
}
