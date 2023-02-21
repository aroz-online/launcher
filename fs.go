package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	cp "github.com/otiai10/copy"
)

func restoreConfigs() {
	restoreIfExists("system/bridge.json")
	restoreIfExists("system/dev.uuid")
	restoreIfExists("system/cron.json")
	restoreIfExists("system/storage.json")
	restoreIfExists("web/SystemAO/vendor/")

	//Restore start script
	if fileExists("./arozos.old/start.sh") {
		copy("./arozos.old/start.sh", "./start.sh")
	}
	if fileExists("./arozos.old/start.bat") {
		copy("./arozos.old/start.bat", "./start.bat")
	}
}

func restoreOldArozOS() {
	fmt.Println("[LAUNCHER] ArozOS unable to launch. Restoring from backup")
	if fileExists("arozos.old") {
		backupfiles, err := filepath.Glob("arozos.old/*")
		if err != nil {
			fmt.Println("[LAUNCHER] Unable to restore backup. Exiting.")
			os.Exit(1)
		}

		for _, thisBackupFile := range backupfiles {
			if isDir(thisBackupFile) {
				cp.Copy(thisBackupFile, "./"+filepath.Base(thisBackupFile))
			} else {
				copy(thisBackupFile, "./"+filepath.Base(thisBackupFile))
			}
		}
	} else {
		fmt.Println("[LAUNCHER] ArozOS backup not found. Exiting.")
		os.Exit(1)
	}

}

func restoreIfExists(fileRelPath string) {
	if fileExists(filepath.Join("arozos.old", fileRelPath)) {
		if !isDir(filepath.Join("arozos.old", fileRelPath)) {
			copy(filepath.Join("arozos.old", fileRelPath), fileRelPath)
		} else {
			cp.Copy(filepath.Join("arozos.old", fileRelPath), fileRelPath)
		}
	}
}

// Auto detect and execute the correct binary
func autoDetectExecutable() string {
	if runtime.GOOS == "windows" {
		if fileExists("arozos.exe") {
			return "arozos.exe"
		}
	} else {
		if fileExists("arozos") {
			return "./arozos"
		}
	}

	//Not build from source. Look for release binary names
	binaryExecPath := "arozos_" + runtime.GOOS + "_" + runtime.GOARCH
	if runtime.GOOS == "windows" {
		binaryExecPath += ".exe"
	}

	binaryExecPath = "./" + binaryExecPath

	if fileExists(binaryExecPath) {
		return binaryExecPath
	} else {
		//Binary not found. Try download it
		err := downloadReleaseBinary()
		if err != nil {
			fmt.Println("[LAUNCHER] Unable to detect ArozOS start binary")
			os.Exit(1)
			return ""
		}

		if runtime.GOOS == "windows" {
			return "./arozos.exe"
		} else {
			return "./arozos"
		}

	}
}

func getUpdateBinaryFilename() (string, error) {
	updateFiles, err := filepath.Glob("./updates/*")
	if err != nil {
		return "", err
	}

	for _, thisFile := range updateFiles {
		if !isDir(thisFile) && filepath.Ext(thisFile) != ".gz" {
			//This might be the file
			return thisFile, nil
		}
	}

	return "", errors.New("file not found")
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, errors.New("invalid file")
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func isDir(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func fileExists(name string) bool {
	_, err := os.Stat(name)
	if err == nil {
		return true
	}
	if errors.Is(err, os.ErrNotExist) {
		return false
	}
	return false
}
