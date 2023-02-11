package server

type HTTPError struct {
	err        error
	StatusCode int
}

func NewHTTPError(err error, code int) HTTPError {
	return HTTPError{
		err:        err,
		StatusCode: code,
	}
}

func (h HTTPError) Error() string {
	return h.err.Error()
}

func (h HTTPError) Unwrap() error {
	return h.err
}
