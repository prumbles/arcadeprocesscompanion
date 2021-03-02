package models

import (
	"github.com/micmonay/keybd_event"
)

type ButtonMappingsInternal struct {
	ButtonMappings
	ButtonsMask uint32
	ButtonsPushed bool
	KeyBonding        *keybd_event.KeyBonding
}