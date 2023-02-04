package kernelparameters

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

type KernelParameters map[string]string

func New(s []string) (KernelParameters, error) {
	kp := make(KernelParameters)

	for _, v := range s {
		err := validateParameter(v)
		if err != nil {
			return kp, err
		}

		splitStr := strings.Split(v, "=")
		if len(splitStr) > 1 {
			// key value parameter, e.g. initrd=initrd.img
			kp[splitStr[0]] = splitStr[1]
		} else {
			// just a value, e.g. noquiet
			kp[splitStr[0]] = ""
		}
	}

	return kp, nil
}

func validateParameter(v string) error {
	pattern := "\\s|=.*="
	matched, err := regexp.MatchString(pattern, v)
	if err != nil {
		return err
	}
	if matched {
		return errors.New(fmt.Sprintf("kernel parameter [%s] contains illegal characters", v))
	}

	return nil
}

func FormatKernelParameters(kp KernelParameters) []string {
	var s []string

	for k, v := range kp {
		if len(v) > 0 {
			s = append(s, fmt.Sprintf("%s=%s", k, v))
		} else {
			s = append(s, k)
		}
	}

	return s
}

// MergeKernelParameters merges two slices of KernelParameter, with values from the second slice taking precedence
func MergeKernelParameters(kp1 map[string]string, kp2 ...map[string]string) map[string]string {
	// Remove duplicate values and overwrite duplicate key-value options
	// Use the first passed slice as the base/result. Loop over the second provided slice.
	// Check if the parameter already exists in the result slice. If it does: key-value parameter, overwrite the value. Value parameter, do nothing.
	for _, kpMap := range kp2 {
		for k, v := range kpMap {
			kp1[k] = v
		}
	}
	return kp1
}
