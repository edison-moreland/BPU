package netlist

import (
	"embed"
	"path"
	"strings"
	"testing"
)

//go:embed testdata/*.json
var testdataFS embed.FS

func TestUnmarshalNetlist(t *testing.T) {
	type test struct {
		filename string
	}

	tests := []test{}

	// Discover tests
	testFiles, err := testdataFS.ReadDir("testdata")
	if err != nil {
		t.Fatal(err)
		return
	}

	for _, file := range testFiles {
		tests = append(tests, test{path.Join("testdata", file.Name())})
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {

			goldenBytes, err := testdataFS.ReadFile(tt.filename)
			if err != nil {
				t.Fatalf("failed to read golden file: %v", err)
			}

			netlist, err := UnmarshalNetlist(goldenBytes)
			if err != nil {
				t.Errorf("failed to unmarshal netlist: %v", err)
				return
			}

			if !strings.HasPrefix(netlist.Creator, "Yosys") {
				t.Errorf("something bad")
			}
		})
	}
}
