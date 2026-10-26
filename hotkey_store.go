package main

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"

	"golang.design/x/hotkey"
)

type Config struct {
	Primary HotkeyJson `json:"primary"`
	Screen  HotkeyJson `json:"screen"`
}

type HotkeyJson struct {
	Modifiers []string `json:"modifiers"`
	Key       string   `json:"key"`
}

var (
	defaultPrimaryHotkey = HotkeyJson{Modifiers: []string{"Ctrl", "Alt"}, Key: "A"}
	defaultScreenHotkey  = HotkeyJson{Modifiers: []string{"Ctrl", "Alt"}, Key: "S"}
)

var (
	configPath string
	origConfig []byte
)

func LoadConfig() (Config, error) {
	var config Config

	dir, err := os.UserConfigDir()

	if err != nil {
		return config, err
	}

	configPath = filepath.Join(dir, "barcode-reader-emulator", "barcode_hotkeys.json")

	file, err := os.ReadFile(configPath)
	if os.IsNotExist(err) {
		origConfig = nil
		return config, nil
	}

	if err != nil {
		return config, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return config, err
	}

	origConfig = file
	return config, nil
}

func changeConfig(config *Config) error {
	file, err := json.MarshalIndent(config, "", "    ")

	if err != nil {
		return err
	}

	return saveConfig(file)
}

func saveConfig(config []byte) error {
	if bytes.Equal(config, origConfig) {
		return nil
	}

	if configPath == "" {
		return os.ErrInvalid
	}

	err := os.MkdirAll(filepath.Dir(configPath), 0700)
	if err != nil {
		return err
	}

	err = os.WriteFile(configPath, config, 0600)
	if err != nil {
		return err
	}

	origConfig = config
	return nil
}

func configToIndexes(config HotkeyJson) ([]int, int, error) {
	modifierIndexes := make([]int, 0, len(config.Modifiers))

	for _, label := range config.Modifiers {
		index, err := modifierToIndex(label)
		if err != nil {
			return nil, -1, err
		}
		modifierIndexes = append(modifierIndexes, index)
	}

	keyIndex, err := hotkeyToIndex(config.Key)
	if err != nil {
		return nil, -1, err
	}

	return modifierIndexes, keyIndex, nil
}

func indexesToConfig(modifierIndexes []int, keyIndex int) (HotkeyJson, error) {
	config := HotkeyJson{
		Modifiers: make([]string, 0, len(modifierIndexes)),
	}

	for _, index := range modifierIndexes {
		config.Modifiers = append(
			config.Modifiers,
			modifierToString(AllModifiers[index]),
		)
	}

	config.Key = hotkeyToString(AllHotkeys[keyIndex])
	return config, nil
}

func bindSavedOrDefault(bindings chan *hotkey.Hotkey, saved HotkeyJson, fallback HotkeyJson) (string, error) {
	label, err := bindHotkey(bindings, saved)
	if err == nil {
		return label, nil
	}
	return bindHotkey(bindings, fallback)
}

func bindHotkey(bindings chan *hotkey.Hotkey, config HotkeyJson) (string, error) {
	modifierIndexes, keyIndex, err := configToIndexes(config)
	if err != nil {
		return "", err
	}
	return SetNewBinding(bindings, modifierIndexes, keyIndex)
}
