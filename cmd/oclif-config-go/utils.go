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

	// log.Println("Successfully Opened " + filepath)
	byteValue, _ := ioutil.ReadAll(file)

	return byteValue, err
}

func GetAllKeys(inputMap map[string]interface{}) []string {
	keys := make([]string, 0, len(inputMap))
	for k := range inputMap {
		keys = append(keys, k)
	}

	return keys
}
