package models

type ButtonMappings struct {
	Buttons   []int
	Axis      []int
	Key       string
	VKKeyCode int
	Shift     bool
	Ctrl      bool
	Alt       bool
}

type ControllerMappings struct {
	Id       int
	Mappings []ButtonMappings
}
