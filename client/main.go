package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"strings"
)

const DATA_DIR = "./data/"

func main() {
	req, err := http.NewRequest("GET", "http://localhost:8080/manifest", nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Status code was not 200")
		return
	}

	defer resp.Body.Close()

	manifest, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(string(manifest))

	b := make([]byte, 0)
	for _, byte := range manifest {
		if byte != '\n' {
			b = append(b, byte)
		} else {
			line := string(b)
			fileName := strings.Split(line, " ")[0]
			fmt.Println("Downloading", fileName)

			req, err := http.NewRequest("GET", "http://localhost:8080/data/"+fileName, nil)
			if err != nil {
				fmt.Println(err)
				return
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				fmt.Println(err)
				return
			}

			if resp.StatusCode != http.StatusOK {
				fmt.Println("Status code was not 200. Code:", resp.StatusCode)
				return
			}

			defer resp.Body.Close()

			file, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
				return
			}

			err = os.MkdirAll(path.Dir(DATA_DIR+fileName), 0770)
			if err != nil {
				fmt.Println(err)
				return
			}

			savedFile, err := os.Create(DATA_DIR + fileName)
			if err != nil {
				fmt.Println(err)
				return
			}

			// write downloaded file to disk in chunks
			for {
				n, err := savedFile.Write(file)
				if err != nil {
					fmt.Println(err)
					return
				}

				if n == 0 {
					break
				}

				file = file[n:]
			}

			b = b[:0]
		}
	}
}
