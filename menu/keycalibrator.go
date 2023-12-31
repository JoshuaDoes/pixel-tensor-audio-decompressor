package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/JoshuaDoes/json"
)

type MenuKeycodeBinding struct {
	Keycode   uint16 `json:"keycode"`
	Action    string `json:"action"`
	OnRelease bool   `json:"onRelease"`
}

func bindKeys() {
	for keyboard, bindings := range keyCalibration {
		kl, err := NewKeycodeListener(keyboard)
		if err != nil {
			panic(fmt.Sprintf("error listening to keyboard %s: %v", keyboard, err))
		}
		for _, binding := range bindings {
			var action func()
			switch binding.Action {
			case "prevItem":
				action = menuEngine.PrevItem
			case "nextItem":
				action = menuEngine.NextItem
			case "selectItem":
				action = menuEngine.Action
			default:
				panic("unknown action: " + binding.Action)
			}
			kl.Bind(binding.Keycode, binding.OnRelease, action)
		}
		go kl.Run()
	}
}

type KeyCalibration struct {
	Ready  bool
	Action string
	KLs    []*KeycodeListener
}

func (kc *KeyCalibration) Input(keyboard string, keycode uint16, onRelease bool) {
	if !kc.Ready {
		os.Exit(0)
	}
	if kc.Action == "" || kc.Action == "cancel" {
		kc.Action = ""
		return
	}
	if onRelease {
		return
	}
	if keyCalibration[keyboard] == nil {
		keyCalibration[keyboard] = make([]*MenuKeycodeBinding, 0)
	}
	keyCalibration[keyboard] = append(keyCalibration[keyboard], &MenuKeycodeBinding{
		Keycode:   keycode,
		Action:    kc.Action,
		OnRelease: true,
	})
	kc.Action = ""
}

func calibrate() {
	//Generate a key calibration file if one doesn't exist yet
	calibrator := &KeyCalibration{KLs: make([]*KeycodeListener, 0)}

	//Get a list of keyboards
	keyboards := make([]string, 0)
	err := filepath.Walk("/dev/input", func(path string, info os.FileInfo, err error) error {
		if len(path) < 16 || string(path[:16]) != "/dev/input/event" {
			return nil
		}
		keyboards = append(keyboards, path)
		return nil
	})
	if err != nil {
		panic(fmt.Sprintf("error walking inputs: %v", err))
	}

	//Bind all keyboards to calibrator input
	for _, keyboard := range keyboards {
		kl, err := NewKeycodeListener(keyboard)
		if err != nil {
			panic(fmt.Sprintf("error listening to walked keyboard %s: %v", keyboard, err))
		}
		kl.RootBind = calibrator.Input
		calibrator.KLs = append(calibrator.KLs, kl)
		go kl.Run()
	}

	//Start calibrating!
	stages := 6
	for stage := 0; stage < stages; stage++ {
		switch stage {
		case 0:
			keyCalibrationJSON, err := ioutil.ReadFile(keyCalibrationFile)
			if err == nil {
				keyCalibration = make(map[string][]*MenuKeycodeBinding)
				err = json.Unmarshal(keyCalibrationJSON, &keyCalibration)
				if err != nil {
					stage = 1
					continue
				}

				clear(3)
				fmt.Println("\t\tYou've calibrated before!\n")
				time.Sleep(time.Second * 2)
				fmt.Println("\t\tPress any key within\n\t\t5 seconds to recalibrate.\n")
				calibrator.Ready = true
				calibrator.Action = "cancel"
				timeout := time.Now()
				for calibrator.Action != "" {
					if time.Now().Sub(timeout).Seconds() > 5 {
						break
					}
					time.Sleep(time.Millisecond * 100)
				}
				if time.Now().Sub(timeout).Seconds() < 5 {
					calibrator.Action = ""
					fmt.Println("\t\tRecalibration time!")
					time.Sleep(time.Second * 2)
					continue
				}
				stage = stages-1 //Skip to the end of the stages
			}
		case 1:
			calibrator.Ready = false
			keyCalibration = make(map[string][]*MenuKeycodeBinding)
			clear(5)
			fmt.Println("\t\tWelcome to the calibrator!\n")
			fmt.Println("\t\tPress any key to cancel.\n")
			time.Sleep(time.Second * 2)
			fmt.Println("\t\tControllers and remotes\n\t\tare also supported.\n")
			time.Sleep(time.Second * 2)
			fmt.Println("\t\tThis is a guided process.\n")
			time.Sleep(time.Second * 2)
			fmt.Println("\t\tGet ready!\n")
			time.Sleep(time.Second * 3)
		case 2:
			clear(1)
			calibrator.Ready = true
			calibrator.Action = "nextItem"
			fmt.Printf("\n")
			fmt.Println("\t\tPress any key to use to\n\t\tnavigate down in a menu.\n")
			fmt.Println("\t\t\tRecommended: volume down")
			for calibrator.Action != "" {
			}
		case 3:
			calibrator.Action = "prevItem"
			fmt.Printf("\n")
			fmt.Println("\t\tPress any key to use to\n\t\tnavigate up in a menu.\n")
			fmt.Println("\t\t\tRecommended: volume up")
			for calibrator.Action != "" {
			}
		case 4:
			calibrator.Action = "selectItem"
			fmt.Printf("\n")
			fmt.Println("\t\tPress any key to use to\n\t\tselect a menu item.\n")
			fmt.Println("\t\t\tRecommended: touch screen")
			for calibrator.Action != "" {
			}
		case 5:
			clear(5)
			fmt.Println("\t\tSaving results...\n")
			keyboards, err := json.Marshal(keyCalibration, true)
			if err != nil {
				panic(fmt.Sprintf("error encoding calibration results: %v", err))
			}
			keyboardsFile, err := os.Create(keyCalibrationFile)
			if err != nil {
				panic(fmt.Sprintf("error creating calibration file: %v", err))
			}
			defer keyboardsFile.Close()
			_, err = keyboardsFile.Write(keyboards)
			if err != nil {
				panic(fmt.Sprintf("error writing calibration file: %v", err))
			}
			//fmt.Println(string(keyboards))
			fmt.Println("\t\tCalibration complete!")
			time.Sleep(time.Second * 2)
			//calibrator.Ready = false
		}
	}

	for i := 0; i < len(calibrator.KLs); i++ {
		calibrator.KLs[i].Close()
	}
}
