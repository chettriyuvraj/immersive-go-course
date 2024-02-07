package ourbytesbuffer

import (
	"bytes"
	"fmt"
	"io"
	"math"
)

const EXPFACTOR = 2

type OurBytesBuffer struct {
	buffer    []byte
	offset    int
	endOffset int
}

func NewBuffer(data []byte) (*OurBytesBuffer, error) {
	b := OurBytesBuffer{buffer: make([]byte, EXPFACTOR*len(data))}
	_, err := b.Write(data)
	if err != nil {
		return nil, fmt.Errorf("func ourbytesbuffer.OurBytesBuffer.NewBuffer: %w", err)
	}
	b.endOffset = len(data)
	return &b, nil
}

func (b *OurBytesBuffer) Read(inputBuffer []byte) (int, error) {
	/* Check if all data already read */
	inputBufferLen := len(inputBuffer)
	if b.offset == b.endOffset {
		return 0, io.EOF
	}

	/* If input buffer cannot accomodate all data */
	if inputBufferLen < b.endOffset-b.offset {
		err := b.overwriteFromStart(inputBuffer, b.offset, inputBufferLen)
		if err != nil {
			return inputBufferLen, fmt.Errorf("func ourbytesbuffer.OurBytesBuffer.Read: %w", err)
		}
		b.offset += inputBufferLen
		return inputBufferLen, nil
	}

	err := b.overwriteFromStart(inputBuffer, b.offset, b.endOffset-b.offset)
	if err != nil {
		return b.endOffset - b.offset, fmt.Errorf("func ourbytesbuffer.OurBytesBuffer.Read: %w", err)
	}
	bytesWritten := b.endOffset - b.offset
	b.offset = b.endOffset
	return bytesWritten, nil
}

func (b *OurBytesBuffer) Bytes() []byte {
	return b.buffer[b.offset:b.endOffset]
}

func (b *OurBytesBuffer) Write(toAppend []byte) (int, error) {
	bufferLen, appendLen := len(b.buffer), len(toAppend)
	freeSpace := bufferLen - b.endOffset

	/* Allot new buffer if current one is unable to accomodate data*/
	if freeSpace < appendLen {
		newBufferLen := (bufferLen + appendLen) * EXPFACTOR
		/* This is our hard limit on buffer length */
		if float64(newBufferLen) > math.Pow(2, 31) {
			return 0, bytes.ErrTooLarge
		}

		newBuffer := make([]byte, newBufferLen)
		copy(newBuffer, b.buffer)
		b.buffer = newBuffer
	}

	for i := 0; i < appendLen; b.endOffset, i = b.endOffset+1, i+1 {
		b.buffer[b.endOffset] = toAppend[i]
	}

	return appendLen, nil
}

func (b *OurBytesBuffer) overwriteFromStart(inputBuffer []byte, offset int, bytesToWrite int) error { // providing explicit offset because we may not always want to use b.offset
	inputBufferLen := len(inputBuffer)

	if offset+bytesToWrite > b.endOffset {
		return fmt.Errorf("func ourbytesbuffer.OurBytesBuffer.overWriteFromStart: offset %d does not exist", offset+bytesToWrite)
	}

	if inputBufferLen < bytesToWrite {
		return fmt.Errorf("func ourbytesbuffer.OurBytesBuffer.overWriteFromStart: input buffer space %d, string to write %s", len(inputBuffer), inputBuffer[offset:offset+bytesToWrite])
	}

	for i := 0; i < bytesToWrite; i, offset = i+1, offset+1 {
		inputBuffer[i] = b.buffer[offset]
	}
	return nil
}

func (b *OurBytesBuffer) ResetOffset() {
	b.offset = 0
}
