package lisp

type (
	env   []frame
	frame map[string]*object
)

// newGlobalEnv returns a new env with single given frame (i.e. global env).
func newGlobalEnv(globalEnv frame) env {
	return []frame{globalEnv}
}

// newEnv appends the given new frame to the existing env, not modifying the base env.
func (e env) newEnv(newFrame frame) env {
	newEnv := make([]frame, len(e)+1)
	newEnv[0] = newFrame
	for i := range e {
		newEnv[i+1] = e[i]
	}
	return newEnv
}

// define adds or overrides a key value pair in this env.
func (e env) define(key string, value *object) {
	for _, frame := range e {
		if _, ok := frame[key]; ok {
			frame[key] = value
			return
		}
	}
	// no match, so define a new key-value pair to the top frame
	e[0][key] = value
}

// lookup looks up for the key in this env.
func (e env) lookup(key string) (value *object, ok bool) {
	for _, frame := range e {
		if v, ok := frame[key]; ok {
			return v, true
		}
	}
	return nil, false
}

// emptyFrame returns an empty new frame.
func emptyFrame() frame {
	return make(map[string]*object)
}

// newBindingFrame returns a new frame binding the given arguments.
// Expects the length of argNames and the length of objects to be the same.
func newBindingFrame(argNames []string, objects []*object) frame {
	f := emptyFrame()
	if len(argNames) != len(objects) {
		panic("assertion error: len(argNames) == len(objects)")
	}
	for i, arg := range argNames {
		f[arg] = objects[i]
	}
	return f
}
