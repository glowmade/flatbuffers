package flatbuffers

import (
	"fmt"
)

// FlatBuffer is the interface that represents a flatbuffer.
type FlatBuffer interface {
	Table() Table
	Init(buf []byte, i UOffsetT)
}

// GetRootAs is a generic helper to initialize a FlatBuffer with the provided buffer bytes and its data offset.
func GetRootAs(buf []byte, offset UOffsetT, fb FlatBuffer) {
	n := GetUOffsetT(buf[offset:])
	fb.Init(buf, n+offset)
}

// carefully cut out a 4-byte string from the byte buffer, breaking if we roll outside of ASCII range or hit a null terminator
func extract4ByteStringID(c []byte) string {
	n := -1
	for i, b := range c {
		if b == 0 || b < 32 || b > 126 {
			break
		}
		n = i
		if n >= 3 {
			break
		}
	}
	return string(c[:n+1])
}

// GetIdentifier extracts the embedded type/file identifier from a buffer, if it is present (ie FinishWithID was called to embed)
func GetIdentifier(buf []byte) (string, error) {
	if len(buf) < SizeUOffsetT+SizeIdentifier {
		return "", fmt.Errorf("buffer too small (%d) to contain identifier", len(buf))
	}

	idString := extract4ByteStringID(buf[SizeUOffsetT : SizeUOffsetT+SizeIdentifier])

	if len(idString) != 4 {
		return "", fmt.Errorf("string id too small (%d), expected 4 bytes", len(idString))
	}

	return idString, nil
}
