package models

type RootMappings struct {
	ProcessPriorityOrder []string
	Mappings             []ControllerMappings
}

type ButtonMappings struct {
	Buttons   []int
	Axis      []int
	Key       *string
	VKKeyCode int
	Shift     bool
	Ctrl      bool
	Alt       bool
	Mouse     []int32
	MouseClick int8
	Command   *string
}

type ControllerMappings struct {
	Id       int
	Mappings []ButtonMappings
	Filters  *MappingFilters
	MouseSimulation *MouseSimulation
}

type MappingFilters struct {
	Processes []string
}

type MouseSimulation struct {
	Acceleration float64 //Number between 1 and 2
	MaxSpeed float64
	StartSpeed float64
}
