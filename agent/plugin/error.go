package plugin

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

type unit struct{}

// DeprecatedNameErrors contains a set of DeprecatedNameError
type DeprecatedNameErrors struct {
	errs map[DeprecatedNameError]unit
}

// IsEmpty return true if and only if `e` contains no errors
func (e *DeprecatedNameErrors) IsEmpty() bool {
	return e == nil || len(e.errs) == 0
}

// Errors returns the contained set of errors in sorted order
func (e *DeprecatedNameErrors) Errors() []DeprecatedNameError {
	if e == nil {
		return nil
	}

	if e.errs == nil {
		return []DeprecatedNameError{}
	}

	errs := make([]DeprecatedNameError, 0, len(e.errs))
	for err := range e.errs {
		errs = append(errs, err)
	}

	sort.Slice(errs, func(i, j int) bool {
		if errs[i].old == errs[j].old {
			return errs[i].new < errs[j].new
		}
		return errs[i].old < errs[j].old
	})

	return errs
}

// Error returns each contained error on a new line
func (e *DeprecatedNameErrors) Error() string {
	builder := strings.Builder{}
	for i, err := range e.Errors() {
		_, _ = builder.WriteString(err.Error())
		if i < len(e.errs)-1 {
			_, _ = builder.WriteRune('\n')
		}
	}
	return builder.String()
}

// Append adds DeprecatedNameError contained set and returns the reciver.
// Returning the reveiver is necessary to support appending to nil. So this
// should be used just like the builtin `append` function.
func (e *DeprecatedNameErrors) Append(errs ...DeprecatedNameError) *DeprecatedNameErrors {
	if e == nil {
		e = &DeprecatedNameErrors{errs: map[DeprecatedNameError]unit{}}
	} else if e.errs == nil {
		e.errs = map[DeprecatedNameError]unit{}
	}

	for _, err := range errs {
		e.errs[err] = unit{}
	}

	return e
}

// Is returns true if and only if a error that is wrapped in target
// contains the same set of DeprecatedNameError as the receiver.
func (e *DeprecatedNameErrors) Is(target error) bool {
	if e == nil {
		return target == nil
	}

	var targetErr *DeprecatedNameErrors
	if !errors.As(target, &targetErr) {
		return false
	}

	if len(e.errs) != len(targetErr.errs) {
		return false
	}

	for err := range e.errs {
		if _, exists := targetErr.errs[err]; !exists {
			return false
		}
	}

	return true
}

// DeprecatedNameError contains information about environment variable names that
// are deprecated. Both the deprecated name and its replacement are held
type DeprecatedNameError struct {
	old string
	new string
}

func NewDeprecatedNameError(oldName, newName string) DeprecatedNameError {
	return DeprecatedNameError{old: oldName, new: newName}
}

func (e *DeprecatedNameError) Error() string {
	return fmt.Sprintf(" deprecated: %q\nreplacement: %q\n", e.old, e.new)
}

func (e *DeprecatedNameError) Is(target error) bool {
	if e == nil {
		return target == nil
	}

	var targetErr *DeprecatedNameError
	if !errors.As(target, &targetErr) {
		return false
	}

	return e.old == targetErr.old && e.new == targetErr.new
}
