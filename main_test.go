package main

import (
	"os"
	"testing"
)

func TestFileNotFound(t *testing.T) {
	_, err := buildPhoneMap("")
	if _, ok := err.(*os.PathError); !ok {
		t.Errorf("Unexpected error: %v", err)
	}
}

func TestOrdinaryFlow(t *testing.T) {
	inputFileName := "test_ordinary.csv"
	outputFileName := "test_ordinary.out.csv"

	userPhoneMap, err := buildPhoneMap(inputFileName)
	if err != nil {
		t.Error(err)
	}

	if err := writeToFile(outputFileName, *userPhoneMap); err != nil {
		t.Error(err)
	}
}

func TestBigFile(t *testing.T) {
	inputFileName := "test_big_file.csv"
	outputFileName := "test_big_file.out.csv"

	userPhoneMap, err := buildPhoneMap(inputFileName)
	if err != nil {
		t.Error(err)
	}

	if err := writeToFile(outputFileName, *userPhoneMap); err != nil {
		t.Error(err)
	}
}

func BenchmarkBigFile(b *testing.B) {
	b.ReportAllocs()
	for counter := 0; counter < b.N; counter++ {
		inputFileName := "test_big_file.csv"
		outputFileName := "test_big_file.out.csv"

		userPhoneMap, err := buildPhoneMap(inputFileName)
		if err != nil {
			b.Error(err)
		}

		if err := writeToFile(outputFileName, *userPhoneMap); err != nil {
			b.Error(err)
		}
	}
}
