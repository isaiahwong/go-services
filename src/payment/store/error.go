package store

type idError struct {
	s string
}

func (e *idError) Error() string {
	return e.s
}

type connectError struct {
	s string
}

func (e *connectError) Error() string {
	return e.s
}
