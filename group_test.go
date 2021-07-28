package multierror

import (
	"errors"
	"sort"
	"testing"
)

func TestGroupCompletes(t *testing.T) {
	group := new(Group)

	group.Go(func() error {
		return nil
	})

	group.Go(func() error {
		return nil
	})

	_ = group.Wait()
}

func TestGroupGoRunsFunctions(t *testing.T) {
	group := new(Group)
	items := []bool{false, false}

	group.Go(func() error {
		items[0] = true
		return nil
	})

	group.Go(func() error {
		items[1] = true
		return nil
	})

	_ = group.Wait()

	if !items[0] {
		t.Error("goroutine 0 didn't run")
	}
	if !items[1] {
		t.Error("goroutine 1 didn't run")
	}
}

type errorSort []error

func (e errorSort) Len() int           { return len(e) }
func (e errorSort) Less(i, j int) bool { return e[i].Error() < e[j].Error() }
func (e errorSort) Swap(i, j int)      { e[i], e[j] = e[j], e[i] }

func sortErrors(errs []error) {
	sort.Sort(errorSort(errs))
}

func TestGroupCollectsErrors(t *testing.T) {
	group := new(Group)

	err1 := errors.New("first error")
	err2 := errors.New("second error")

	group.Go(func() error { return err1 })
	group.Go(func() error { return err2 })

	errs := group.Wait().Errors()
	sortErrors(errs)

	if !errors.Is(errs[0], err1) || !errors.Is(errs[1], err2) {
		t.Errorf("incorrect errors returned: %+v", errs)
	}
	if errs[0] != err1 || errs[1] != err2 {
		t.Errorf("unexpected error identities")
	}
}

func TestGroupCollectsPanics(t *testing.T) {
	group := new(Group)

	err1 := errors.New("first error")
	err2 := errors.New("second error")

	group.Go(func() error { panic(err1) })
	group.Go(func() error { panic(err2) })

	errs := group.Wait().Errors()
	sortErrors(errs)

	if !errors.Is(errs[0], err1) || !errors.Is(errs[1], err2) {
		t.Errorf("incorrect errors returned: %+v", errs)
	}

	if errs[0].Error()[:16] != "captured panic: " {
		t.Errorf("expected 'captured panic: ', got '%s'", errs[0].Error()[:20])
	}
	if errs[1].Error()[:16] != "captured panic: " {
		t.Errorf("expected 'captured panic: ', got '%s'", errs[1].Error()[:20])
	}
}

func TestGroupSuccessProducesNil(t *testing.T) {
	group := new(Group)

	group.Go(func() error { return nil })
	group.Go(func() error { return nil })
	group.Go(func() error { return nil })

	if err := group.Wait(); err != nil {
		t.Errorf("expected nil error, got %#v\n", err)
	}
}
