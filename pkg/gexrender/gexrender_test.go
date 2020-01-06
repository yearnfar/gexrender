package gexrender

import (
	"testing"
)

func TestRender(t *testing.T) {
	options := []Option{
		WithConfigFile("./testdata/myjob.json"),
		WithBinary("aerender"),
		WithSkipCleanup(),
		WithNoLicense(),
	}

	err := Render(options...)
	if err != nil {
		t.Fatal(err)
	}
}
