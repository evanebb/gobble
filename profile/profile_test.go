package profile

import (
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"reflect"
	"testing"
)

func TestNewProfile(t *testing.T) {
	expected := Profile{
		Id:               uuid.Nil,
		Name:             "TestProfile",
		Description:      "",
		Kernel:           "kernel",
		Initrd:           "initrd",
		KernelParameters: kernelparameters.KernelParameters{},
	}

	actual, err := New(uuid.Nil, "TestProfile", "", "kernel", "initrd", kernelparameters.KernelParameters{})
	if err != nil || !reflect.DeepEqual(actual, expected) {
		t.Fatalf(`New() = %v, %v, expected: %v, nil`, actual, err, expected)
	}
}

func TestNewProfile_InvalidName(t *testing.T) {
	actual, err := New(uuid.Nil, "invalid name", "", "kernel", "initrd", kernelparameters.KernelParameters{})
	if err == nil {
		t.Fatalf(`Expected New() to return invalid name error, got: %v, %v`, actual, err)
	}
}

func TestNewProfile_EmptyKernel(t *testing.T) {
	actual, err := New(uuid.Nil, "TestProfile", "", "", "initrd", kernelparameters.KernelParameters{})
	if err == nil || err.Error() != "kernel cannot be empty" {
		t.Fatalf(`Expected New() to return empty kernel error, got: %v, %v`, actual, err)
	}
}

func TestNewProfile_InvalidKernel(t *testing.T) {
	actual, err := New(uuid.Nil, "TestProfile", "", "invalid kernel", "initrd", kernelparameters.KernelParameters{})
	if err == nil || err.Error() != "kernel contains illegal characters" {
		t.Fatalf(`Expected New() to return invalid kernel error, got: %v, %v`, actual, err)
	}
}

func TestNewProfile_EmptyInitrd(t *testing.T) {
	actual, err := New(uuid.Nil, "TestProfile", "", "kernel", "", kernelparameters.KernelParameters{})
	if err == nil || err.Error() != "initrd cannot be empty" {
		t.Fatalf(`Expected New() to return invalid initrd error, got: %v, %v`, actual, err)
	}
}

func TestNewProfile_InvalidInitrd(t *testing.T) {
	actual, err := New(uuid.Nil, "TestProfile", "", "kernel", "invalid initrd", kernelparameters.KernelParameters{})
	if err == nil || err.Error() != "initrd contains illegal characters" {
		t.Fatalf(`Expected New() to return invalid initrd error, got: %v, %v`, actual, err)
	}
}
