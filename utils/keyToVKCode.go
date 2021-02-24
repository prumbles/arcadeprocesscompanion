package utils

import (
	"errors"

	"github.com/micmonay/keybd_event"
)

func GetVKCode(keyStr string) (int, error) {
	switch keyStr {
	case "ESC":
		return keybd_event.VK_ESC, nil

	case "ENTER":
		return keybd_event.VK_ENTER, nil

	case "SPACE":
		return keybd_event.VK_SPACE, nil

	case "BACKSPACE":
		return keybd_event.VK_BACKSPACE, nil

	case "TAB":
		return keybd_event.VK_TAB, nil

	case "UP":
		return keybd_event.VK_UP, nil

	case "DOWN":
		return keybd_event.VK_DOWN, nil

	case "LEFT":
		return keybd_event.VK_LEFT, nil

	case "RIGHT":
		return keybd_event.VK_RIGHT, nil

	case "F1":
		return keybd_event.VK_F1, nil

	case "F2":
		return keybd_event.VK_F2, nil

	case "F3":
		return keybd_event.VK_F3, nil

	case "F4":
		return keybd_event.VK_F4, nil

	case "F5":
		return keybd_event.VK_F5, nil

	case "F6":
		return keybd_event.VK_F6, nil

	case "F7":
		return keybd_event.VK_F7, nil

	case "F8":
		return keybd_event.VK_F8, nil

	case "F9":
		return keybd_event.VK_F9, nil

	case "F10":
		return keybd_event.VK_F10, nil

	case "F11":
		return keybd_event.VK_F11, nil

	case "F12":
		return keybd_event.VK_F12, nil

	default:
		return 0, errors.New("Key " + keyStr + " doesn't map to any key")
	}
}
