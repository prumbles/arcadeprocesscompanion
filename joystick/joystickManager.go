package joystick

import (
	"arcadeprocesscompanion/models"
	"arcadeprocesscompanion/utils"
	"log"
	"sync"
	"time"
)

const FRAME_DURATION time.Duration = 20 * time.Millisecond

type JoystickManager struct {
	RootConfig *models.RootMappings

	filterUpdatingMutex      *sync.Mutex
	isPolling                bool
	activeControllerMappings []models.ControllerMappings
	activeJoystickReaders    []JoystickReader
	processNameMatchFilter   string
	exitPoll                 bool
}

func NewJoystickManager(rootConfig *models.RootMappings) JoystickManager {
	manager := JoystickManager{}

	manager.filterUpdatingMutex = &sync.Mutex{}

	manager.RootConfig = rootConfig

	manager.isPolling = false
	manager.activeControllerMappings = nil
	manager.processNameMatchFilter = ""
	manager.exitPoll = false

	//apparently you have to do this for the keybd_event package
	time.Sleep(2 * time.Second)

	return manager
}

func (joystickManager *JoystickManager) SetProcessFilter(processNameMatchFilter string) error {
	joystickManager.filterUpdatingMutex.Lock()
	defer joystickManager.filterUpdatingMutex.Unlock()

	//do caching stuff here

	var newMappings []models.ControllerMappings

	for i := range joystickManager.RootConfig.Mappings {
		mappings := joystickManager.RootConfig.Mappings[i]

		if mappings.Filters == nil {
			if processNameMatchFilter == "" {
				newMappings = append(newMappings, mappings)
			}
		} else {
			for j := range mappings.Filters.Processes {
				if mappings.Filters.Processes[j] == processNameMatchFilter {
					newMappings = append(newMappings, mappings)
					break
				}
			}
		}

	}

	joystickManager.activeControllerMappings = newMappings

	newJSReaders := make([]JoystickReader, len(newMappings))

	for i := range newMappings {
		controllerMapping := &newMappings[i]

		for j := range controllerMapping.Mappings {
			mapping := &controllerMapping.Mappings[j]
			if mapping.Key == nil {
				continue
			}

			vk, err := utils.GetVKCode(*mapping.Key)

			if err != nil {
				log.Fatal(err)
				return err
			}

			mapping.VKKeyCode = vk
		}

		joystickReader, err := NewJoystickReader(*controllerMapping)

		if err != nil {
			log.Fatal(err)
		}

		newJSReaders[i] = joystickReader
	}

	joystickManager.activeJoystickReaders = newJSReaders

	return nil
}

func (joystickManager *JoystickManager) StopPolling() {
	joystickManager.filterUpdatingMutex.Lock()
	defer joystickManager.filterUpdatingMutex.Unlock()

	joystickManager.exitPoll = true
}

func (joystickManager *JoystickManager) StartPolling() {
	joystickManager.filterUpdatingMutex.Lock()
	defer joystickManager.filterUpdatingMutex.Unlock()

	if !joystickManager.isPolling {
		joystickManager.isPolling = true
		go joystickManager.startPolling()
	}
}

func (joystickManager *JoystickManager) startPolling() {
	for {
		time.Sleep(FRAME_DURATION)

		joystickManager.filterUpdatingMutex.Lock()
		if joystickManager.exitPoll {
			for i := range joystickManager.activeJoystickReaders {
				joystickReader := &joystickManager.activeJoystickReaders[i]

				joystickReader.CleanUp()
			}
			break
		}

		for i := range joystickManager.activeJoystickReaders {
			joystickReader := &joystickManager.activeJoystickReaders[i]

			joystickReader.ProcessState()
		}

		joystickManager.filterUpdatingMutex.Unlock()
	}
}
