package records

import "errors"

var (
	// ErrUnSupportedKind is a sentinel error thats returned when parsing values with unsupported kind. use errors.Is() to compare
	ErrUnSupportedKind = errors.New("unsupported kind")
)

// A KindErr records failed parsing for values with unsupported underlying kind.
type KindErr struct {
	WrappedErr error
	Message    string
}

func (k KindErr) Error() string {
	return k.Message
}

func (k KindErr) Unwrap() error {
	return k.WrappedErr
}
