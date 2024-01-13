package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Message string `yaml:"message"`
}

var config Config

func main() {
	config, err := loadConfig("config")
	if err != nil {
		fmt.Println("Error loading configuration:", err)
		return
	}

	http.HandleFunc("/api/hello", helloHandler)

	port := 8080
	fmt.Printf("Server is listening on port %d...\n", port)
	err = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	if err != nil {
		fmt.Println("Error starting the server:", err)
	}
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}

	var requestData map[string]interface{}
	err = json.Unmarshal(body, &requestData)
	if err != nil {
		http.Error(w, "Error decoding JSON", http.StatusBadRequest)
		return
	}

	// Access data from the JSON request, if needed
	// For example: requestData["key"]

	// Respond with the configured message
	fmt.Fprint(w, config.Message)
}

func loadConfig(directory string) (Config, error) {
	var loadedConfig Config

	files, err := ioutil.ReadDir(directory)
	if err != nil {
		return loadedConfig, fmt.Errorf("error reading directory: %v", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yaml" {
			filePath := filepath.Join(directory, file.Name())
			data, err := ioutil.ReadFile(filePath)
			if err != nil {
				return loadedConfig, fmt.Errorf("error reading file: %v", err)
			}

			err = yaml.Unmarshal(data, &loadedConfig)
			if err != nil {
				return loadedConfig, fmt.Errorf("error unmarshalling YAML: %v", err)
			}

			fmt.Printf("Configuration loaded from %s\n", filePath)
		}
	}

	return loadedConfig, nil
}