package distro

import (
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestNewDistro(t *testing.T) {
	expected := Distro{
		Id:               uuid.Nil,
		Name:             "TestDistro",
		Description:      "",
		Kernel:           "kernel",
		Initrd:           "initrd",
		KernelParameters: kernelparameters.KernelParameters{},
	}

	actual, err := New(uuid.Nil, "TestDistro", "", "kernel", "initrd", kernelparameters.KernelParameters{})
	if err != nil || !reflect.DeepEqual(actual, expected) {
		t.Fatalf(`New() = %v, %v, expected: %v, nil`, actual, err, expected)
	}
}

func TestNewDistroInvalidName(t *testing.T) {
	actual, err := New(uuid.Nil, "invalid name", "", "kernel", "initrd", kernelparameters.KernelParameters{})
	if err == nil {
		t.Fatalf(`Expected New() to return invalid name error, got: %v, %v`, actual, err)
	}
}
