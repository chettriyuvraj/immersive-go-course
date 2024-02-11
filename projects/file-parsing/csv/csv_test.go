package csv

import (
	"encoding/csv"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadNextCSVRecord(t *testing.T) {
	type args struct {
		data      string // techincally not an arg
		readCount int    // total number of reads to get to record - technically not an arg
	}
	tests := []struct {
		name       string
		args       args
		wantRecord []string
		wantErr    bool
		err        error
	}{
		{
			name: "Parse CSV title line",
			args: args{
				data: `first_name,last_name,username
"Rob","Pike",rob`,
				readCount: 1,
			},
			wantRecord: []string{"first_name", "last_name", "username"},
		},
		{
			name: "Parse CSV first line",
			args: args{
				data: `first_name,last_name,username
"Rob","Pike",rob`,
				readCount: 2,
			},
			wantRecord: []string{"Rob", "Pike", "rob"},
		},
		{
			name: "Parse CSV EOF",
			args: args{
				data: `first_name,last_name,username
"Rob","Pike",rob`,
				readCount: 3,
			},
			wantErr: true,
			err:     io.EOF,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			csvRdr := csv.NewReader(strings.NewReader(tc.args.data))
			for i := 0; i < tc.args.readCount-1; i++ {
				ReadNextCSVRecord(csvRdr)
			}
			gotRecord, err := ReadNextCSVRecord(csvRdr)
			if tc.wantErr {
				require.ErrorIs(t, err, tc.err)
			} else {
				require.NoError(t, err)
				require.Equal(t, tc.wantRecord, gotRecord)
			}

		})
	}
}
