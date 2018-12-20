package testing

import (
	"testing"
)

func TestJSONComparation(t *testing.T) {
	type testCase struct {
		name   string
		got    string
		want   string
		errors bool
	}

	tests := []testCase{
		{
			"valid empty json string",
			"{}",
			"{}",
			false,
		},
		{
			"malformed got",
			"{...}",
			"{}",
			true,
		},
		{
			"empty want",
			"{}",
			"",
			true,
		},
		{
			"empty got",
			"",
			"{}",
			true,
		},
		{
			"order insensitive",
			`{"hello":"world","ciao":"mondo"}`,
			`{"ciao":"mondo","hello":"world"}`,
			false,
		},
		{
			"case sensitive",
			`{"lower":"case"}`,
			`{"UPPER":"CASE"}`,
			true,
		},
	}

	for _, tt := range tests {
		if err := CompareJSON(tt.got, tt.want, tt.name, "CompareJSON"); tt.errors != (err != nil) {
			t.Errorf("%q.%s", "testing", tt.name)
		}
	}
}
