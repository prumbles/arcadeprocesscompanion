package joystick

import (
	. "arcadeprocesscompanion/models"
	"errors"
	"fmt"
	"math"
	"strconv"
	"os/exec"

	"github.com/micmonay/keybd_event"
	"github.com/simulatedsimian/joystick"
	. "github.com/simulatedsimian/joystick"
	"gopkg.in/bendahl/uinput.v1"
)


type JoystickReader struct {
	joystickReference Joystick
	previousButtons   uint32
	previousAxis      []int
	joystickId        int
	buttonMappings    []ButtonMappingsInternal
	processId int
	mouse             uinput.Mouse
	mouseSimulation *MouseSimulation
	mouseSpeed float64
	mouseAccelerationSeed float64
}

var bitmask = make([]uint32, 40)

var initAlreadyExecuted = false

func initialize() {
	if !initAlreadyExecuted {
		for i := 0; i < 40; i++ {
			bitmask[i] = uint32(math.Pow(2, float64(i)))
		}
	}

	initAlreadyExecuted = true
}

func NewJoystickReader(mappings ControllerMappings, pId int) (JoystickReader, error) {
	initialize()
	var reader = JoystickReader{}
	reader.joystickId = mappings.Id
	reader.processId = pId

	if mappings.MouseSimulation != nil {
		reader.mouseSimulation = mappings.MouseSimulation

		if reader.mouseSimulation.Acceleration < 0 || reader.mouseSimulation.Acceleration > 4 {
			reader.mouseSimulation.Acceleration = 1.0
		}

		if reader.mouseSimulation.MaxSpeed < 1 || reader.mouseSimulation.MaxSpeed > 100 {
			reader.mouseSimulation.MaxSpeed = 60
		}

		if reader.mouseSimulation.StartSpeed <=0 || reader.mouseSimulation.StartSpeed > 30{
			reader.mouseSimulation.StartSpeed = 1.0
		}
	}

	reader.buttonMappings = make([]ButtonMappingsInternal, len(mappings.Mappings))

	for i := range mappings.Mappings {
		var mask uint32 = 0

		for _,btn := range mappings.Mappings[i].Buttons {
			mask = mask | uint32(math.Pow(2, float64(btn)))
		}

		reader.buttonMappings[i] = ButtonMappingsInternal{
			ButtonMappings: mappings.Mappings[i],
			ButtonsMask: mask,
		}

		kb, err := keybd_event.NewKeyBonding()
		if err != nil {
			panic(err)
		}

		reader.buttonMappings[i].KeyBonding = &kb
	}

	if mappings.Mappings == nil || len(mappings.Mappings) == 0 {
		return reader, errors.New("Mappings cannot be empty")
	}

	js, err := joystick.Open(mappings.Id)
	if err != nil {
		fmt.Printf("reader: %v\n", reader)
		fmt.Printf("mappings: %v\n", mappings)
		panic(err)
	}

	fmt.Printf("Joystick Name: %s\n", js.Name())
	fmt.Printf("   Axis Count: %d\n", js.AxisCount())
	fmt.Printf(" Button Count: %d\n", js.ButtonCount())

	reader.joystickReference = js
	reader.previousButtons = 0
	reader.previousAxis = make([]int, 2)
	reader.previousAxis[0] = 0
	reader.previousAxis[1] = 0

	mouse, mouseError := uinput.CreateMouse("/dev/uinput", []byte("arcadeprocesscompanionmouse"))

	if mouseError == nil {
		reader.mouse = mouse
	} else {
		reader.mouse = nil
	}

	return reader, nil
}

func (reader *JoystickReader) CleanUp() {
	if reader.joystickReference != nil {
		reader.joystickReference.Close()
	}

	if reader.mouse != nil {
		reader.mouse.Close()
	}

	for _,mapping := range reader.buttonMappings {
		mapping.KeyBonding.Clear()
	}
}

func (reader *JoystickReader) ProcessState() {
	if reader.joystickReference == nil {
		return
	}

	state, err := reader.joystickReference.Read()
	if err != nil {
		panic(err)
	}

	if state.Buttons > 0 || state.AxisData[0] != 0 || state.AxisData[1] != 0 {
		for i := range reader.buttonMappings {
			buttonMapping := &reader.buttonMappings[i]
			prevButtonsPushed := buttonMapping.ButtonsPushed

			buttonsHaveNotChanged, allButtonsPushed := getButtonsPushed(buttonMapping.ButtonsMask, state.Buttons, reader.previousButtons,
				buttonMapping.Axis, state.AxisData, reader.previousAxis)

			if allButtonsPushed && !buttonsHaveNotChanged {
				keyPressed := false

				if buttonMapping.Key != nil {
					buttonMapping.KeyBonding.SetKeys(buttonMapping.VKKeyCode)
					keyPressed = true
				}

				if buttonMapping.Alt {
					buttonMapping.KeyBonding.HasALT(true)
					keyPressed = true
				}

				if buttonMapping.Ctrl {
					buttonMapping.KeyBonding.HasCTRL(true)
					keyPressed = true
				}

				if buttonMapping.Shift {
					buttonMapping.KeyBonding.HasSHIFT(true)
					keyPressed = true
				}

				if keyPressed {
					buttonMapping.ButtonsPushed = true
					buttonMapping.KeyBonding.Press()					
				}

				if buttonMapping.Command != nil {
					cmdStr := *buttonMapping.Command

					if reader.processId > 0 {
						cmdStr = "PROCESS_ID=" + strconv.Itoa(reader.processId) + " && " + cmdStr
					}
					fmt.Printf("%v\n", cmdStr)
					cmd := exec.Command("bash", "-c", cmdStr)

					err := cmd.Run()

					if err != nil {
						fmt.Printf(err.Error())
					}
				}

				if len(buttonMapping.Mouse) == 2 {
					if buttonMapping.Mouse[0] > 0 {
						reader.mouse.MoveRight(buttonMapping.Mouse[0])

					} else if buttonMapping.Mouse[0] < 0 {
						reader.mouse.MoveLeft(buttonMapping.Mouse[0] * -1)
					}

					if buttonMapping.Mouse[1] > 0 {
						reader.mouse.MoveDown(buttonMapping.Mouse[1])
					} else if buttonMapping.Mouse[1] < 0 {
						reader.mouse.MoveUp(buttonMapping.Mouse[1] * -1)
					}
				}

				if buttonMapping.MouseClick > 0 {
					if buttonMapping.MouseClick == 1 {
						reader.mouse.LeftClick()
					} else if buttonMapping.MouseClick == 2 {
						reader.mouse.RightClick()
					}
				}

			} else if prevButtonsPushed && !buttonsHaveNotChanged{
				buttonMapping.ButtonsPushed = false
				buttonMapping.KeyBonding.Release()
			}
		}
	} else {
		for i := range reader.buttonMappings {
			buttonMapping := &reader.buttonMappings[i]
			if buttonMapping.ButtonsPushed {
				buttonMapping.ButtonsPushed = false
				buttonMapping.KeyBonding.Release()
			}
		}
	}

	if reader.mouseSimulation != nil {
		if state.AxisData[0] == 0 && state.AxisData[1] == 0 {
			reader.mouseSpeed = 0
			reader.mouseAccelerationSeed = 0
		} else {
			if reader.mouseSpeed < reader.mouseSimulation.MaxSpeed {
				reader.mouseSpeed += (math.Pow(reader.mouseSimulation.Acceleration,reader.mouseAccelerationSeed)) * reader.mouseSimulation.StartSpeed

				if reader.mouseSpeed > reader.mouseSimulation.MaxSpeed {
					reader.mouseSpeed = reader.mouseSimulation.MaxSpeed
				}
			}


			if state.AxisData[0] > 0 {
				reader.mouse.MoveRight(int32(reader.mouseSpeed))
			} else if state.AxisData[0] < 0 {
				reader.mouse.MoveLeft(int32(reader.mouseSpeed))
			}

			if state.AxisData[1] < 0 {
				reader.mouse.MoveUp(int32(reader.mouseSpeed))
			} else if state.AxisData[1] > 0 {
				reader.mouse.MoveDown(int32(reader.mouseSpeed))
			}

			if reader.mouseAccelerationSeed < 1000 && reader.mouseSpeed <= reader.mouseSimulation.MaxSpeed {
				reader.mouseAccelerationSeed ++
			}
		}
	}

	reader.previousButtons = state.Buttons
	copy(reader.previousAxis, state.AxisData)
}

func getButtonsPushed(buttonsMask uint32, state uint32, previousState uint32, axisDefinition []int, axisData []int, previousAxisData []int) (noChange bool, allButtonsPushed bool) {
	buttonsHaveNotChanged := (state & buttonsMask) == (previousState & buttonsMask)
	buttonsAreAllPushed := (state & buttonsMask) == buttonsMask

	if (buttonsMask != 0) {
		if !buttonsAreAllPushed {
			return buttonsHaveNotChanged, false
		}
	}

	if axisDefinition != nil && len(axisDefinition) > 0 {
		for i := range axisDefinition {
			if axisDefinition[i] == 0 {
				continue
			}

			if (axisData[i] * axisDefinition[i]) <= 0 {
				return buttonsHaveNotChanged, false
			}

			if buttonsHaveNotChanged && (axisData[i] != previousAxisData[i]) {
				buttonsHaveNotChanged = false
			}
		}
	}

	return buttonsHaveNotChanged, buttonsAreAllPushed
}
