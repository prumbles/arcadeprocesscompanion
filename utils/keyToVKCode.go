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

	case "DELETE":
		return keybd_event.VK_DELETE, nil

	case "CAPS":
		return keybd_event.VK_CAPSLOCK, nil

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

	case "A":
		return keybd_event.VK_A, nil

	case "B":
		return keybd_event.VK_B, nil

	case "C":
		return keybd_event.VK_C, nil

	case "D":
		return keybd_event.VK_D, nil

	case "E":
		return keybd_event.VK_E, nil

	case "F":
		return keybd_event.VK_F, nil

	case "G":
		return keybd_event.VK_G, nil

	case "H":
		return keybd_event.VK_H, nil

	case "I":
		return keybd_event.VK_I, nil

	case "J":
		return keybd_event.VK_J, nil

	case "K":
		return keybd_event.VK_K, nil

	case "L":
		return keybd_event.VK_L, nil

	case "M":
		return keybd_event.VK_M, nil

	case "N":
		return keybd_event.VK_N, nil

	case "O":
		return keybd_event.VK_O, nil

	case "P":
		return keybd_event.VK_P, nil

	case "Q":
		return keybd_event.VK_Q, nil

	case "R":
		return keybd_event.VK_R, nil

	case "S":
		return keybd_event.VK_S, nil

	case "T":
		return keybd_event.VK_T, nil

	case "U":
		return keybd_event.VK_U, nil

	case "V":
		return keybd_event.VK_V, nil

	case "W":
		return keybd_event.VK_W, nil

	case "X":
		return keybd_event.VK_X, nil

	case "Y":
		return keybd_event.VK_Y, nil

	case "Z":
		return keybd_event.VK_Z, nil

	case "0":
		return keybd_event.VK_0, nil

	case "1":
		return keybd_event.VK_1, nil

	case "2":
		return keybd_event.VK_2, nil

	case "3":
		return keybd_event.VK_3, nil

	case "4":
		return keybd_event.VK_4, nil

	case "5":
		return keybd_event.VK_5, nil

	case "6":
		return keybd_event.VK_6, nil

	case "7":
		return keybd_event.VK_7, nil

	case "8":
		return keybd_event.VK_8, nil

	case "9":
		return keybd_event.VK_9, nil

	default:
		return 0, errors.New("Key " + keyStr + " doesn't map to any key")
	}
}
