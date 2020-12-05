package lisp

type (
	env struct {
		frame frame
		upper *env
	}
	frame map[string]*object
)

// newGlobalEnv returns a new env with single given frame (i.e. global env).
func newGlobalEnv(globalEnv frame) *env {
	return &env{frame: globalEnv}
}

// newEnv appends the given new frame to the existing env, not modifying the base env.
func (e *env) newEnv(newFrame frame) *env {
	return &env{
		frame: newFrame,
		upper: e,
	}
}

// define adds or overrides a key value pair in this env.
func (e *env) define(key string, value *object) {
	cur := e
	for cur != nil {
		if _, ok := cur.frame[key]; ok {
			cur.frame[key] = value
			return
		}
		cur = cur.upper
	}
	// no match, so define a new key-value pair to the top frame
	e.frame[key] = value
}

// lookup looks up for the key in this env.
func (e *env) lookup(key string) (value *object, ok bool) {
	cur := e
	for cur != nil {
		if v, ok := cur.frame[key]; ok {
			return v, true
		}
		cur = cur.upper
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
