// 				02/21/26 			//
//   	qiwi@roundinternet.money 	//

package pb

import (
	"fmt"
)

// reads byte values between " <--> " and interpets them as an int64
// returns bytes consumed as int
func ReadInt64At(b []byte, i int, dilineator byte) (int, int64, error) {
	if b[i]-'0' >= 10 {
		return 0, 0, fmt.Errorf(`expected < 10 at %d (got %q)`, i, string(b[i]))
	}

	// unsafe output buffer. length assert if using.
	var buffer int64
	setAt := func(x byte) {
		buffer = buffer*10 + int64(x-'0')
	}

	start := i
	for ; i < len(b); i++ {
		if b[i] == dilineator {
			setAt(b[i+1])
			return i + 1, buffer, nil
		}
		setAt(b[i])
	}

	return 0, 0, fmt.Errorf("unterminated string starting at %d", start-1)
}

// reads byte values between " <--> " and interpets them as a string
// returns bytes consumed as int
func ReadStringAt(b []byte, i int) (int, string, error) {
	if b[i] != '"' {
		return 0, "", fmt.Errorf(`expected '"' at %d (got %q)`, i, b[i])
	}
	i++
	start := i
	for ; i < len(b); i++ {
		if b[i] == '"' {
			return i + 1, string(b[start:i]), nil
		}
	}
	return 0, "", fmt.Errorf("unterminated string starting at %d", start-1)
}

// reads byte values between [ <--> ] and interpets them as two strings of a protobuf DiffLevel
// returns bytes consumed as int
func decodeDiffLevelMsg(b []byte, out *DiffLevel) (int, error) {
	var i int
	if b[i] != '[' {
		return 0, fmt.Errorf(`expected '"' at %d (got %q)`, i, b[i])
	}
	i++

	// price
	var err error
	var s string // to float or not to float
	i, s, err = ReadStringAt(b, i)
	if err != nil {
		return 0, err
	}
	out.Price = s

	if b[i] != ',' {
		return 0, fmt.Errorf("expected ',' after price at %d (got %q)", i, b[i])
	}
	i++ // consume ,

	// size
	i, s, err = ReadStringAt(b, i)
	if err != nil {
		return 0, err
	}
	out.Size = s

	i++ // consume ]

	return i, nil
}

// reads []... values between [ <--> ] and interpets them as levels of a protobuf BookDiff
// returns bytes consumed as int
func (diff *BookDiff) DecodeDiffSideMsg(b []byte, s bool) (int, error) {
	var buffer = make([]*DiffLevel, 0)
	var i int
	if b[i] != '[' { // consume [
		return 0, fmt.Errorf(`expected '"' at %d (got %q)`, i, b[i])
	}
	i++

	for ; i < len(b); i++ {
		switch b[i] {
		case '[':
			var end int
			var err error

			level := &DiffLevel{}
			if end, err = decodeDiffLevelMsg(b[i:], level); err != nil {
				panic(err)
			}
			buffer = append(buffer, level)

			i += end
			i -= 1 // unconsume ]
		case ',':
		case ']':
			if s {
				diff.Asks = buffer
			} else {
				diff.Bids = buffer
			}

			return i + 1, nil // consume ]
		default:
			return 0, fmt.Errorf("Unexpected token: %c in bytes: %s", b[i], string(b[i:]))
		}
	}
	return 0, fmt.Errorf("Buffer full or invalit bytes: %v | %v", b, b[i:])
}
