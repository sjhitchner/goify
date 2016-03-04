package main

import (
	"encoding/json"
	"fmt"
	//	"go/ast"
	//	"go/parser"
	//	"go/token"
	"bytes"
	"io"
	"log"
	"reflect"
	"sort"
	"strings"
	"time"
)

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
const STRUCT = `
type {{.Name}} struct {
	{{range .Fields}}
		
	{{end}}
}
`
const SLICE = "{{.Name}} []{{.Type}} `json:\"{{.JSON}}\"`"
const PRIMATIVE = "{{.Name}} {{.Type}} `json:\"{{.JSON}}\"`"

var m map[string]interface{}

func Goify(reader io.Reader, structName string, packageName string) ([]byte, error) {
	buf := &bytes.Buffer{}

	var m interface{}
	dec := json.NewDecoder(reader)
	if err := dec.Decode(&m); err != nil {
		return nil, fmt.Errorf("Invalid JSON %v", err)
	}

	fmt.Fprintf(buf, "package %s\n", packageName)
	fmt.Fprintf(buf, "type %s", structName)

	if err := generate(buf, m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func generate(buf *bytes.Buffer, m interface{}) error {

	switch mt := m.(type) {
	case map[string]interface{}:
		return generateStruct(buf, mt, 0)

	case []map[string]interface{}:
		if len(mt) > 0 {
			return generateStruct(buf, mt[0], 0)
		}
		return fmt.Errorf("json array is empty")

	case []interface{}:
		return fmt.Errorf("json array not implemented")

	default:
		return fmt.Errorf("invalid type: %T", mt)
	}
}

func generateStruct(buf *bytes.Buffer, m map[string]interface{}, depth int) error {

	keys := make([]string, 0)
	for key, _ := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := m[key]

		buf.WriteString(strings.Repeat("  ", depth))
		buf.WriteString(jsonToGoName(key))
		buf.WriteString(" ")

		switch vt := value.(type) {
		case []map[string]interface{}:
			log.Println("[]map[string]interface{}")

			buf.WriteString("[]struct {\n")
			if len(vt) > 0 {
				if err := generateStruct(buf, vt[0], depth+1); err != nil {
					return err
				}
			} else {
				return fmt.Errorf("empty json array")
			}
			buf.WriteString("}")

		case map[string]interface{}:
			log.Println("map[string]interface{}")

			buf.WriteString("struct {\n")
			if err := generateStruct(buf, vt, depth+1); err != nil {
				return err
			}
			buf.WriteString("}")

		case []interface{}:
			log.Println("[]interface{}")

			if len(vt) > 0 {
				buf.WriteString("[]")
				buf.WriteString(getTypeForValue(vt[0]))
			} else {
				buf.WriteString("[]interface{}")
			}

		default:
			buf.WriteString(getTypeForValue(vt))
		}

		fmt.Fprintf(buf, " `json:\"%s\"`\n", key)
	}

	return nil
}

func jsonToGoName(key string) string {
	// TODO use replacer
	ss := strings.Split(key, "_")
	return strings.Join(MapStringSlice(ss, strings.Title), "")
}

func MapStringSlice(ss []string, f func(string) string) []string {
	for i, s := range ss {
		ss[i] = f(s)
	}
	return ss
}

func getTypeForValue(value interface{}) string {
	//return reflect.TypeOf(value).Name()
	switch vt := value.(type) {
	case string:
		_, err := time.Parse(dateFormat, vt)
		if err == nil {
			return "time.Time"
		}
		return "string"
	case float32, float64:
		return "float32"
	case int, int64:
		return "int"
	case bool:
		return "bool"
	default:
		fmt.Println(reflect.TypeOf(value).Elem(), reflect.TypeOf(value).Name())
		return "Xinterface{}"
	}
}

/*
	if m, isMap := value.(map[string]interface{}); isMap {

	} else if s, isSlice := value.([]interface{}); isSlice {


	} else if reflect.TypeOf(value) == nil {
		return "interface{}"
	}

	return reflect.TypeOf(value).Name()
}
*/

//	fmt.Printf("type %s struct {\n", name)
//	goify(m)
//	fmt.Printf("}\n")
//
//	return nil

func goify(m map[string]interface{}) { //(interface{}, error) {

	for k, v := range m {
		typ := reflect.TypeOf(v).Kind()
		switch typ {
		case reflect.Map:
			fmt.Printf("%s struct {\n", strings.Title(strings.Replace(k, "_", "", 10)))
			goify(v.(map[string]interface{}))
			fmt.Printf("} `json:\"%s\"`\n", k)

		case reflect.Slice:
			fmt.Println("====>", reflect.ValueOf(v))
			fmt.Println("====>", reflect.ValueOf(v), reflect.TypeOf(v).Elem(), reflect.TypeOf(reflect.TypeOf(v).Elem()))
			fmt.Printf("%s []%s `json:\"%s\"`\n", strings.Title(strings.Replace(k, "_", "", 10)), typ, k)
		default:
			fmt.Printf("%s %s `json:\"%s\"`\n", strings.Title(strings.Replace(k, "_", "", 10)), typ, k)
		}
	}
}
