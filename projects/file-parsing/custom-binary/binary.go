package binary

import (
	"bufio"
	"fmt"
	"os"
)

type PlayerData struct { // repeated in all packages bring to one place
	Name      string `json:"name"`
	HighScore int    `json:"high_score"`
}

type PlayerDataBinary struct {
	isBigEndian bool          // if not then it's little endian
	reader      *bufio.Reader // reader with it's starting endianess byte already read
}

func NewPlayerDataBinary(path string) (*PlayerDataBinary, error) {
	var pdb PlayerDataBinary

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r := bufio.NewReader(f)
	pdb.reader = r
	pdb.isBigEndian, err = isBigEndian(r)
	if err != nil {
		return nil, err
	}

	return &pdb, nil
}

func isBigEndian(r *bufio.Reader) (isBE bool, err error) { // r must be untouched i.e at the first byte, otherwise won't be able to find endianess mark
	b1, err := r.ReadByte()
	if err != nil {
		return false, fmt.Errorf("error reading first byte: %w", err)
	}
	b2, err := r.ReadByte()
	if err != nil {
		return false, fmt.Errorf("error reading second byte: %w", err)
	}

	if b1 == 0xFE && b2 == 0xFF {
		isBE = true
		return isBE, nil
	}

	if !isBE && b1 != 0xFF && b2 != 0xFE {
		return isBE, fmt.Errorf("no endianess mark found")
	}

	return isBE, nil
}

func GetNextBinaryRecord(d PlayerDataBinary) ([]byte, error) {
	score := []byte{}

	for i := 0; i < 4; i++ {
		scoreByte, err := d.reader.ReadByte()
		if err != nil {
			return []byte{}, err
		}
		score = append(score, scoreByte)
	}

	name, err := d.reader.ReadBytes(0x00)
	if err != nil {
		return name, err
	}

	return append(score, name...), nil
}
