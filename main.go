package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var (
	fileIn     string
	fileOut    string
	name       string
	dateFormat string
)

func init() {
	flag.StringVar(&name, "name", "Unknown", "Name of type")
	flag.StringVar(&fileIn, "in", "", "Input file")
	flag.StringVar(&fileOut, "out", "", "Output file")
	flag.StringVar(&dateFormat, "date", time.RFC3339, "Preferred date format")
}

func main() {
	flag.Parse()

	_, err := time.Parse(dateFormat, time.Now().Format(dateFormat))
	if err != nil {
		ExitOnError(err)
	}

	reader, err := GetReader(fileIn)
	if err != nil {
		ExitOnError(err)
	}

	b, err := Goify(reader, name, "test")
	if err != nil {
		ExitOnError(err)
	}

	fmt.Println(string(b))
}

func ExitOnError(err error) {
	fmt.Println(err)
	os.Exit(-1)
}

func GetReader(fileIn string) (io.ReadCloser, error) {
	if fileIn == "" {
		return os.Stdin, nil
	}
	return os.Open(fileIn)
}

func GetWriter(fileOut string) (io.WriteCloser, error) {
	if fileOut == "" {
		return os.Stdout, nil
	}
	return os.Create(fileOut)
}
