package kernelparameters

import (
	"errors"
	"reflect"
	"testing"
)

func TestNewKernelParameters(t *testing.T) {
	v := []string{
		"test1",
		"test2=value",
	}

	expected := KernelParameters{
		"test1": "",
		"test2": "value",
	}

	actual, err := ParseStringSlice(v)
	if !reflect.DeepEqual(actual, expected) || err != nil {
		t.Fatalf(`New() = %v, %v, expected: %v, nil`, actual, err, expected)
	}
}

func TestMergeKernelParameters(t *testing.T) {
	kp1 := KernelParameters{
		"test1": "",
		"test2": "value",
		"test3": "value2",
	}

	kp2 := KernelParameters{
		"test2": "newvalue",
		"test4": "value3",
	}

	expected := KernelParameters{
		"test1": "",
		"test2": "newvalue",
		"test3": "value2",
		"test4": "value3",
	}

	actual := MergeKernelParameters(kp1, kp2)
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf(`MergeKernelParameters() = %v, expected: %v`, actual, expected)
	}
}

func TestInvalidParameterError(t *testing.T) {
	var actualErr *InvalidParameterError
	expectedErr := NewInvalidParameterError("this is invalid")

	v := []string{
		"this is invalid",
		"test2=value",
	}

	expectedValue := make(KernelParameters)

	actualValue, err := ParseStringSlice(v)
	if err == nil || !errors.As(err, &actualErr) {
		t.Fatalf(`New() = %v, %v, expected: %v, %v`, actualValue, actualErr, expectedValue, expectedErr)
	}
}
