package csvUtil

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

//CsvWrapper wraps the csv package with extra functions.
type CsvWrapper struct {
	FileHandle *os.File
	Reader     *csv.Reader
	hdr        []string
	itm        []map[string]string
	includeHdr bool
}

//Headers returns the loaded headers
func (pc CsvWrapper) Headers() []string {
	return pc.hdr
}

func (pc *CsvWrapper) SetHeaders(headersArray []string) ([]string, error) {
	itmLen := len(pc.Items()[0])
	if len(headersArray) != itmLen {
		return []string{}, errors.New("Headers length should match Items Length")
	}
	pc.hdr = headersArray
	return pc.Headers(), nil

}

//Items returns all the loaded rows
func (pc CsvWrapper) Items() []map[string]string {
	return pc.itm
}

//Load the file and create a new csvReader
func (pc *CsvWrapper) Load(file string, withHeaders bool) {
	pc.FileHandle, _ = os.Open(file)
	pc.Reader = csv.NewReader(pc.FileHandle)
	pc.Reader.FieldsPerRecord = -1
	pc.includeHdr = withHeaders
	// pc.Parse() // set headers and items
}

//Close all the handles
func (pc CsvWrapper) Close() {
	pc.FileHandle.Close()
}

//Parse will set the headers and items from the loaded csvFile through the
func (pc *CsvWrapper) Parse() {
	rawData, err := pc.Reader.Read()

	if pc.includeHdr {
		for x := range rawData {
			rawData[x] = normalizeString(rawData[x])
		}
		pc.hdr = rawData
		rawData, err = pc.Reader.Read() // header is considered as an extra raw, move the pointer forward.
	} else {
		for x := 0; x < len(rawData); x++ {
			item := fmt.Sprintf("%v%v", "item", x)
			pc.hdr = append(pc.hdr, item)
		}
		//when Seeking to the begining of the file i get an error.
		//so ive rearranged the Reads.
		// pc.FileHandle.Seek(0, os.SEEK_SET) //going back to the beginning of file.
	}
	for {
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err)
			break
		}
		tempRes := make(map[string]string)
		for z, a := range pc.hdr {
			tempRes[a] = rawData[z]
		}
		pc.itm = append(pc.itm, tempRes)
		rawData, err = pc.Reader.Read()

	}
}

//GetElements returns all the elements of the specified headers
func (pc CsvWrapper) GetElements(names ...string) (result map[string][]string, err map[string]error) {
	result = map[string][]string{}
	err = map[string]error{}
	for _, rows := range pc.itm {
		for _, name := range names {
			if i, ok := rows[name]; ok {
				result[name] = append(result[name], i)
			} else {
				err[name] = errors.New("Key not found")
			}
		}

	}
	return
}

//FilterElements returns all the rows with only the specified fields
func (pc CsvWrapper) FilterElements(names ...string) (results []map[string]string, err map[string]error) {
	results = []map[string]string{}
	err = map[string]error{}
	for _, rows := range pc.itm {
		tempMap := make(map[string]string, len(rows))
		for _, name := range names {
			if i, ok := rows[name]; ok {
				tempMap[name] = i
			} else {
				err[name] = errors.New("Key not found")
			}
		}
		results = append(results, tempMap)
	}
	return
}

//normalizeString convert a string to a snake case string
func normalizeString(string string) string {
	string = strings.ToLower(strings.Replace(strings.Trim(string, " "), " ", "_", -1))
	return string

}
