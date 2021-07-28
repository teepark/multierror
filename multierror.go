package multierror

import (
	"bytes"
	"errors"
	"fmt"
)

// MultiError is a container for multiple errors, and also an error itself.
//
// Implementations should support the Is() and As() protocols so that
// instances evaluate as equivalent to, and resolve into, any of their
// contained errors.
//
// Unwrap() is not necessarily supported however, as it can be impossible
// for a multi error to provide a single wrapped error.
type MultiError interface {
	error

	// Add another error into the MultiError, producing a new instance.
	Add(error) MultiError

	// Errors returns all wrapped errors contained by the MultiError.
	Errors() []error
}

type multiError []error

// Wrap creates a MultiError around initial error instances.
//
// Wrapping an empty list or all-nils returns nil.
func Wrap(errs ...error) MultiError {
	var errors multiError
	for _, err := range errs {
		errors = errors.Add(err).(multiError)
	}
	if len(errors) == 0 {
		return nil
	}
	return errors
}

func (me multiError) Error() string {
	switch len(me) {
	case 0:
		return ""
	case 1:
		return me[0].Error()
	default:
		b := new(bytes.Buffer)
		for i, err := range me {
			if i > 0 {
				b.Write([]byte{';', ' '})
			}
			fmt.Fprintf(b, "error #%d: %s", i+1, err.Error())
		}
		return b.String()
	}
}

func (me multiError) Add(err error) MultiError {
	if err == nil {
		return me
	}
	var errs []error
	if merr, ok := err.(MultiError); ok {
		errs = append(me.Errors(), merr.Errors()...)
	} else {
		errs = append(me.Errors(), err)
	}
	return multiError(errs)
}

func (me multiError) Errors() []error {
	if len(me) == 0 {
		return nil
	}
	result := make([]error, len(me))
	copy(result, me)
	return result
}

func (me multiError) Is(target error) bool {
	// When comparing to another multiError instance we're looking
	// for an exact match: all contained errors must match, in order.
	if tgt, ok := target.(*multiError); ok {
		target = *tgt
	}
	if tgt, ok := target.(multiError); ok {
		if len(me) != len(tgt) {
			return false
		}
		for i, a := range me {
			b := tgt[i]
			if !errors.Is(a, b) {
				return false
			}
		}

		return true
	}

	for _, err := range me {
		if errors.Is(err, target) {
			return true
		}
	}
	return false
}

func (me multiError) As(target interface{}) bool {
	if tgt, ok := target.(*multiError); ok {
		if cap(*tgt) >= len(me) {
			// target is a *multiError with enough capacity
			*tgt = (*tgt)[:len(me)]
			copy(*tgt, me)
			return true
		} else {
			*tgt = multiError(me.Errors())
		}
	}
	for _, err := range me {
		if errors.As(err, target) {
			return true
		}
	}
	return false
}
