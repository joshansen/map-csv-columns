package mapCSVColumns

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ctessum/macreader"
)

//The Converter type hold the field map, the input io.Reader and the output io.Writer
type Converter struct {
	FieldMap map[string]string
	Input    io.Reader
	Output   io.Writer
}

//Creates a new converter.
func NewConverter(FieldMap map[string]string) *Converter {
	return &Converter{FieldMap: FieldMap}
}

//Sets Converter Input
func (c *Converter) SetInput(input io.Reader) {
	c.Input = input
}

//Sets Converter Output
func (c *Converter) SetOutput(output io.Reader) {
	c.Output = output
}

//Sets Converter Input by opening a file with filname.
func (c *Converter) SetInputWithFilename(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	c.Input = file
	return nil
}

//Sets Converter Output by creating a file with filname.
func (c *Converter) SetOutputWithFilename(filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	c.Output = file
	return nil
}

//Converts an input csv file into an output file, mapping column names given in FieldMap
func (c *Converter) Convert() error {
	//Fix csv issue with mac cr with macreader
	r := csv.NewReader(macreader.New(c.Input))

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
	w := csv.NewWriter(c.Output)
	if err := w.WriteAll(convertedCSV); err != nil {
		return err
	}

	return nil
}
