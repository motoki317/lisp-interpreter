package object

type (
	Env struct {
		frame Frame
		upper *Env
	}
	Frame map[string]Object
)

// NewGlobalEnv returns a new Env with single given frame (i.e. global Env).
func NewGlobalEnv(globalEnv Frame) *Env {
	return &Env{frame: globalEnv}
}

// NewEnv appends the given new frame to the existing Env, not modifying the base Env.
func (e *Env) NewEnv(newFrame Frame) *Env {
	return &Env{
		frame: newFrame,
		upper: e,
	}
}

// Define adds a key value pair in the top Frame.
func (e *Env) Define(key string, value Object) {
	e.frame[key] = value
}

// Set overrides a key value pair in this Env.
// Returns false if the key isn't this Env.
func (e *Env) Set(key string, value Object) (ok bool) {
	cur := e
	for cur != nil {
		if _, ok := cur.frame[key]; ok {
			cur.frame[key] = value
			return true
		}
		cur = cur.upper
	}
	return false
}

// Lookup looks up for the key in this Env.
func (e *Env) Lookup(key string) (value Object, ok bool) {
	cur := e
	for cur != nil {
		if v, ok := cur.frame[key]; ok {
			return v, true
		}
		cur = cur.upper
	}
	return nil, false
}

// EmptyFrame returns an empty new frame.
func EmptyFrame() Frame {
	return make(map[string]Object)
}

// NewBindingFrame returns a new Frame binding the given arguments.
// Expects the length of argNames and the length of objects to be the same.
func NewBindingFrame(argNames []string, objects []Object) Frame {
	f := EmptyFrame()
	if len(argNames) != len(objects) {
		panic("assertion error: len(argNames) == len(objects)")
	}
	for i, arg := range argNames {
		f[arg] = objects[i]
	}
	return f
}
