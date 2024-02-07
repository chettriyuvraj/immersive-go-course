package ourbytesbuffer

import (
	"bytes"
	"testing"
)

func TestBufferBytes(t *testing.T) {
	tc := []byte("Hey")
	tc2 := []byte("Dude!")
	tcplustc2 := append(tc, tc2...)
	b, err := NewBuffer(tc)
	if err != nil {
		t.Errorf("%v", err)
	}

	/* Testing if buffer.Bytes() returns the initialized byte slice */
	got := b.Bytes()
	if !bytes.Equal(tc, got) {
		t.Errorf("error: func: bytes.Buffer.Bytes: want %v got %v", string(tc), string(got))
	}

	/* Testing if buffer.Bytes() returns the initialized byte slice + appended byte slice */
	_, err = b.Write(tc2)
	if err != nil {
		t.Errorf("error: func: bytes.Buffer.Write: %v", err)
	}

	got = b.Bytes()
	if !bytes.Equal(tcplustc2, got) {
		t.Errorf("error: func: bytes.Buffer.Bytes: want %v got %v", string(append(tc, tc2...)), string(got))
	}

}

func TestBufferRead(t *testing.T) {
	tc := []byte("HeyGuys!")
	b, err := NewBuffer(tc)
	if err != nil {
		t.Errorf("%v", err)
	}
	smallSlice := make([]byte, len(tc)/2)
	largeSlice := make([]byte, len(tc)+1)

	/* Testing if bytes.Buffer.Read() reads correct number of bytes on a slice smaller than buffer - confirms that only some of the bytes are read */
	n, err := b.Read(smallSlice)
	if err != nil {
		t.Errorf("error: func: bytes.Buffer.Read: %v", err)
	}

	if n > len(smallSlice) {
		t.Errorf("error: func: bytes.Buffer.Read: input buffer len %d, length read %d", len(smallSlice), n)
	}

	if !bytes.Equal(tc[:n], smallSlice) {
		t.Errorf("error: func: bytes.Buffer.Read: expected to read %v, actually read %v", string(tc[:n]), string(smallSlice))
	}

	/* Testing if next bytes are read on the next Read() */
	n2, err := b.Read(smallSlice)
	if err != nil {
		t.Errorf("error: func: bytes.Buffer.Read: %v", err)
	}

	if !bytes.Equal(tc[n:n+n2], smallSlice) {
		t.Errorf("error: func: bytes.Buffer.Read: expected to read %v, actually read %v", string(tc[n:n+n2]), string(smallSlice))
	}

	/* Testing that a large enough slice reads all of the bytes int the buffer */
	b, err = NewBuffer(tc)
	if err != nil {
		t.Errorf("%v", err)
	}
	n, err = b.Read(largeSlice)
	if err != nil {
		t.Errorf("error: func: bytes.Buffer.Read: %v", err)
	}
	if !bytes.Equal(tc, largeSlice[:n]) {
		t.Errorf("error: func: bytes.Buffer.Read: expected to read %v, actually read %v", string(tc), string(largeSlice[:n]))
	}

}
