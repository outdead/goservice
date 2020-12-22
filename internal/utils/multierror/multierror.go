package multierror

import "strings"

type Error interface {
	error
	Append(err error)
	Errors() []error
	Len() int
}

type multierr struct {
	errs []error
}

func New(errs ...error) Error {
	return &multierr{
		errs: errs,
	}
}

func (merr *multierr) Error() string {
	errs := make([]string, len(merr.errs))
	for i, err := range merr.errs {
		errs[i] = err.Error()
	}

	return strings.Join(errs, ", ")
}

func (merr *multierr) Append(err error) {
	if err != nil {
		merr.append(err)
	}
}

func (merr *multierr) Len() int {
	return len(merr.errs)
}

func (merr *multierr) Errors() []error {
	return merr.errs
}

func (merr *multierr) append(err error) {
	merr.errs = append(merr.errs, err)
}
