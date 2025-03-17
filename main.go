package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

type ModrinthIndex struct {
	Files []FileEntry `json:"files"`
}

type FileEntry struct {
	Path      string   `json:"path"`
	Downloads []string `json:"downloads"`
}

// unmarshal the json file into a struct
func getModrinthIndex(filePath string) (*ModrinthIndex, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, errors.New(color.RedString("[ERR] File not found: %s", filePath))
	}
	defer file.Close()

	var index ModrinthIndex
	if err := json.NewDecoder(file).Decode(&index); err != nil {
		return nil, errors.New(color.RedString("[ERR] Invalid JSON format in %s", filePath))
	}
	return &index, nil
}


// download the files from the modrinth index
func downloadFiles(toSave map[string]FileEntry) {
	total := len(toSave)
	saved := 0

	for name, entry := range toSave {
		url := entry.Downloads[0]
		fpath := entry.Path

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(color.RedString("[ERR] Failed to fetch %s: %s", name, err))
			continue
		}

		var codeColored string
		if resp.StatusCode == 200 {
			codeColored = color.GreenString("%d", resp.StatusCode)
		} else {
			codeColored = color.RedString("%d", resp.StatusCode)
		}

		fmt.Printf("[%s] Saving %s\n",
			codeColored, color.YellowString(name),
		)

		outputFilePath := path.Join("returns", fpath)
		if err := os.MkdirAll(filepath.Dir(outputFilePath), os.ModePerm); err != nil {
			fmt.Println(color.RedString("[ERR] Failed to create directories for %s: %s", name, err))
			continue
		}

		outFile, err := os.Create(outputFilePath)
		if err != nil {
			fmt.Println(color.RedString("[ERR] Failed to create file %s: %s", name, err))
			continue
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, resp.Body); err != nil {
			fmt.Println(color.RedString("[ERR] Failed to save %s: %s", name, err))
			continue
		}
		saved++
	}

	fmt.Println(color.MagentaString("Saved %d/%d files", saved, total))
}

func main() {
	defaultPath := "./modrinth.index.json"
	var filePath string

	// try creating the file if it doesn't exist
	file, err := os.OpenFile(defaultPath, os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		if !os.IsExist(err) {
			panic(err)
		}
	} else {
		file.Close()
	}

	fmt.Printf("Enter the path to the %s file\n", color.GreenString("modrinth.index.json"))
	fmt.Printf("Or press Enter to use default %s:\n", color.YellowString(defaultPath))
	fmt.Print(color.GreenString("> "))

	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)

	if input == "" {
		filePath = defaultPath
	} else {
		filePath = input
	}

	mrIndex, err := getModrinthIndex(filePath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	toSave := make(map[string]FileEntry)
	for _, entry := range mrIndex.Files {
		if len(entry.Downloads) > 0 {
			toSave[filepath.Base(entry.Path)] = entry
		}
	}

	downloadFiles(toSave)
}
