package kernelparameters

import (
	"fmt"
	"regexp"
	"strings"
)

type KernelParameters map[string]string

// String returns the string representation of the KernelParameters, e.g. 'initrd=initrd quiet splash'
func (k KernelParameters) String() string {
	return strings.Join(k.StringSlice(), " ")
}

// StringSlice returns the set of KernelParameters as a slice of strings, e.g. ["initrd=initrd", "quiet", "splash"]
func (k KernelParameters) StringSlice() []string {
	var s []string

	for k, v := range k {
		if len(v) > 0 {
			s = append(s, fmt.Sprintf("%s=%s", k, v))
		} else {
			s = append(s, k)
		}
	}

	return s
}

// ParseString will parse and validate s into a set of KernelParameters, and return an error if it encounters an invalid kernel parameter.
func ParseString(s string) (KernelParameters, error) {
	return ParseStringSlice(strings.Split(s, " "))
}

// ParseStringSlice will parse and validate s into a set of KernelParameters, and return an error if it encounters an invalid kernel parameter.
func ParseStringSlice(s []string) (KernelParameters, error) {
	kp := make(KernelParameters)

	for _, v := range s {
		err := validateParameter(v)
		if err != nil {
			return kp, err
		}

		splitStr := strings.Split(v, "=")
		if len(splitStr) > 1 {
			// key value parameter, e.g. 'initrd=initrd'
			kp[splitStr[0]] = splitStr[1]
		} else {
			// just a value, e.g. 'quiet'
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
		return NewInvalidParameterError(v)
	}

	return nil
}

// MergeKernelParameters merges multiple sets of KernelParameters, with values from maps being overwritten in the order that they're passed into the function
func MergeKernelParameters(kp1 KernelParameters, kp2 ...KernelParameters) KernelParameters {
	for _, kpMap := range kp2 {
		for k, v := range kpMap {
			kp1[k] = v
		}
	}
	return kp1
}
