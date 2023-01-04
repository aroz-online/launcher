package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

/*
	Download.go

	This script help download the latest release from Github, if endpoint is not specified
*/

func downloadReleaseBinary() error {
	//Check if the src.txt file exists and overwrite default release if found
	src := "https://github.com/tobychui/arozos/releases/latest/download/"
	srcFileContent, err := os.ReadFile("src.txt")
	if err == nil {
		src = strings.TrimSpace(string(srcFileContent))
	}

	if src[len(src)-1:] != "/" {
		src = src + "/"
	}

	fmt.Println("[LAUNCHER] Downloading release from " + src)

	//Download binary
	binaryTag := "arozos_" + runtime.GOOS + "_" + runtime.GOARCH
	if runtime.GOOS == "windows" {
		binaryTag += ".exe"
	}

	binaryName := "arozos"
	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	fmt.Println("[LAUNCHER] Downloading arozos release binary from: " + src + binaryTag)
	err = downloadFile(binaryName, src+binaryTag)
	if err != nil {
		fmt.Println("[LAUNCHER] Trying to download binary from release page but failed: ", err.Error())
		return err
	}

	//Download web.tar.gz
	fmt.Println("[LAUNCHER] Downloading web.tar.gz from: " + src + "web.tar.gz (this might take a while)...")
	err = downloadFile("web.tar.gz", src+"web.tar.gz")
	if err != nil {
		fmt.Println("[LAUNCHER] Trying to download binary from release page but failed: ", err.Error())
		return err
	}

	return nil
}

func downloadFile(filepath string, url string) (err error) {

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
