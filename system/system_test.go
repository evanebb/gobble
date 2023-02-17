package system

import (
	"github.com/evanebb/gobble/kernelparameters"
	"github.com/google/uuid"
	"net"
	"reflect"
	"testing"
)

func TestNewProfile(t *testing.T) {
	mac, err := net.ParseMAC("11:22:33:44:55:66")
	if err != nil {
		t.Fatalf(`New(): failed to parse MAC address`)
	}

	expected := System{
		Id:               uuid.Nil,
		Name:             "TestSystem",
		Description:      "",
		Profile:          uuid.Nil,
		Mac:              mac,
		KernelParameters: kernelparameters.KernelParameters{},
	}

	actual := New(uuid.Nil, "TestSystem", "", uuid.Nil, mac, kernelparameters.KernelParameters{})
	if !reflect.DeepEqual(actual, expected) {
		t.Fatalf(`New() = %v, expected: %v`, actual, expected)
	}
}
