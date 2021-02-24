package main

import (
	. "arcadeprocesscompanion/joystick"
	"arcadeprocesscompanion/models"
	. "arcadeprocesscompanion/utils"

	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/micmonay/keybd_event"
)

/**
Command-line arguments
arcadeprocesscompanion PROCESS_NAME_LIKE JOY_2_KEY_MAPPINGS

	PROCESS_NAME_LIKE: wait for process with name that contains this string to end before exiting this program
	JOY_2_KEY_MAPPINGS_FILE: file path of joy-to-key mappings JSON
		[
			{
				Id: 1,
				Mappings: [
					{
						Buttons: [0],
						Key: "ESC",
						Shift: false,
						Ctrl: false,
						Alt: false
					}
				]
			}
		]
*/
func main() {
	args := os.Args[1:]

	// processNameLike := args[0]
	joy2KeyMappingsFilePath := args[1]

	content, err := ioutil.ReadFile(joy2KeyMappingsFilePath)

	if err != nil {
		log.Fatal(err)
	}

	var controllerMappings []models.ControllerMappings

	unmarshalErr := json.Unmarshal(content, &controllerMappings)

	if unmarshalErr != nil {
		log.Fatal(unmarshalErr)
	}

	joystickReaders := make([]JoystickReader, len(controllerMappings))

	kb, err := keybd_event.NewKeyBonding()
	if err != nil {
		panic(err)
	}

	//apparently you have to do this for the keybd_event package
	time.Sleep(2 * time.Second)

	for i := range controllerMappings {
		controllerMapping := &controllerMappings[i]

		for j := range controllerMapping.Mappings {
			mapping := &controllerMapping.Mappings[j]
			vk, err := GetVKCode(mapping.Key)

			if err != nil {
				log.Fatal(err)
			}

			mapping.VKKeyCode = vk
		}

		joystickReader, err := NewJoystickReader(*controllerMapping, &kb)

		if err != nil {
			log.Fatal(err)
		}

		joystickReaders[i] = joystickReader
	}

	watiTime := 50 * time.Millisecond
	for {
		time.Sleep(watiTime)

		for i := range joystickReaders {
			joystickReader := &joystickReaders[i]

			joystickReader.ProcessState()
		}
	}

}
