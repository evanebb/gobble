package system

import (
	"github.com/evanebb/gobble/kernelparameters"
	"testing"
)

func TestPxeConfig_Render(t *testing.T) {
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

func TestRenderNotFound(t *testing.T) {
	expected := `#!ipxe

echo No matching profile found for system!
prompt --timeout 30000 Press any key or wait 30 seconds to continue ||

`

	actual := RenderNotFound()
	if actual != expected {
		t.Fatalf("RenderNotFound() = %v, expected %v", actual, expected)
	}
}
