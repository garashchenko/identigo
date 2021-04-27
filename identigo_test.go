package identigo

import (
	"reflect"
	"testing"
)

type testLine struct {
	in  string
	out string
}

func TestHashString(t *testing.T) {
	result := hashString("test")
	switch {
	case len(result) == 0:
		t.Error("Wrong hash length")
	case len(result) < 32:
		t.Error("Hash was not generated")
	}
}

func TestGetColorFromHash(t *testing.T) {
	var red, blue, green byte
	red = 0
	blue = 100
	green = 200

	colorBytes := []byte{red, blue, green}
	_, err := getColorFromBytes(colorBytes)
	if err != nil {
		t.Error("Error getting color from bytes")
	}

	colorBytes = []byte{red}
	_, err = getColorFromBytes(colorBytes)
	if err == nil {
		t.Error("Not enough bytes for color but no error was raised")
	}
}

func TestGetCellsToFill(t *testing.T) {
	type testLine struct {
		inBytes []byte
		inCount int
		out     []bool
	}

	var testCases = []testLine{
		{[]byte{154, 146, 61, 30, 165, 126, 24, 254, 65, 220, 181, 67, 226, 196, 0, 92, 65, 255, 33, 8, 100, 167, 16, 176, 251, 178, 101, 76, 17},
			15,
			[]bool{false, true, false, true, true, false, false, true, false, true, false, false, true, false, false},
		},
	}

	for _, tt := range testCases {
		out := getCellsToFill(tt.inBytes, tt.inCount)
		if reflect.DeepEqual(out, tt.out) == false {
			t.Error("Wrong fill grid")
		}
	}
}
