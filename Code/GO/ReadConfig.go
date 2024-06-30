package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/oauth2/jwt"
)

// Function to read the JSON file and unmarshal it into the jwt.Config struct
func ReadConfig(filename string) (*jwt.Config, error) {
	// Open the JSON file
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	// Read the file content
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	// Read the file content
	var data map[string]interface{}

	json.Unmarshal(content, &data)

	// Unmarshal the JSON data into the jwt.Config struct
	config := jwt.Config{
		Email:        data["client_email"].(string),
		PrivateKeyID: data["private_key_id"].(string),
		PrivateKey:   []byte(data["private_key"].(string)),
		TokenURL:     data["token_uri"].(string),
		Scopes: []string{
			"https://www.googleapis.com/auth/spreadsheets.readonly",
		},
	}

	return &config, nil
}