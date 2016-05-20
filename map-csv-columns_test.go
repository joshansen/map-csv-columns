package mapCSVColumns

import (
	"testing"
	"strings"
	"bytes"
	"reflect"
)

func TestConvert(t *testing.T) {
	m := map[string]string{
		"column1": "newColumn1",
		"column2": "newColumn2",
	}

	c := NewConverter(m)
	compareConverter := &Converter{fieldMap: m}
	if !reflect.DeepEqual(c, compareConverter) {
		t.Errorf("func NewConverter did not behave as expected:\n\tExpected %+v\n\tRecieved %+v", c, compareConverter)
		return
	}

	inCsv := `column1,column2,column3,column4
"value1,1","value1,2","value1,3","value1,4"
"value2,1","value2,2","value2,3","value2,4"`
	inputReader := strings.NewReader(inCsv)

	c.SetInput(inputReader)
	compareConverter = &Converter{fieldMap: m, input: inputReader}
	if !reflect.DeepEqual(c, compareConverter) {
		t.Errorf("func Converter.SetInput did not behave as expected:\n\tExpected %+v\n\tRecieved %+v", c, compareConverter)
		return
	}

	buff := new(bytes.Buffer)

	c.SetOutput(buff)
	compareConverter = &Converter{fieldMap: m, input: inputReader, output: new(bytes.Buffer)}
	if !reflect.DeepEqual(c, compareConverter) {
		t.Errorf("func Converter.SetOutput did not behave as expected:\n\tExpected %+v\n\tRecieved %+v", c, compareConverter)
		return
	}

	if err := c.Convert(); err != nil {
		t.Log(err)
	}

	expectedCsvOption1 := `newColumn1,newColumn2
"value1,1","value1,2"
"value2,1","value2,2"
`
	expectedCsvOption2 := `newColumn2,newColumn1
"value1,2","value1,1"
"value2,2","value2,1"
`
	//This test will fail depending on sort order of the map above
	convertedCsv := buff.String()
	if !(expectedCsvOption1 == convertedCsv || expectedCsvOption2 == convertedCsv){
		t.Errorf("Failed Conversion:\n\tExpected %q or %q\n\tRecieved %q", expectedCsvOption1, expectedCsvOption2, convertedCsv)
	}

	invalidFieldMap := map[string]string{
		"column1": "newColumn1",
		"column5": "newColumn5",
	}
	c.fieldMap = invalidFieldMap
	inputReader.Seek(0,0)
	err := c.Convert();

	expectedError := "could not find the following column name(s) in the csv header:\n\tcolumn5\n"

	if err.Error() !=  expectedError{
		t.Errorf("Unexpected Error:\n\tExpected %q\n\tRecieved %q", expectedError, err.Error())
	}
}