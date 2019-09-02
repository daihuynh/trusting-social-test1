package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"
)

func man() {
	fmt.Println(`Usage: main {csv input path | required} {csv output path}
		output: default is './output.csv'
		Author - Huynh Phuc Dai
	`)
}

func main() {
	start := time.Now()
	defer func() {
		if err := recover(); err != nil {
			log.Fatalln(err)
		}
	}()

	args := os.Args[1:]

	if len(args) == 0 {
		man()
		os.Exit(0)
	}

	inputFileName := "input.csv"
	outputFileName := "output.csv"

	if len(args) > 0 {
		inputFileName = args[0]
	}
	if len(args) > 1 {
		outputFileName = args[1]
	}

	userPhoneMap, err := buildPhoneMap(inputFileName)
	if err != nil {
		panic(err)
	}

	// Write to file
	if err := writeToFile(outputFileName, *userPhoneMap); err != nil {
		panic(err)
	}

	fmt.Printf("Output: %s \r\n", outputFileName)
	fmt.Printf("Time execute: %s \r\n", time.Since(start))
}

// buildPhoneMap will read row by row to build a hashmap for each phone number
func buildPhoneMap(inputFileName string) (*map[string]map[string]string, error) {
	file, err := os.Open(inputFileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	userPhoneMap := map[string]map[string]string{}
	didReadHeader := false
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if len(record) < 3 {
			return nil, fmt.Errorf("Wrong format at line %v", record)
		}

		if !didReadHeader {
			didReadHeader = true
			continue
		}
		phone := record[0]
		activate := record[1]
		deactivate := record[2]

		if _, ok := userPhoneMap[phone]; !ok {
			userPhoneMap[phone] = map[string]string{}
		}
		userPhoneMap[phone][deactivate] = activate
	}

	return &userPhoneMap, nil
}

func writeToFile(outFileName string, userPhoneMap map[string]map[string]string) (err error) {
	var f *os.File
	f, err = os.Create(outFileName)
	if err != nil {
		return
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	err = writer.Write([]string{"PHONE_NUMBER", "REAL_ACTIVATION_DATE"})
	if err != nil {
		return
	}
	defer writer.Flush()

	resultCount := len(userPhoneMap)
	phoneNumbers := make([]string, 0, resultCount)
	for phoneNumber := range userPhoneMap {
		phoneNumbers = append(phoneNumbers, phoneNumber)
	}

	sort.Strings(phoneNumbers)

	// for phoneNumber, result := range phoneResult {
	for counter := 0; counter < resultCount; counter++ {
		phoneNumber := phoneNumbers[counter]
		result := userPhoneMap[phoneNumber]
		// This is an ordinary situtation
		if latestActivateDate, ok := result[""]; ok {
			activateDate := latestActivateDate
			for len(result[activateDate]) > 0 {
				activateDate = result[activateDate]
			}
			writer.Write([]string{phoneNumber, activateDate})
		} else {
			keys := make([]string, 0, len(result))
			for deactivate := range result {
				keys = append(keys, deactivate)
			}
			sort.Strings(keys)

			activateDate := result[keys[len(keys)-1]]
			for len(result[activateDate]) > 0 {
				activateDate = result[activateDate]
			}
			writer.Write([]string{phoneNumber, activateDate})
		}
		// latestDeactivation := latestDeactivationMap[phoneNumber]
		// if latestDeactivation != nil {
		// 	activateDate := result[*latestDeactivation]
		// 	for len(result[activateDate]) > 0 {
		// 		activateDate = result[activateDate]
		// 	}
		// 	writer.Write([]string{phoneNumber, activateDate})
		// }
	}

	return
}
