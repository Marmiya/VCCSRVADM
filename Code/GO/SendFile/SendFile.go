package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"golang.org/x/crypto/ssh"
)

func SendFile(ip_string, username_string, file_path string) (error) {
// Use SSH key authentication from the auth package
	// we ignore the host key in this example, please change this if you use this library
	clientConfig, _ := auth.PrivateKey(username_string, "/root/.ssh/id_rsa", ssh.InsecureIgnoreHostKey())

	// For other authentication methods see ssh.ClientConfig and ssh.AuthMethod

	// Create a new SCP client
	client := scp.NewClient(ip_string+":22", &clientConfig)

	// Connect to the remote server
	err := client.Connect()
	if err != nil {
		return  fmt.Errorf("error connecting to remote server: %v", err)
	}

	// Open a file
	f, err := os.Open(file_path)
	if err != nil {
		return fmt.Errorf("error opening file: %v", err)
	}

	// Close client connection after the file has been copied
	defer client.Close()

	// Close the file after it has been copied
	defer f.Close()

	// Finally, copy the file over
	// Usage: CopyFromFile(context, file, remotePath, permission)

        // the context can be adjusted to provide time-outs or inherit from other contexts if this is embedded in a larger application.
	err = client.CopyFromFile(context.Background(), *f, "/home/" + username_string + "/" + filepath.Base(file_path), "0655")

	if err != nil {
		return fmt.Errorf("error copying file: %v", err)
	}

	return nil
}