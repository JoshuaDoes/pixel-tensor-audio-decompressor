package main

import (
	//	"fmt"

	"github.com/MarinX/keylogger"
)

//KeycodeBinding holds a binding between a Linux keycode and a bare Go handler
type KeycodeBinding struct {
	Handler   func() //The binding handler function that will be called when this binding activates
	Keycode   uint16 //The Linux-designated keycode for this binding
	OnRelease bool   //If this binding should activate when the button is released instead of when pressed
}

//KeycodeListener holds a Linux keycode listener
type KeycodeListener struct {
	RootBind  func(keyboard string, keycode uint16, onRelease bool) //Fallback if no other bindings match an event
	Bindings  []*KeycodeBinding
	Keyboard  string
	KeyLogger *keylogger.KeyLogger

	running bool
	closed  bool
}

//Bind binds a keycode to a handler, bind nil to remove all bindings to the keycode
func (kl *KeycodeListener) Bind(keycode uint16, onRelease bool, handler func()) {
	if kl.closed {
		return
	}
	if handler == nil {
		return
	}

	kl.Bindings = append(kl.Bindings, &KeycodeBinding{
		Handler:   handler,
		Keycode:   keycode,
		OnRelease: onRelease,
	})
}

//RemoveBind removes all bindings to a keycode
func (kl *KeycodeListener) RemoveBind(keycode uint16) {
	if kl.closed {
		return
	}
	newBindings := make([]*KeycodeBinding, 0)
	for _, binding := range kl.Bindings {
		if binding.Keycode == keycode {
			continue
		}
		newBindings = append(newBindings, binding)
	}
	kl.Bindings = newBindings
}

//NewKeycodeListener returns a new keycode listener
func NewKeycodeListener(keyboard string) (*KeycodeListener, error) {
	k, err := keylogger.New(keyboard)
	if err != nil {
		return nil, err
	}

	return &KeycodeListener{
		Bindings:  make([]*KeycodeBinding, 0),
		Keyboard:  keyboard,
		KeyLogger: k,
	}, nil
}

//Run starts the keycode listener and blocks until it's closed
func (kl *KeycodeListener) Run() {
	if kl.running {
		return
	}
	kl.closed = false
	kl.running = true

	events := kl.KeyLogger.Read()
	for e := range events {
		if !kl.running {
			break //Exit the keylogger if we're done
		}

		switch e.Type {
		case keylogger.EvKey:
			if e.KeyPress() || e.KeyRelease() {
				//fmt.Printf("<> Handling key (%v|%v): %d\n", e.KeyPress(), e.KeyRelease(), e.Code)
				binded := false
				for _, binding := range kl.Bindings {
					if binding.Keycode == e.Code {
						if e.KeyPress() && !binding.OnRelease {
							binding.Handler()
							binded = true
						} else if e.KeyRelease() && binding.OnRelease {
							binding.Handler()
							binded = true
						}
					}
				}
				if !binded && kl.RootBind != nil {
					kl.RootBind(kl.Keyboard, e.Code, e.KeyRelease())
				}
			}
		}
	}
}

//Close closes the keycode listener
func (kl *KeycodeListener) Close() {
	if kl.closed {
		return
	}
	kl.running = false
	kl.closed = true

	kl.KeyLogger.Close()
}
