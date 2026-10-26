package main

import (
	"fmt"
	"strings"
	"time"

	"golang.design/x/hotkey"
)

// Listen for the already bind hotkey event. Blocks the calling thread.
func ListenForHotkey(bindings chan *hotkey.Hotkey, actionCallback func()) {
	binding := <-bindings
	for {
		select {
		case newBinding := <-bindings:
			if binding != newBinding {
				err := binding.Unregister()
				if err != nil {
					panic(err)
				}
				binding = newBinding
			}
			continue
		case _, ok := <-binding.Keyup():
			if !ok {
				continue
			}
		default:
			time.Sleep(time.Millisecond * 100)
			continue
		}

		actionCallback()
	}
}

// Creates a new binding for a hotkey with optional modifiers
func SetNewBinding(bindings chan *hotkey.Hotkey, modifierIndexes []int, keyIndex int) (string, error) {
	var hkMods = make([]hotkey.Modifier, len(modifierIndexes))
	for i, selectedIdx := range modifierIndexes {
		hkMods[i] = AllModifiers[selectedIdx]
	}
	key := AllHotkeys[keyIndex]
	hkBind := hotkey.New(hkMods, key)
	err := hkBind.Register()
	if err != nil {
		return "", err
	}

	bindings <- hkBind
	return bindingToLabel(hkMods, key), nil
}

// All possible modifiers for a hotkey binding
var AllModifiers = []hotkey.Modifier{hotkey.ModCtrl, hotkey.ModAlt, hotkey.ModWin, hotkey.ModShift}

// All possible hotkeys for a binding
var AllHotkeys = []hotkey.Key{
	hotkey.Key0,
	hotkey.Key1,
	hotkey.Key2,
	hotkey.Key3,
	hotkey.Key4,
	hotkey.Key5,
	hotkey.Key6,
	hotkey.Key7,
	hotkey.Key8,
	hotkey.Key9,
	hotkey.KeyA,
	hotkey.KeyB,
	hotkey.KeyC,
	hotkey.KeyD,
	hotkey.KeyE,
	hotkey.KeyF,
	hotkey.KeyG,
	hotkey.KeyH,
	hotkey.KeyI,
	hotkey.KeyJ,
	hotkey.KeyK,
	hotkey.KeyL,
	hotkey.KeyM,
	hotkey.KeyN,
	hotkey.KeyO,
	hotkey.KeyP,
	hotkey.KeyQ,
	hotkey.KeyR,
	hotkey.KeyS,
	hotkey.KeyT,
	hotkey.KeyU,
	hotkey.KeyV,
	hotkey.KeyW,
	hotkey.KeyX,
	hotkey.KeyY,
	hotkey.KeyZ,
}

// some utility methods...

func bindingToLabel(mods []hotkey.Modifier, key hotkey.Key) string {
	s := new(strings.Builder)
	for _, v := range mods {
		if v == 0x0 {
			continue
		}
		s.WriteString(modifierToString(v) + " + ")
	}
	s.WriteString(hotkeyToString(key))
	return s.String()
}

func modifierToString(md hotkey.Modifier) string {
	switch md {
	case hotkey.ModAlt:
		return "Alt"
	case hotkey.ModCtrl:
		return "Ctrl"
	case hotkey.ModShift:
		return "Shift"
	case hotkey.ModWin:
		return "Win"
	default:
		panic("unexpected modifier " + string(md))
	}
}

func hotkeyToString(v hotkey.Key) string {
	if v <= 0x39 {
		return string((v - 0x30) + '0')
	} else if v <= 0x5A {
		return string((v - 0x41) + 'A')
	} else {
		panic("unexpected key found: " + string(v))
	}
}

func modifierToIndex(label string) (int, error) {
	for i, modifier := range AllModifiers {
		if modifierToString(modifier) == label {
			return i, nil
		}
	}

	return -1, fmt.Errorf("unknown modifier: %q", label)
}

func hotkeyToIndex(label string) (int, error) {
	for i, key := range AllHotkeys {
		if hotkeyToString(key) == label {
			return i, nil
		}
	}

	return -1, fmt.Errorf("unknown hotkey: %q", label)
}
