package bytesbuffer

import (
	"bytes"
	"testing"
)

func TestBufferBytes(t *testing.T) {
	tc := []byte("Hey")
	tc2 := []byte("Dude!")
	tcplustc2 := append(tc, tc2...)
	b := bytes.NewBuffer(tc)

	got := b.Bytes()
	if !bytes.Equal(tc, got) {
		t.Errorf("error: func: bytes.Buffer.Bytes: want %v got %v", string(tc), string(got))
	}

	_, err := b.Write(tc2)
	if err != nil {
		t.Errorf("error: func: bytes.Buffer.Write: %v", err)
	}

	got = b.Bytes()
	if !bytes.Equal(tcplustc2, got) {
		t.Errorf("error: func: bytes.Buffer.Bytes: want %v got %v", string(append(tc, tc2...)), string(got))
	}

}
