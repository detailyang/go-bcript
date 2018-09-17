package bscript

type Boolean bool

func NewBoolean(d []byte) Boolean {
	rv := false
	for i := range d {
		if d[i] != 0 {
			// Negative 0 is also considered false.
			if i == len(d)-1 && d[i] == 0x80 {
				rv = false
			}

			rv = true
		}
	}

	return Boolean(rv)
}

// Byte vectors are interpreted as Booleans
// where False is represented by any representation of zero and True is represented by any representation of non-zero.
func (b Boolean) Bytes() []byte {
	if b {
		return []byte{1}
	}

	return []byte{0}
}
