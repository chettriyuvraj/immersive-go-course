package json

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodePlayerDataSetJSONFile(t *testing.T) {
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		want    []PlayerData
		wantErr bool
	}{
		{
			name: "decode valid json in file",
			args: args{
				path: "./assets/json.txt",
			},
			want: []PlayerData{
				{Name: "Aya", HighScore: 10},
				{Name: "Prisha", HighScore: 30},
				{Name: "Charlie", HighScore: -1},
				{Name: "Margot", HighScore: 25},
			},
			wantErr: false,
		},
		{
			name: "invalid filepath",
			want: []PlayerData{},
			args: args{
				path: "./assets/jsoninvalidpath.txt",
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := DecodePlayerDataSetJSONFile(tc.args.path)
			require.Equal(t, err != nil, tc.wantErr)
			require.Equal(t, tc.want, got)
		})
	}
}

func TestDecodePlayerDataSet(t *testing.T) {
	data := strings.NewReader(`[{"name": "Aya", "high_score": 10},{"name": "Prisha", "high_score": 30}]`)
	marshalledData := []PlayerData{{Name: "Aya", HighScore: 10}, {Name: "Prisha", HighScore: 30}}
	got, _ := DecodePlayerDataSet(data)
	require.Equal(t, marshalledData, got)
}
