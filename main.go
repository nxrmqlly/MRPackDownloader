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

// Unmarshal the JSON file into a struct
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

// Download the files from the Modrinth index
func downloadFiles(toSave map[string]FileEntry) {
	total := len(toSave)
	saved := 0

	for name, entry := range toSave {
		if len(entry.Downloads) == 0 {
			fmt.Println(color.RedString("[ERR] No download URL for %s", name))
			continue
		}

		url := entry.Downloads[0]
		fpath := entry.Path

		resp, err := http.Get(url)
		if err != nil {
			fmt.Println(color.RedString("[ERR] Failed to fetch %s: %s", name, err))
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Printf("[%s] Failed to fetch %s\n", color.RedString("%d", resp.StatusCode), color.YellowString(name))
			continue
		}

		fmt.Printf("[%s] Saving %s\n", color.GreenString("%d", resp.StatusCode), color.YellowString(name))

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

		if _, err := io.Copy(outFile, resp.Body); err != nil {
			fmt.Println(color.RedString("[ERR] Failed to save %s: %s", name, err))
			outFile.Close()
			continue
		}
		outFile.Close()

		saved++
	}

	fmt.Println(color.MagentaString("Saved %d/%d files", saved, total))
}

func main() {
	defaultPath := "./modrinth.index.json"
	var filePath string

	// ensure the modrinth.index.json exists
	if _, err := os.Stat(defaultPath); os.IsNotExist(err) {
		_, err := os.Create(defaultPath)
		if err != nil {
			fmt.Println(color.RedString("[ERR] Failed to create default file: %s", err))
			os.Exit(1)
		}
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
