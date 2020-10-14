package cmd

import (
	"reflect"
	"testing"

	goregistry "github.com/micro/micro/v3/service/registry"
)

type parseCase struct {
	args     []string
	values   *goregistry.Value
	expected map[string]interface{}
}

func TestDynamicFlagParsing(t *testing.T) {
	cases := []parseCase{
		{
			args: []string{"--ss=a,b"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "ss",
						Type: "[]string",
					},
				},
			},
			expected: map[string]interface{}{
				"ss": []interface{}{"a", "b"},
			},
		},
		{
			args: []string{"--ss=a", "--ss=b"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "ss",
						Type: "[]string",
					},
				},
			},
			expected: map[string]interface{}{
				"ss": []interface{}{"a", "b"},
			},
		},
		{
			args: []string{"--bs=true,false"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "bs",
						Type: "[]bool",
					},
				},
			},
			expected: map[string]interface{}{
				"bs": []interface{}{true, false},
			},
		},
		{
			args: []string{"--bs=true", "--bs=false"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "bs",
						Type: "[]bool",
					},
				},
			},
			expected: map[string]interface{}{
				"bs": []interface{}{true, false},
			},
		},
		{
			args: []string{"--is=10,20"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "is",
						Type: "[]int32",
					},
				},
			},
			expected: map[string]interface{}{
				"is": []interface{}{int32(10), int32(20)},
			},
		},
		{
			args: []string{"--is=10", "--is=20"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "is",
						Type: "[]int32",
					},
				},
			},
			expected: map[string]interface{}{
				"is": []interface{}{int32(10), int32(20)},
			},
		},
		{
			args: []string{"--is=10,20"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "is",
						Type: "[]int64",
					},
				},
			},
			expected: map[string]interface{}{
				"is": []interface{}{int64(10), int64(20)},
			},
		},
		{
			args: []string{"--is=10", "--is=20"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "is",
						Type: "[]int64",
					},
				},
			},
			expected: map[string]interface{}{
				"is": []interface{}{int64(10), int64(20)},
			},
		},
		{
			args: []string{"--fs=10.1,20.2"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "fs",
						Type: "[]float64",
					},
				},
			},
			expected: map[string]interface{}{
				"fs": []interface{}{float64(10.1), float64(20.2)},
			},
		},
		{
			args: []string{"--fs=10.1", "--fs=20.2"},
			values: &goregistry.Value{
				Values: []*goregistry.Value{
					{
						Name: "fs",
						Type: "[]float64",
					},
				},
			},
			expected: map[string]interface{}{
				"fs": []interface{}{float64(10.1), float64(20.2)},
			},
		},
	}
	for _, c := range cases {
		_, flags, err := splitCmdArgs(c.args)
		if err != nil {
			t.Fatal(err)
		}
		req, err := flagsToRequest(flags, c.values)
		if err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(c.expected, req) {
			t.Fatalf("Expected %v, got %v", c.expected, req)
		}
	}
}
