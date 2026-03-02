package engine

import "github.com/hajimehoshi/ebiten/v2"

// inputManager handles keyboard and mouse input
type InputManager struct {
	// mouse
	MouseX, MouseY int
	MousePressed   [3]bool // left, right, middle
	MouseReleased  [3]bool
	MouseHeld      [3]bool
	MousePrevHeld  [3]bool

	// keyboard
	KeysPressed  map[ebiten.Key]bool
	KeysReleased map[ebiten.Key]bool
	KeysHeld     map[ebiten.Key]bool
	KeysPrevHeld map[ebiten.Key]bool

	// mouse wheel delta for the current frame (consumed per game tick).
	WheelX, WheelY float64
}

func NewInputManager() *InputManager {
	return &InputManager{
		KeysPressed:  make(map[ebiten.Key]bool),
		KeysReleased: make(map[ebiten.Key]bool),
		KeysHeld:     make(map[ebiten.Key]bool),
		KeysPrevHeld: make(map[ebiten.Key]bool),
	}
}

// update should be called each frame to refresh input state
func (im *InputManager) Update() {
	// mouse position
	im.MouseX, im.MouseY = ebiten.CursorPosition()
	im.WheelX, im.WheelY = ebiten.Wheel()

	// mouse buttons
	for i := 0; i < 3; i++ {
		im.MousePrevHeld[i] = im.MouseHeld[i]
		btn := ebiten.MouseButton(i)
		im.MouseHeld[i] = ebiten.IsMouseButtonPressed(btn)
		// latch edge events until ConsumeEvents runs during a game tick.
		if im.MouseHeld[i] && !im.MousePrevHeld[i] {
			im.MousePressed[i] = true
		}
		if !im.MouseHeld[i] && im.MousePrevHeld[i] {
			im.MouseReleased[i] = true
		}
	}

	// copy previous keys
	for k := range im.KeysPrevHeld {
		delete(im.KeysPrevHeld, k)
	}
	for k, v := range im.KeysHeld {
		im.KeysPrevHeld[k] = v
	}

	// check all keys
	allKeys := []ebiten.Key{
		ebiten.KeyA, ebiten.KeyB, ebiten.KeyC, ebiten.KeyD, ebiten.KeyE,
		ebiten.KeyF, ebiten.KeyG, ebiten.KeyH, ebiten.KeyI, ebiten.KeyJ,
		ebiten.KeyK, ebiten.KeyL, ebiten.KeyM, ebiten.KeyN, ebiten.KeyO,
		ebiten.KeyP, ebiten.KeyQ, ebiten.KeyR, ebiten.KeyS, ebiten.KeyT,
		ebiten.KeyU, ebiten.KeyV, ebiten.KeyW, ebiten.KeyX, ebiten.KeyY,
		ebiten.KeyZ,
		ebiten.KeyDigit0, ebiten.KeyDigit1, ebiten.KeyDigit2, ebiten.KeyDigit3,
		ebiten.KeyDigit4, ebiten.KeyDigit5, ebiten.KeyDigit6, ebiten.KeyDigit7,
		ebiten.KeyDigit8, ebiten.KeyDigit9,
		ebiten.KeySpace, ebiten.KeyEnter, ebiten.KeyEscape, ebiten.KeyBackspace,
		ebiten.KeyTab, ebiten.KeyShiftLeft, ebiten.KeyControlLeft,
		ebiten.KeyArrowUp, ebiten.KeyArrowDown, ebiten.KeyArrowLeft, ebiten.KeyArrowRight,
		ebiten.KeyF1, ebiten.KeyF2, ebiten.KeyF3, ebiten.KeyF4,
		ebiten.KeyF5, ebiten.KeyF6, ebiten.KeyF7, ebiten.KeyF8,
		ebiten.KeyF9, ebiten.KeyF10, ebiten.KeyF11, ebiten.KeyF12,
	}

	for _, key := range allKeys {
		held := ebiten.IsKeyPressed(key)
		im.KeysHeld[key] = held

		if held && !im.KeysPrevHeld[key] {
			im.KeysPressed[key] = true
		}
		if !held && im.KeysPrevHeld[key] {
			im.KeysReleased[key] = true
		}
	}
}

// mouseLeftPressed returns true on the frame left mouse was clicked
func (im *InputManager) MouseLeftPressed() bool {
	return im.MousePressed[0]
}

// mouseRightPressed returns true on the frame right mouse was clicked
func (im *InputManager) MouseRightPressed() bool {
	return im.MousePressed[1]
}

// mouseLeftHeld returns true while left mouse is held
func (im *InputManager) MouseLeftHeld() bool {
	return im.MouseHeld[0]
}

// keyPressed returns true on the frame a key was pressed
func (im *InputManager) KeyPressed(key ebiten.Key) bool {
	return im.KeysPressed[key]
}

// keyHeld returns true while a key is held
func (im *InputManager) KeyHeld(key ebiten.Key) bool {
	return im.KeysHeld[key]
}

// keyReleased returns true on the frame a key was released
func (im *InputManager) KeyReleased(key ebiten.Key) bool {
	return im.KeysReleased[key]
}

// anyKeyPressed returns true if any key was pressed this frame
func (im *InputManager) AnyKeyPressed() bool {
	for _, v := range im.KeysPressed {
		if v {
			return true
		}
	}
	return false
}

// mouseLeftReleased returns true on the frame left mouse was released
func (im *InputManager) MouseLeftReleased() bool {
	return im.MouseReleased[0]
}

// mouseRightReleased returns true on the frame right mouse was released
func (im *InputManager) MouseRightReleased() bool {
	return im.MouseReleased[1]
}

// wheelDelta returns mouse wheel delta captured this frame.
func (im *InputManager) WheelDelta() (float64, float64) {
	return im.WheelX, im.WheelY
}

// consumeEvents clears all single-frame input events (pressed/released).
// called after each game tick to prevent the same press from firing
// multiple times when multiple ticks run in one frame.
func (im *InputManager) ConsumeEvents() {
	for i := 0; i < 3; i++ {
		im.MousePressed[i] = false
		im.MouseReleased[i] = false
	}
	for k := range im.KeysPressed {
		delete(im.KeysPressed, k)
	}
	for k := range im.KeysReleased {
		delete(im.KeysReleased, k)
	}
	im.WheelX = 0
	im.WheelY = 0
}
