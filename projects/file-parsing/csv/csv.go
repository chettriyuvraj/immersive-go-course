package csv

import (
	"encoding/csv"
)

func ReadNextCSVRecord(csvRdr *csv.Reader) (record []string, err error) {
	record, err = csvRdr.Read()
	if err != nil {
		return record, err
	}

	return record, nil
}
