package actions

import "strings"

const VOLUME_DEFAULT_CHANGE = 5

func PerformAction(actionStr string) {
	if strings.HasPrefix(actionStr, "VOLUME_UP") {
		ChangeVolume(VOLUME_DEFAULT_CHANGE)
	} else if strings.HasPrefix(actionStr, "VOLUME_DOWN") {
		ChangeVolume(VOLUME_DEFAULT_CHANGE * -1)
	} else if strings.HasPrefix(actionStr, "MUTE") {
		Mute()
	}
}
