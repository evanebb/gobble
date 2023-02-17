package system

import (
	"github.com/evanebb/gobble/kernelparameters"
	"testing"
)

func TestRenderPxeConfig(t *testing.T) {
	expected := `#!ipxe

kernel testkernel param1
initrd testinitrd

boot
`

	kp, err := kernelparameters.New([]string{"param1"})
	if err != nil {
		t.Fatalf(`NewPxeConfig(): failed to instantiate KernelParameters, error: %v`, err)
	}

	pxeConfig := NewPxeConfig("testkernel", "testinitrd", kp)
	actual := pxeConfig.Render()
	if actual != expected {
		t.Fatalf("PxeConfig.Render() = %v, expected %v", actual, expected)
	}
}
