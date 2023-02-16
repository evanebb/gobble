package kernelparameters

import (
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
		return NewInvalidParameterError(v)
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

// MergeKernelParameters merges multiple KernelParameter maps, with values from maps being overwritten in the order that they're passed into the function
func MergeKernelParameters(kp1 KernelParameters, kp2 ...KernelParameters) KernelParameters {
	for _, kpMap := range kp2 {
		for k, v := range kpMap {
			kp1[k] = v
		}
	}
	return kp1
}
