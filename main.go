package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

/*
	ArozOS Launcher
	For auto update and future extension purpose

	Author: tobychui

*/

const (
	launcherVersion   = "1.3"
	restoreRetryCount = 3 //The number of retry before restore old version, if not working after restoreRetryCount + 1 launcher will exit
)

var (
	norestart bool = false
)

func main() {
	//Print basic information
	fmt.Println("[LAUNCHER] ArozOS Launcher ver " + launcherVersion)
	binaryName := autoDetectExecutable()
	fmt.Println("[LAUNCHER] Choosing binary executable: " + binaryName)

	//Check if updates exists. If yes, overwrite it
	updateIfExists(binaryName)

	//Check launch paramter for norestart
	for _, arg := range os.Args[1:] {
		if arg == "-h" || arg == "-help" {
			//help argument, do not restart
			norestart = true
		} else if arg == "-version" || arg == "-v" {
			//version argument, no restart
			norestart = true
		}
	}

	//Register the binary start path
	cmd := exec.Command(binaryName, os.Args[1:]...)
	cmd.Dir = filepath.Dir(binaryName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	//Register the http server to notify ArozOS there is a launcher will handle the update
	go func() {
		http.HandleFunc("/chk", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("LauncherA v" + launcherVersion))
			fmt.Println("[LAUNCHER] CHK RECV - DONE")
		})

		http.ListenAndServe("127.0.0.1:25576", nil)
	}()

	retryCounter := 0
	//Start the cmd
	for {
		startTime := time.Now().Unix()
		err := cmd.Run()
		endTime := time.Now().Unix()

		if err != nil {
			panic(err)
		}
		if norestart {
			return
		}
		if endTime-startTime < 3 {
			//Less than 3 seconds, shd be crashed. Add to retry counter
			fmt.Println("[LAUNCHER] ArozOS Crashed. Restarting in 3 seconds... ")
			retryCounter++
		} else {
			fmt.Println("[LAUNCHER] ArozOS Exited. Restarting in 3 seconds... ")
			retryCounter = 0
		}

		time.Sleep(3 * time.Second)

		if retryCounter > restoreRetryCount+1 {
			//Fail to start. Exit program
			log.Fatal("Unable to start ArozOS. Exiting to OS")
			return
		} else if retryCounter > restoreRetryCount {
			//Restore from old version of the binary
			restoreOldArozOS()
		} else {
			updateIfExists(binaryName)
		}

		//Rebuild the start paramters
		cmd = exec.Command(binaryName, os.Args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

}
