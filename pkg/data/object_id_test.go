package data_test

import (
	"fmt"
	"testing"

	"github.com/shuvava/treehub/data"
)

func TestObjectIDValidate(t *testing.T) {
	cases := []struct {
		ID          string
		ExpectError bool
	}{
		{"some_invalid_str", true},
		{"aec070645.some-type", true},
		{"aec070645fe53ee3b3763059376134f058cc337247c978add178b6ccdfb0019f.commit", false},
	}
	for _, test := range cases {
		exStr := "invalid"
		if test.ExpectError != true {
			exStr = "valid"
		}
		name := fmt.Sprintf("objectId '%s' is %s", test.ID, exStr)
		t.Run(name, func(t *testing.T) {
			objID := data.ObjectID(test.ID)
			got := objID.Validate()
			if (got == nil && test.ExpectError) ||
				(got != nil && !test.ExpectError) {
				t.Errorf("for %s objectId got error '%v'", exStr, got)
			}
		})
	}
}

func TestObjectIDPath(t *testing.T) {
	cases := []struct {
		name   string
		parent string
		value  string
		want   string
	}{
		{"Should return split by first two characters path", "/", "abcd", "/ab/cd"},
		{"Should work with none root parent", "/ttt", "abcd", "/ttt/ab/cd"},
		{"Should work with trailing slash", "/ttt/", "abcd", "/ttt/ab/cd"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			objID := data.ObjectID(test.value)
			got := objID.Path(test.parent)
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}

func TestObjectIDFilename(t *testing.T) {
	cases := []struct {
		name  string
		value string
		want  string
	}{
		{"Should return split by first two characters", "abcd", "cd"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			objID := data.ObjectID(test.value)
			got := objID.Filename()
			if got != test.want {
				t.Errorf("got %v, want %v", got, test.want)
			}
		})
	}
}
