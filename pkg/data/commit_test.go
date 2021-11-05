package data_test

import (
	"fmt"
	"testing"

	"github.com/shuvava/treehub/pkg/data"
)

func TestCommitValidate(t *testing.T) {
	cases := []struct {
		ID          string
		ExpectError bool
	}{
		{"some_invalid_str", true},
		{"aec070645e", true},
		{"aec070645fe53ee3b3763059376134f058cc337247c978add178b6ccdfb0019f", false},
	}
	for _, test := range cases {
		exStr := "invalid"
		if test.ExpectError != true {
			exStr = "valid"
		}
		name := fmt.Sprintf("objectId '%s' is %s", test.ID, exStr)
		t.Run(name, func(t *testing.T) {
			objID := data.Commit(test.ID)
			got := objID.Validate()
			if (got == nil && test.ExpectError) ||
				(got != nil && !test.ExpectError) {
				t.Errorf("for %s objectId got error '%v'", exStr, got)
			}
		})
	}
}
