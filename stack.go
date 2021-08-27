package multierror

import "reflect"

func (ch errorChain) StackTrace() []uintptr {
	// This function provides an interface that sentry-go recognizes.
	// It makes a best effort to provide stack information for wrapped
	// errors from a variety of popular libraries.

	if len(ch) == 0 {
		return nil
	}
	err := ch[0]

	// golang.org/x/xerrors
	if pc := xerrorsStackTrace(err); pc != nil {
		return pc
	}

	// github.com/pkg/errors
	if pc := pkgerrorsStackTrace(err); pc != nil {
		return pc
	}

	// github.com/go-errors/errors
	if pc := goerrorsStackTrace(err); pc != nil {
		return pc
	}

	return nil
}

func (ch errorChain) Callers() []uintptr {
	// One of the interfaces supported by bugsnag wants
	// a Callers() function with the same signature.
	return ch.StackTrace()
}

func xerrorsStackTrace(err error) []uintptr {
	val := dereference(reflect.ValueOf(err))
	if val.Kind() != reflect.Struct {
		return nil
	}

	val = dereference(val.FieldByName("frame"))
	if val.Kind() != reflect.Struct || val.IsZero() {
		return nil
	}

	val = dereference(val.FieldByName("frames"))
	if !isArray(val) || val.IsZero() {
		return nil
	}
	val = val.Slice(1, val.Len())

	pc := make([]uintptr, val.Len())
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if !isUint(item) {
			return nil
		}
		pc[i] = uintptr(item.Uint())
	}

	return pc
}

func dereference(val reflect.Value) reflect.Value {
	kind := val.Kind()
	for kind == reflect.Ptr || kind == reflect.Interface {
		val = val.Elem()
		kind = val.Kind()
	}
	return val
}

func isUint(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Uint, reflect.Uintptr, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return true
	default:
		return false
	}
}

func isArray(val reflect.Value) bool {
	k := val.Kind()
	return k == reflect.Slice || k == reflect.Array
}

func pkgerrorsStackTrace(err error) []uintptr {
	val := dereference(reflect.ValueOf(err))
	if val.Kind() != reflect.Struct {
		return nil
	}

	val = val.MethodByName("StackTrace")
	if !val.IsValid() {
		return nil
	}

	ret := val.Call(make([]reflect.Value, 0))
	if len(ret) == 0 || !isArray(ret[0]) {
		return nil
	}
	val = ret[0]

	pc := make([]uintptr, val.Len())
	for i := 0; i < val.Len(); i++ {
		item := val.Index(i)
		if !isUint(item) {
			return nil
		}
		pc[i] = uintptr(item.Uint())
	}

	return pc
}

func goerrorsStackTrace(err error) []uintptr {
	if callers, ok := err.(interface{ Callers() []uintptr }); ok {
		return callers.Callers()
	}
	return nil
}
