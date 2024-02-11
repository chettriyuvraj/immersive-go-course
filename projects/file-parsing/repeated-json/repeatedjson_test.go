package repeatedjson

import (
	"bufio"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsLineAComment(t *testing.T) {
	b := []byte("# Ignore this line ffs")
	b2 := []byte(`{"name": "Charlie", "high_score": -1}`)
	isComment := IsLineAComment(b)
	require.Equal(t, true, isComment)
	isComment = IsLineAComment(b2)
	require.Equal(t, false, isComment)
}

func TestScanNextLine(t *testing.T) {
	type args struct {
		data      string
		scanCount int // not technically an argument
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
		err     error
	}{
		{
			name: "Test first line scanning",
			args: args{
				data:      "This is the first line\nThis is the second line\nThis is the third line",
				scanCount: 1,
			},
			want:    []byte("This is the first line"),
			wantErr: false,
		},
		{
			name: "Test second line scanning",
			args: args{
				data:      "This is the first line\nThis is the second line\nThis is the third line",
				scanCount: 2,
			},
			want:    []byte("This is the second line"),
			wantErr: false,
		},
		{
			name: "Test EOF scanning",
			args: args{
				data:      "This is the first line\nThis is the second line\nThis is the third line",
				scanCount: 4,
			},
			want:    nil,
			wantErr: true,
			err:     io.EOF,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			sc := bufio.NewScanner(strings.NewReader(tc.args.data))
			for i := 0; i < tc.args.scanCount-1; i++ {
				ScanNextLine(sc)
			}

			got, err := ScanNextLine(sc)
			if tc.wantErr {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
			}
			require.Equal(t, tc.want, got)
		})
	}
}

func TestParseRepeatedJSONFile(t *testing.T) {
	path := "./assets/repeated-json.txt"
	reqdData := []PlayerData{
		{Name: "Aya", HighScore: 10},
		{Name: "Prisha", HighScore: 30},
	}
	parsedData, err := ParseRepeatedJSONFile(path)
	require.NoError(t, err)
	require.Equal(t, reqdData, parsedData)

}
