package eth

import (
	"encoding/json"
	"testing"
)

func TestTranslateTopicsRecursionFlattensInputs(t *testing.T) {
	tests := [][]interface{}{
		{
			"0x0",
		},
		{
			nil,
		},
		{
			"0x0",
			"0x1",
		},

		{
			[]interface{}{
				"0x0",
				"0x1",
			},
		},

		{
			[]interface{}{
				"0x0",
			},
			[]interface{}{
				"0x1",
			},
		},
	}
	expected := []interface{}{
		[]string{
			"0",
		},
		[]interface{}{
			nil,
		},
		[]string{
			"0",
			"1",
		},
		[]string{
			"0",
			"1",
		},
		[]string{
			"0",
			"1",
		},
	}

	for i := 0; i < len(tests); i++ {
		test := tests[i]
		expect := expected[i]

		output, err := TranslateTopics(test)
		if err != nil {
			t.Error(err)
		} else {
			result, err := json.Marshal(output)
			if err != nil {
				t.Error(err)
				continue
			}
			expectedResult, err := json.Marshal(expect)
			if err != nil {
				t.Fatal(err)
				continue
			}
			if string(result) != string(expectedResult) {
				t.Errorf("%s != %s", string(result), string(expectedResult))
			}
		}
	}
}
