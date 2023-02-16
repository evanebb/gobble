package kernelparameters

import "fmt"

type InvalidParameterError struct {
	Parameter string
}

func NewInvalidParameterError(p string) *InvalidParameterError {
	return &InvalidParameterError{
		Parameter: p,
	}
}

func (i *InvalidParameterError) Error() string {
	return fmt.Sprintf("invalid kernel parameter [%s] provided", i.Parameter)
}
