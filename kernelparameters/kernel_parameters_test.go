package kernelparameters

import (
	"reflect"
	"testing"
)

func TestNewKernelParameters(t *testing.T) {
	v := []string{"test1", "test2=value"}

	expected := KernelParameters{
		"test1": "",
		"test2": "value",
	}

	kp, err := New(v)
	if !reflect.DeepEqual(kp, expected) || err != nil {
		t.Fatalf(`New([]string{"test1 test2=value"}) = %v, %v, expected: %v, nil`, kp, err, expected)
	}
}
