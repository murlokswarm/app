package app

// ControllerInput describes a controller input.
type ControllerInput int

// Constants that describe the controller inputs.
const (
	DirectionalPad ControllerInput = iota
	LeftThumbstick
	RightThumbstick
	A
	B
	X
	Y
	L1
	L2
	R1
	R2
	Pause
)

// ControllerConfig is a struct that describes a controller.
type ControllerConfig struct {
	// The function that is called when a directional pad or a thumbstick is
	// used.
	OnDirectionChange func(in ControllerInput, x float64, y float64) `json:"-"`

	// The function that is called when a button in pressed.
	OnButtonPressed func(in ControllerInput, value float64, pressed bool) `json:"-"`

	// The function that is called when the controller is connected.
	OnConnected func() `json:"-"`

	// The function that is called when the controller is disconnected.
	OnDisconnected func() `json:"-"`

	// The function that is called when the pause button is pressed.
	OnPause func() `json:"-"`

	// The function that is called when the controller is closed.
	OnClose func() `json:"-"`
}

// Controller is the interface that describes a controller.
type Controller interface {
	Elem
	Closer
}
