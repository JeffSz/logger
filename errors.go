package logger

type LError struct {
	error
	err string
}

func (err LError) Error() string {
	return err.err
}

func NewError(err string) LError {
	return LError{err: err}
}
