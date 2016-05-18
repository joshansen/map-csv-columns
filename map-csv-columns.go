package mapCSVColumns

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ctessum/macreader"
)

type Converter struct {
	FieldMap  map[string]string
	ReadFile  io.Reader
	WriteFile io.Writer
}

//Creates a new converter taking two file names and a field map
func NewConverter(readFile, writeFile string, FieldMap map[string]string) (*Converter, error) {
	c := &Converter{FieldMap: FieldMap}
	if err := c.setReadFile(readFile); err != nil {
		return nil, err
	}
	if err := c.setWriteFile(writeFile); err != nil {
		return nil, err
	}
	return c, nil
}

//Sets ReadFile on converter given a file name.
func (c *Converter) setReadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	c.ReadFile = macreader.New(file)
	return nil
}

//Sets WriteFile on converter given a file name.
func (c *Converter) setWriteFile(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	c.WriteFile = file
	return nil
}

//Converts an input csv file into an output file, mapping column names given in FieldMap
func (c *Converter) Convert() error {
	r := csv.NewReader(c.ReadFile)

	//Read the first row
	firstRow, err := r.Read()
	if err != nil {
		return err
	}

	var (
		missingFieldList       []string
		convertedColumnIndices []int
		convertedCSV           [][]string
		convertedRow           []string
	)

	//Loop through map looking for each key in the first row.
	//If there are any missing fields, return an error with all missing fields.
	//Record all columns that have been changed, delete other columns.
	for inValue, outValue := range c.FieldMap {
		for columnIndex, columnName := range firstRow {
			if columnName == inValue {
				convertedRow = append(convertedRow, outValue)
				convertedColumnIndices = append(convertedColumnIndices, columnIndex)
				break
			}
			//Last pass through loop
			if columnIndex == len(firstRow)-1 {
				missingFieldList = append(missingFieldList, inValue)
			}
		}
	}

	//if fields in FieldMap weren't present, return the fields' name
	if len(missingFieldList) > 0 {
		return fmt.Errorf("could not find the following column name(s) in the csv header.:\n\t%s\n", strings.Join(missingFieldList, "\n\t"))
	}

	//append the first row
	convertedCSV = append(convertedCSV, convertedRow)
	for {

		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		//zero out converted row
		convertedRow = []string{}
		//move record values into converted row
		for _, columnIndex := range convertedColumnIndices {
			convertedRow = append(convertedRow, record[columnIndex])
		}

		//append the converted row
		convertedCSV = append(convertedCSV, convertedRow)
	}

	//write new csv file
	w := csv.NewWriter(c.WriteFile)
	if err := w.WriteAll(convertedCSV); err != nil {
		return err
	}

	return nil
}
