package multierror

import "errors"

func (me multiError) Unwrap() error {
	return errorChain(me)
}

type errorChain []error

func (ch errorChain) Error() string              { return ch[0].Error() }
func (ch errorChain) Is(target error) bool       { return errors.Is(ch[0], target) }
func (ch errorChain) As(target interface{}) bool { return errors.As(ch[0], target) }

func (ch errorChain) Unwrap() error {
	child := errors.Unwrap(ch[0])
	if child != nil {
		newChain := make(errorChain, len(ch))
		copy(newChain[1:], ch[1:])
		newChain[0] = child
		return newChain
	}

	if len(ch) == 1 {
		return nil
	}
	return ch[1:]
}
