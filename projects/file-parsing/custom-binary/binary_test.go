package binary

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetNextBinaryRecord(t *testing.T) {
	record1 := []byte{0x00, 0x00, 0x00, 0x0a, 0x41, 0x79, 0x61, 0x00}                   // <Score>Aya0x00
	record2 := []byte{0x00, 0x00, 0x00, 0x1e, 0x50, 0x72, 0x69, 0x73, 0x68, 0x61, 0x00} // <Score>Prisha0x00
	b := []byte{0xFE, 0xFF}                                                             // Big Endian
	b = append(b, record1...)
	b = append(b, record2...)

	r := bufio.NewReader(bytes.NewReader(b))
	isBE, _ := isBigEndian(r)
	pdataBinary := PlayerDataBinary{reader: r, isBigEndian: isBE}

	record, _ := GetNextBinaryRecord(pdataBinary)
	require.Equal(t, record, record1)
	record, _ = GetNextBinaryRecord(pdataBinary)
	require.Equal(t, record, record2)
	_, err := GetNextBinaryRecord(pdataBinary)
	require.ErrorIs(t, err, io.EOF)
}

func Test_isBigEndian(t *testing.T) {
	type args struct {
		r *bufio.Reader
	}
	tests := []struct {
		name     string
		args     args
		wantIsBE bool
		wantErr  bool
		err      error
	}{
		{
			name: "Big Endian",
			args: args{
				r: bufio.NewReader(bytes.NewReader([]byte{0xFE, 0xFF})),
			},
			wantIsBE: true,
		},
		{
			name: "Litcle Endian",
			args: args{
				r: bufio.NewReader(bytes.NewReader([]byte{0xFF, 0xFE})),
			},
			wantIsBE: false,
		},
		{
			name: "Error",
			args: args{
				r: bufio.NewReader(bytes.NewReader([]byte{0xCC, 0xAC, 0xFF, 0xFE})),
			},
			wantErr: true,
			err:     fmt.Errorf("no endianess mark found"),
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			gotIsBE, err := isBigEndian(tc.args.r)
			if tc.wantErr {
				require.Errorf(t, err, tc.err.Error()) // Simply comparing error message - different instances of fmt.Errorf fail ErrorIs test
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantIsBE, gotIsBE)
			}
		})
	}
}
