package multierror

import (
	"fmt"
	"sync"
	"sync/atomic"
)

// Group represents a sync.WaitGroup which can start multiple goroutines and
// collect up their errors into a MultiError.
//
// The zero value of Group is usable immediately.
type Group struct {
	wg sync.WaitGroup

	mut  sync.Mutex
	errs multiError

	allowPanic uint32
}

// RecoverPanics toggles panic recovery from goroutines started by Go (default true).
//
// If RecoverPanics is not called or is passed true, panicing goroutines will have
// any panicked values recovered, formatted into an error, and added into the
// MultiError returned by Wait().
func (g *Group) RecoverPanics(toggle bool) {
	// the bool is flipped to 'allowPanic' on the struct so that it will
	// default to false for the zero value Group, thus the logical
	// 'recover panics' bool defaults to true.
	var ui uint32
	if toggle {
		ui = 1
	}
	atomic.StoreUint32(&g.allowPanic, ui)
}

// Go starts a new goroutine running the provided function.
// If the error it returns is non-nil, it will be a part of the
// errors returned by Group.Errors(), and part of the MultiError
// returned by Group.Error().
func (g *Group) Go(f func() error) {
	g.wg.Add(1)

	go func() {
		defer g.wg.Done()

		if atomic.LoadUint32(&g.allowPanic) == 0 {
			defer func() {
				r := recover()
				if r == nil {
					return
				}
				var wrapped error
				if err, ok := r.(error); ok {
					wrapped = fmt.Errorf("panic in goroutine: %w", err)
				} else {
					wrapped = fmt.Errorf("panic in goroutine: %#v", r)
				}

				g.addError(wrapped)
			}()
		}

		g.addError(f())
	}()
}

func (g *Group) addError(err error) {
	if err == nil {
		return
	}
	g.mut.Lock()
	defer g.mut.Unlock()
	g.errs = append(g.errs, err)
}

// Wait blocks until all goroutines started by this group have completed.
func (g *Group) Wait() MultiError {
	g.wg.Wait()

	return g.Error()
}

// Error returns a MultiError wrapping all errors from the Group's goroutines.
func (g *Group) Error() MultiError {
	if g.errs == nil {
		return nil
	}
	return g.errs
}
