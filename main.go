// https://github.com/Ahmad-Magdy/CSV-To-JSON-Converter/blob/master/main.go
package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func main() {
	path := flag.String("path", "./data.csv", "Path of the file")
	flag.Parse()
	fileBytes, fileNPath := ReadCSV(path)
	SaveFile(fileBytes, fileNPath)
	fmt.Println(strings.Repeat("=", 10), "Done", strings.Repeat("=", 10))
}

// ReadCSV to read the content of CSV File
func ReadCSV(path *string) ([]byte, string) {
	csvFile, err := os.Open(*path)

	if err != nil {
		log.Fatal("The file is not found || wrong root")
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	content, _ := reader.ReadAll()

	if len(content) < 1 {
		log.Fatal("Something wrong, the file maybe empty or length of the lines are not the same")
	}

	headersArr := make([]string, 0)
	for _, headE := range content[0] {
		headersArr = append(headersArr, headE)
	}

	//Remove the header row
	content = content[1:]

	var buffer bytes.Buffer
	buffer.WriteString("[")
	for i, d := range content {
		buffer.WriteString("{")
		for j, y := range d {
			buffer.WriteString(`"` + headersArr[j] + `":`)
			if y == "" {
				buffer.WriteString((`""`))
			} else {
				_, fErr := strconv.ParseFloat(y, 32)
				_, bErr := strconv.ParseBool(y)
				// when csv value like as 0xxxx or 123_456, MarshalIndent put out nothings.
				if y[0] != '0' && bytes.IndexAny([]byte{'_'}, y) == -1 && fErr == nil {
					buffer.WriteString(y)
				} else if bErr == nil {
					buffer.WriteString(strings.ToLower(y))
				} else {
					buffer.WriteString((`"` + y + `"`))
				}
			}

			//end of property
			if j < len(d)-1 {
				buffer.WriteString(",")
			}
		}
		//end of object of the array
		buffer.WriteString("}")
		if i < len(content)-1 {
			buffer.WriteString(",")
		}
	}

	buffer.WriteString(`]`)
	rawMessage := json.RawMessage(buffer.String())

	x, _ := json.MarshalIndent(rawMessage, "", "  ")
	if len(x) == 0 {
		log.Fatal("Something wrong, on the way to convert CSV to JSON.")
	}

	newFileName := filepath.Base(*path)
	newFileName = newFileName[0:len(newFileName)-len(filepath.Ext(newFileName))] + ".json"
	r := filepath.Dir(*path)
	return x, filepath.Join(r, newFileName)
}

// SaveFile Will Save the file, magic right?
func SaveFile(myFile []byte, path string) {
	if err := ioutil.WriteFile(path, myFile, os.FileMode(0644)); err != nil {
		panic(err)
	}
}
