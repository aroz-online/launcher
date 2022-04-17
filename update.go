package main

import (
	"fmt"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

func updateIfExists(binaryName string) {
	if fileExists("./updates") && fileExists("./updates/web/") && fileExists("./updates/system") {
		//All component exists. Update it
		newArozBinary, err := getUpdateBinaryFilename()
		if err != nil {
			fmt.Println("[LAUNCHER] Unable to access update files: ", err.Error())
		} else {
			//Binary file got. Update it
			//Backup the current executables and system files
			fmt.Println("[LAUNCHER] Starting system backup process (to ./arozos.old)")
			os.MkdirAll("./arozos.old", 0775)
			copy(binaryName, filepath.Join("./arozos.old", filepath.Base(binaryName)))

			//Backup the start.sh script and start.bat script if exists
			if fileExists("./start.sh") {
				copy("./start.sh", "./arozos.old/start.sh")
			}

			if fileExists("./start.bat") {
				copy("./start.bat", "./arozos.old/start.bat")
			}

			//Backup the important system files
			fmt.Println("[LAUNCHER] Backing up system files")
			cp.Copy("./system", "./arozos.old/system/")

			fmt.Println("[LAUNCHER] Backing up web server")
			cp.Copy("./web", "./arozos.old/web/")

			//Success. Continue binary replacement
			fmt.Println("[LAUNCHER] Copying updates to runtime environment")
			copy(newArozBinary, binaryName)

			cp.Copy("./updates/system", "./system/")
			cp.Copy("./updates/web", "./web/")

			//Restore the configs from the arozos.old
			fmt.Println("[LAUNCHER] Restoring previous configurations")
			restoreConfigs()

			fmt.Println("[LAUNCHER] Update Completed. Removing the update files")
			os.RemoveAll("./updates/")
		}

	} else if fileExists("./updates") && (!fileExists("./updates/web/") || !fileExists("./updates/system")) {
		//Update folder exists but some components is broken
		fmt.Println("[LAUNCHER] Detected damaged / incomplete update package. Skipping update process")
	}
}
