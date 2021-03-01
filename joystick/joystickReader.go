package joystick

import (
	. "arcadeprocesscompanion/models"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/micmonay/keybd_event"
	"github.com/simulatedsimian/joystick"
	. "github.com/simulatedsimian/joystick"
)

type JoystickReader struct {
	joystickReference Joystick
	previousButtons   uint32
	previousAxis      []int
	joystickId        int
	buttonMappings    []ButtonMappings
	keyBonding        *keybd_event.KeyBonding
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

func NewJoystickReader(mappings ControllerMappings, keyBonding *keybd_event.KeyBonding) (JoystickReader, error) {
	initialize()
	var reader = JoystickReader{}
	reader.joystickId = mappings.Id
	reader.buttonMappings = mappings.Mappings

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
	reader.keyBonding = keyBonding

	return reader, nil
}

func (reader *JoystickReader) CleanUp() {
	if reader.joystickReference != nil {
		reader.joystickReference.Close()
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
			if areAllButtonsPushedForFirstTime(buttonMapping.Buttons, state.Buttons, reader.previousButtons,
				buttonMapping.Axis, state.AxisData, reader.previousAxis) {
				reader.keyBonding.SetKeys(buttonMapping.VKKeyCode)

				if buttonMapping.Alt {
					reader.keyBonding.HasALT(true)
				}

				if buttonMapping.Ctrl {
					reader.keyBonding.HasCTRL(true)
				}

				if buttonMapping.Shift {
					reader.keyBonding.HasSHIFT(true)
				}

				reader.keyBonding.Press()
				time.Sleep(10 * time.Millisecond)
				reader.keyBonding.Release()
			}
		}
	}

	reader.previousButtons = state.Buttons
	copy(reader.previousAxis, state.AxisData)
}

func areAllButtonsPushedForFirstTime(buttons []int, state uint32, previousState uint32, axisDefinition []int, axisData []int, previousAxisData []int) bool {
	buttonsHaveNotChanged := true
	for i := range buttons {
		btn := buttons[i]

		mask := bitmask[btn]
		if (state & mask) != mask {
			return false
		}

		if buttonsHaveNotChanged && ((state & mask) != (previousState & mask)) {
			buttonsHaveNotChanged = false
		}
	}

	if axisDefinition != nil && len(axisDefinition) > 0 {
		for i := range axisDefinition {
			if axisDefinition[i] == 0 {
				continue
			}

			if (axisData[i] * axisDefinition[i]) <= 0 {
				return false
			}

			if buttonsHaveNotChanged && (axisData[i] != previousAxisData[i]) {
				buttonsHaveNotChanged = false
			}
		}
	}

	if buttonsHaveNotChanged {
		return false
	}

	return true
}
