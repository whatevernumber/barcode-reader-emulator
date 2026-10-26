package main

import (
	"github.com/gen2brain/beeep"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"golang.design/x/hotkey"

	"time"
)

func main() {
	primaryBindings := make(chan *hotkey.Hotkey, 1)
	screenBindings := make(chan *hotkey.Hotkey, 1)

	var value *walk.TextEdit
	var hotkeyLabel *walk.Label
	var screenHotkeyLabel *walk.Label
	var slider *walk.Slider
	var mainWnd *walk.MainWindow
	var sendEnter *walk.CheckBox

	config, err := LoadConfig()
	if err != nil {
		config = Config{}
	}

	primaryLabel, err := bindSavedOrDefault(primaryBindings, config.Primary, defaultPrimaryHotkey)
	if err != nil {
		walk.MsgBox(nil, "Failed to bind a default hotkey", err.Error()+"\nTry to change the hotkey to something another", walk.MsgBoxApplModal)
	}

	screenLabel, err := bindSavedOrDefault(screenBindings, config.Screen, defaultScreenHotkey)
	if err != nil {
		walk.MsgBox(nil, "Failed to bind a default scan hotkey", err.Error()+"\nTry to change the hotkey to something another", walk.MsgBoxApplModal)
	}

	showScanResult := func(val string, err error) {
		if err != nil {
			walk.MsgBox(mainWnd, "Failed to scan", "Could not find or parse a barcode. Try to zoom in a bit or choose another one.\nParser error: "+err.Error(), walk.MsgBoxApplModal)
		} else {
			value.SetText(val)
		}
	}

	go ListenForHotkey(primaryBindings, func() {
		time.Sleep(time.Millisecond * 150)

		valueToWrite := value.Text()
		if len(valueToWrite) == 0 {
			return
		}
		beeep.Beep(beeep.DefaultFreq*2, beeep.DefaultDuration/4) // emulate scanner sound :)

		EmulateTyping(valueToWrite, slider.Value(), sendEnter.CheckState() == walk.CheckChecked)
	})

	go ListenForHotkey(screenBindings, func() {
		val, err := ScanScreenBarcode()

		if err == nil {
			beeep.Beep(beeep.DefaultFreq*2, beeep.DefaultDuration/4)
		}

		mainWnd.Synchronize(func() {
			showScanResult(val, err)
		})
	})

	ico, err := walk.NewIconFromResourceId(2) // icon id from manifest
	if err != nil {
		panic("could not load icon")
	}

	changeHotkey := func(bindings chan *hotkey.Hotkey, label *walk.Label, hotkeyConfig *HotkeyJson) {
		var lb *walk.ListBox
		var cb *walk.ComboBox
		var dlg *walk.Dialog

		modLabels := make([]string, len(AllModifiers))
		for i := 0; i < len(AllModifiers); i++ {
			modLabels[i] = modifierToString(AllModifiers[i])
		}

		hotkeyLabels := make([]string, len(AllHotkeys))
		for i := 0; i < len(AllHotkeys); i++ {
			hotkeyLabels[i] = hotkeyToString(AllHotkeys[i])
		}

		Dialog{
			AssignTo: &dlg,
			Title:    "Select a new hotkey",
			Layout:   HBox{},
			MinSize:  Size{Width: 350, Height: 165},
			Children: []Widget{
				ListBox{
					AssignTo:       &lb,
					MultiSelection: true,
					Model:          modLabels,
				},
				ComboBox{
					Model:        hotkeyLabels,
					AssignTo:     &cb,
					CurrentIndex: 0,
				},
				PushButton{
					Text: "OK", OnClicked: func() {
						selModIdx := lb.SelectedIndexes()
						selHotkeyIdx := cb.CurrentIndex()
						bindingLabel, err := SetNewBinding(bindings, selModIdx, selHotkeyIdx)
						if err != nil {
							walk.MsgBox(mainWnd, "Failed", err.Error(), walk.MsgBoxApplModal)
							return
						}

						label.SetText(bindingLabel)
						dlg.Close(walk.DlgCmdOK)

						newHotkey, err := indexesToConfig(selModIdx, selHotkeyIdx)
						if err == nil {
							*hotkeyConfig = newHotkey
							err = changeConfig(&config)
						}
						if err != nil {
							walk.MsgBox(mainWnd, "Could not save the hotkey", "The new hotkey works until you close the app.\n"+err.Error(), walk.MsgBoxApplModal)
						}
					},
				},
			},
		}.Run(mainWnd)
	}

	MainWindow{
		AssignTo: &mainWnd,
		Title:    "Barcode Reader Emulator",
		Size:     Size{Width: 350, Height: 135},
		Icon:     ico,
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					Label{Text: "Value"},
					TextEdit{
						AssignTo:      &value,
						Text:          "test",
						MaxLength:     128,
						CompactHeight: true,
						VScroll:       true,
					},
					PushButton{
						Text:    "Scan from the screen",
						MaxSize: Size{Width: 40, Height: 20},
						OnClicked: func() {
							showScanResult(ScanScreenBarcode())
						},
					},
				},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "Hotkey"},
					Label{AssignTo: &hotkeyLabel, Text: primaryLabel},
					PushButton{
						Text:    "Change",
						MaxSize: Size{Width: 40, Height: 20},
						OnClicked: func() {
							changeHotkey(primaryBindings, hotkeyLabel, &config.Primary)
						},
					},
				},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "Scan hotkey"},
					Label{AssignTo: &screenHotkeyLabel, Text: screenLabel},
					PushButton{
						Text:    "Change",
						MaxSize: Size{Width: 40, Height: 20},
						OnClicked: func() {
							changeHotkey(screenBindings, screenHotkeyLabel, &config.Screen)
						},
					},
				},
			},
			HSplitter{
				Children: []Widget{
					Label{Text: "Input Key Delay"},
					Slider{AssignTo: &slider, MinValue: 10, MaxValue: 100},
					CheckBox{
						Text:     "Send ENTER at the end",
						AssignTo: &sendEnter,
						MinSize:  Size{Width: 40, Height: 20},
					},
				},
			},
		},
	}.Run()
}
