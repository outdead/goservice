package multierror

import (
	"errors"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("", func(t *testing.T) {
		err1 := errors.New("error 1")
		err2 := errors.New("error 2")
		err3 := errors.New("error 3")
		err4 := errors.New("error 4")

		errs1 := []error{err1, err2}

		multierr := New(errs1...)
		multierr.Append(err3)
		multierr.Append(err4)

		if multierr.Len() != 4 {
			t.Fatalf("got len %d, want len %d", multierr.Len(), 4)
		}

		want := err1.Error() + ", " + err2.Error() + ", " + err3.Error() + ", " + err4.Error()

		if multierr.Error() != want {
			t.Fatalf("got %q, want %q", multierr.Error(), want)
		}
	})
}
