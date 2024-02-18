package main

import (
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
)

const DATA_DIR = "./data/"
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

	startServer()
}

func startServer() {
	mux := http.NewServeMux()

	// Start the server
	mux.HandleFunc("/manifest", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Read the manifest file
		file, err := os.ReadFile(MANIFEST_PATH)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.Write(file)
	})

	mux.HandleFunc("/data/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Handling request for", r.URL.Path)
		if r.Method != "GET" {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		filePath := strings.TrimPrefix(r.URL.Path, "/data/")

		// Read the file from the data directory
		file, err := os.Open(DATA_DIR + filePath)
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		defer file.Close()

		w.Header().Set("Content-Type", "application/octet-stream")

		b := make([]byte, 4096)
		for {
			n, err := file.Read(b)
			if err != nil {
				if err == io.EOF {
					break
				}
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}

			w.Write(b[:n])
		}
	})

	http.ListenAndServe(":8080", mux)
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

	bytes := make([]byte, 500000)
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
		if !info.IsDir() {
			// Generate a hash and add the file to the manifest
			relPath, _ := filepath.Rel(DATA_DIR, path)
			hash := generateHash(path)
			files = append(files, File{relPath, hash})
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
