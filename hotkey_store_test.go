package main

import (
	"reflect"
	"testing"
)

func Test_configRoundTrip(t *testing.T) {
	tests := []struct {
		name   string
		config HotkeyJson
	}{
		{name: "Ctrl + Alt + A", config: HotkeyJson{Modifiers: []string{"Ctrl", "Alt"}, Key: "A"}},
		{name: "Win + 5", config: HotkeyJson{Modifiers: []string{"Win"}, Key: "5"}},
		{name: "no modifiers", config: HotkeyJson{Modifiers: []string{}, Key: "Z"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mods, key, err := configToIndexes(tt.config)
			if err != nil {
				t.Fatalf("configToIndexes() error = %v", err)
			}
			got, err := indexesToConfig(mods, key)
			if err != nil {
				t.Fatalf("indexesToConfig() error = %v", err)
			}
			if !reflect.DeepEqual(got, tt.config) {
				t.Errorf("round trip = %v, want %v", got, tt.config)
			}
		})
	}
}

func Test_configToIndexes_invalid(t *testing.T) {
	tests := []struct {
		name   string
		config HotkeyJson
	}{
		{name: "empty (no saved hotkey)", config: HotkeyJson{}},
		{name: "typo in modifier", config: HotkeyJson{Modifiers: []string{"Ctlr"}, Key: "A"}},
		{name: "unsupported key", config: HotkeyJson{Modifiers: []string{"Ctrl"}, Key: "F1"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, _, err := configToIndexes(tt.config); err == nil {
				t.Errorf("configToIndexes(%v) returned no error", tt.config)
			}
		})
	}
}
