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
	Command   *string
}

type ControllerMappings struct {
	Id       int
	Mappings []ButtonMappings
	Filters  *MappingFilters
}

type MappingFilters struct {
	Processes []string
}
