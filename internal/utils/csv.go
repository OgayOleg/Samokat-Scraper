package utils

import (
	"encoding/csv"
	"os"
)

func SaveToTXT(filename string, rows [][]string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	return writer.WriteAll(rows)
}
