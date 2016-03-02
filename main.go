package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	//"io/ioutil"
	"os"
	"reflect"
	"strings"
)

var (
	fileIn string
	name   string
)

func init() {
	flag.StringVar(&name, "name", "Unknown", "Name of type")
	flag.StringVar(&fileIn, "in", "", "Input file")
}

func main() {
	flag.Parse()

	reader, err := GetReader(fileIn)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(Goify(reader))
}

// func Goify(reader io.Reader) error {
// 	var m map[string]interface{}
//
// 	dec := json.NewDecoder(reader)
// 	if err := dec.Decode(&m); err != nil {
// 		return fmt.Errorf("Invalid JSON %v", err)
// 	}
//
// 	for k, v := range m {
// 		typ := reflect.TypeOf(v)
//
// 		value := reflect.New(typ)
//         value.
// 		// switch typ {
// 		// case reflect.Map:
// 		// 	fmt.Printf("%s struct {\n", strings.Replace(strings.Title(k), "_", "", 10))
// 		// 	goify(v.(map[string]interface{}))
// 		// 	fmt.Printf("} `json:\"%s\"`\n", k)
// 		// default:
// 		// 	fmt.Printf("%s %s `json:\"%s\"`\n", strings.Replace(strings.Title(k), "_", "", 10), typ, k)
// 		// }
// 	}
// }

func Goify(reader io.Reader) error {
	var m map[string]interface{}

	dec := json.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return fmt.Errorf("Invalid JSON %v", err)
	}

	fmt.Printf("type %s struct {\n", name)
	goify(m)
	fmt.Printf("}\n")

	return nil
}

func goify(m map[string]interface{}) {

	for k, v := range m {
		typ := reflect.TypeOf(v).Kind()
		switch typ {
		case reflect.Map:
			fmt.Printf("%s struct {\n", strings.Title(strings.Replace(k, "_", "", 10)))
			goify(v.(map[string]interface{}))
			fmt.Printf("} `json:\"%s\"`\n", k)
		case reflect.Slice:
			fmt.Println("====>", reflect.ValueOf(v))
			fmt.Printf("%s []%s `json:\"%s\"`\n", strings.Title(strings.Replace(k, "_", "", 10)), typ, k)
		default:
			fmt.Printf("%s %s `json:\"%s\"`\n", strings.Title(strings.Replace(k, "_", "", 10)), typ, k)
		}
	}
}

func GetReader(fileIn string) (io.ReadCloser, error) {
	if fileIn == "" {
		return os.Stdin, nil
	}
	return os.Open(fileIn)
}
