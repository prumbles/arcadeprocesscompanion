package main

import (
	"arcadeprocesscompanion/joystick"
	"arcadeprocesscompanion/models"
	"arcadeprocesscompanion/proc"
	"fmt"

	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"
)

/**
Command-line arguments
arcadeprocesscompanion KEEP_ALIVE_MATCHER JOY_2_KEY_MAPPINGS

	KEEP_ALIVE_MATCHER: wait for process with name that contains this string to end before exiting this program
	JOY_2_KEY_MAPPINGS_FILE: file path of joy-to-key mappings JSON
*/

func main() {
	args := os.Args[1:]

	keepAliveMatcher := args[0]
	joy2KeyMappingsFilePath := args[1]

	content, err := ioutil.ReadFile(joy2KeyMappingsFilePath)

	if err != nil {
		log.Fatal(err)
	}

	var rootMappings models.RootMappings

	unmarshalErr := json.Unmarshal(content, &rootMappings)

	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}

	quitChan := make(chan int)
	updateProcessChan := make(chan proc.CheckProcessesResult)
	exitApp := false

	joystickManager := joystick.NewJoystickManager(&rootMappings)
	joystickManager.StartPolling()

	go checkProcesses(rootMappings.ProcessPriorityOrder, "-1", keepAliveMatcher, quitChan, updateProcessChan)
	for {
		select {
		case newProcess := <-updateProcessChan:
			fmt.Printf("New Processs %v\n", newProcess.PriorityProcessMatcher)
			joystickManager.SetProcessFilter(newProcess.PriorityProcessMatcher, newProcess.PID)

		case <-quitChan:
			joystickManager.StopPolling()
			exitApp = true
			break
		}

		if exitApp {
			break
		}
	}

	//Sleep for just a short time so any simulated keyboard events from the joypad don't get propagated to the terminal on exit
	time.Sleep(500 * time.Millisecond)
}

func checkProcesses(processPriorityMatchers []string, previousProcessMatcher string,
	keepAliveMatcher string, quitChan chan int, updateProcessChan chan proc.CheckProcessesResult) {
	result := proc.CheckProcesses(processPriorityMatchers, keepAliveMatcher)

	if !result.KeepAliveProcessFound {
		quitChan <- 0
		return
	}

	processMatcher := result.PriorityProcessMatcher

	if !result.PriorityProcessFound {
		processMatcher = ""
	}

	if previousProcessMatcher != processMatcher {
		updateProcessChan <- result
	}

	time.Sleep(5 * time.Second)
	go checkProcesses(processPriorityMatchers, processMatcher, keepAliveMatcher, quitChan, updateProcessChan)
}
