package internal

type HttpError struct {
	Err  error
	Code int
}

func (hErr HttpError) Error() string {
	return hErr.Err.Error()
}
