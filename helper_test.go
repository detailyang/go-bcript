package bscript

import "testing"

func TestReadNBytes(t *testing.T) {
	tests := []struct {
		data   []byte
		nbytes int
		want   int
	}{
		{[]byte{1}, 1, 1},
		{[]byte{1, 0}, 2, 1},
		{[]byte{1, 0, 0, 0}, 4, 1},
	}

	for _, test := range tests {
		n, err := readNBytes(test.data, test.nbytes)
		if err != nil {
			t.Error(err)
		}

		if n != test.want {
			t.Errorf("Expect %d, want %d", test.want, n)
		}
	}
}
