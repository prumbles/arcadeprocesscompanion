package actions

import (
	"fmt"

	"github.com/itchyny/volume-go"
)

func ChangeVolume(diff int) {
	vol, err := volume.GetVolume()
	if err != nil {
		fmt.Printf("get volume failed: %+v", err)
		return
	}

	vol += diff

	if vol < 0 {
		vol = 0
	} else if vol > 100 {
		vol = 100
	}

	err = volume.IncreaseVolume(diff)
	if err != nil {
		fmt.Printf("set volume failed: %+v", err)
	}
}

func Mute() {
	isMuted, err := volume.GetMuted()

	if err != nil {
		fmt.Printf("get muted failed: %+v", err)
		return
	}

	if isMuted {
		err = volume.Unmute()
		if err != nil {
			fmt.Printf("un-mute failed: %+v", err)
			return
		}
	} else {
		err = volume.Mute()
		if err != nil {
			fmt.Printf("mute failed: %+v", err)
			return
		}
	}

}
