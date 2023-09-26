package kernelparameters

import (
	"errors"
	"reflect"
	"sort"
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

	actual, err := New(v)
	if !reflect.DeepEqual(actual, expected) || err != nil {
		t.Fatalf(`New() = %v, %v, expected: %v, nil`, actual, err, expected)
	}
}

func TestFormatKernelParameters(t *testing.T) {
	v := KernelParameters{
		"test1": "",
		"test2": "value",
	}

	expected := []string{
		"test1",
		"test2=value",
	}

	actual := FormatKernelParameters(v)

	// Sort the slices, since the order is not guaranteed and I don't care about it either
	sort.Strings(actual)
	sort.Strings(expected)

	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf(`FormatKernelParameters() = %v, expected: %v`, actual, expected)
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

	actualValue, err := New(v)
	if err == nil || !errors.As(err, &actualErr) || err.Error() != "invalid kernel parameter [this is invalid] provided" {
		t.Fatalf(`New() = %v, %v, expected: %v, %v`, actualValue, actualErr, expectedValue, expectedErr)
	}
}
