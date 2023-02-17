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
		Distro:           uuid.Nil,
		KernelParameters: kernelparameters.KernelParameters{},
	}

	actual := New(uuid.Nil, "TestProfile", "", uuid.Nil, kernelparameters.KernelParameters{})
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf(`New() = %v, expected: %v`, actual, expected)
	}
}
