package main

import (
	"io/ioutil"
	"log"
	"os"
)

// Reads in a file based on path and converts it to
// 	an array of byte data that can be parsed by the
//  JSON unmarshalling library
func ReadFileAsByte(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	// if we os.Open returns an error then handle it
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	LogPrintln("Successfully Opened " + filepath)
	byteValue, _ := ioutil.ReadAll(file)

	return byteValue, err
}

// Gets all keys from a map, returned as a string array.
func GetAllKeys(inputMap map[string]interface{}) []string {
	keys := make([]string, 0, len(inputMap))
	for k := range inputMap {
		keys = append(keys, k)
	}

	return keys
}

func LogPrintln(msg ...interface{}) {
	debugEnv := os.Getenv("DEBUG")
	if len(debugEnv) > 0 {
		log.Println(msg...)
	}
}
