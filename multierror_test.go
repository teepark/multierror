package multierror_test

import (
	"errors"
	"testing"

	"github.com/teepark/multierror"
)

type SomeError struct{}

func (s *SomeError) Error() string {
	return "yep that's an error"
}

func TestMultiErrorIsWrapped(t *testing.T) {
	err := errors.New("an example error")
	merr := multierror.Wrap(err)
	if !errors.Is(merr, err) {
		t.Errorf("MultiError is not equivalent to its wrapped error")
	}
}

func TestMultiErrorAsWrapped(t *testing.T) {
	err := multierror.Wrap(new(SomeError))
	target := new(SomeError)

	if !errors.As(err, &target) {
		t.Errorf("MultiError does not As() into its wrapped type")
	}
}

func TestMultiErrorIsAnEarlierWrapped(t *testing.T) {
	err1 := errors.New("first example error")
	err2 := errors.New("second example error")

	err := multierror.Wrap(err1)
	err = err.Add(err2)

	if !errors.Is(err, err2) {
		t.Errorf("MultiError is not equivalient to last-added error")
	}

	if !errors.Is(err, err1) {
		t.Errorf("MultiError is not equivalent to an earlier-wrapped error")
	}
}

func TestMultiErrorAsAnEarlierWrapped(t *testing.T) {
	err := multierror.Wrap(new(SomeError))
	err = err.Add(errors.New("second error to maybe hide the SomeError"))

	target := new(SomeError)

	if !errors.As(err, &target) {
		t.Errorf("MultiError does not As() into an earlier-wrapped type")
	}
}

func TestMultiErrorAddFlattensArgument(t *testing.T) {
	err := multierror.Wrap(new(SomeError))
	err2 := multierror.Wrap(new(SomeError), new(SomeError))

	result := err.Add(err2)
	if len(result.Errors()) != 3 {
		t.Errorf("Add should flatten the argument if it is a MultiError")
	}
}

func TestMultiErrorAddFlattensArgumentRecursively(t *testing.T) {
	err := multierror.Wrap(new(SomeError))
	grandchildren := multierror.Wrap(new(SomeError), new(SomeError))
	err2 := multierror.Wrap(new(SomeError), grandchildren)

	result := err.Add(err2)
	if len(result.Errors()) != 4 {
		t.Errorf("Add should recursively flatten a MultiError argument")
	}
}

func TestNilMultiErrorIsNil(t *testing.T) {
	err := multierror.Wrap(nil)
	if err != nil {
		t.Errorf("Wrap(nil) != nil")
	}
}

func TestNilMultiErrorIsNotAnError(t *testing.T) {
	err := multierror.Wrap(nil)
	target := errors.New("testing")
	// really just testing this doesn't dereference nil
	if errors.Is(err, target) {
		t.Errorf("Wrap(nil) equivalent to a real error?")
	}
}

func TestNilMultiErrorAsSomeType(t *testing.T) {
	err := multierror.Wrap(nil)
	target := new(SomeError)
	// really just testing this doesn't dereference nil
	if errors.As(err, &target) {
		t.Errorf("Wrap(nil) unwraps into SomeError type")
	}
}

func TestMultiErrorAsItself(t *testing.T) {
	err := multierror.Wrap(errors.New("foobar"))
	target := multierror.Wrap(errors.New("thing"))

	if !errors.As(err, &target) {
		t.Errorf("MultiError does not As() into its own type")
	}
}

func TestMultiErrorIsItself(t *testing.T) {
	err := multierror.Wrap(errors.New("foo"), errors.New("bar"))

	if !errors.Is(err, err) {
		t.Error("MultiError is not itself")
	}
}
