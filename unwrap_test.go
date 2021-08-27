package multierror_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/teepark/multierror"
)

func TestUnwrapSingleErr(t *testing.T) {
	err := errors.New("A")
	merr := multierror.Wrap(err)

	unwr := errors.Unwrap(merr)
	if unwr == nil {
		t.Fatalf("Unwrap returned nil")
	}

	if unwr.Error() != err.Error() {
		t.Errorf("unwrap(wrap(err)) and error messages don't match")
	}

	if !errors.Is(unwr, err) {
		t.Errorf("unwrap(wrap(err)) *is* not err")
	}
}

func TestUnwrapSingleChain(t *testing.T) {
	err := multierror.Wrap(wrap("C", wrap("B", errors.New("A"))))

	unwr := errors.Unwrap(err)
	assertMsg(t, "C: B: A", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "B: A", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "A", unwr)

	unwr = errors.Unwrap(unwr)
	if unwr != nil {
		t.Errorf("expected nil unwrap, got %+v", unwr)
	}
}

func TestUnwrapTwoChains(t *testing.T) {
	err := multierror.Wrap(
		wrap("C", wrap("B", errors.New("A"))),
		wrap("Z", wrap("Y", errors.New("X"))),
	)

	unwr := errors.Unwrap(err)
	assertMsg(t, "C: B: A", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "B: A", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "A", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "Z: Y: X", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "Y: X", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "X", unwr)
}

func TestNestedMultierr(t *testing.T) {
	err := multierror.Wrap(
		wrap("C", multierror.Wrap(errors.New("B"), errors.New("A"))),
		wrap("Z", wrap("Y", errors.New("X"))),
	)

	unwr := errors.Unwrap(err)
	assertMsg(t, "C: error #1: B; error #2: A", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "error #1: B; error #2: A", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "B", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "A", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "Z: Y: X", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "Y: X", unwr)

	unwr = errors.Unwrap(unwr)
	assertMsg(t, "X", unwr)
}

func wrap(msg string, err error) error {
	return fmt.Errorf(msg+": %w", err)
}

func assertMsg(t *testing.T, msg string, err error) {
	if msg != err.Error() {
		t.Errorf("wrong message on error: expected %q, got %q", msg, err.Error())
	}
}
