package records

type recordsErr string

func (m recordsErr) Error() string {
	return (string)(m)
}

const (
	// ErrUnSupportedKind is a sentinel error thats returned when parsing values with unsupported underlying kind.
	// compare with errors.Is()
	ErrUnSupportedKind = recordsErr("unsupported kind")
)

// A KindErr is returned when parsing values with unsupported underlying kind.
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
