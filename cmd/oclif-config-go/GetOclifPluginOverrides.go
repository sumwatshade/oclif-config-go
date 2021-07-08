package main

import (
	"encoding/json"
	"fmt"
)

type OclifManifestPluginDefinition struct {
	Name string `json:"name"`
	Tag  string `json:"tag"`
	Type string `json:"type"`
	Root string `json:"root"`
}

type OclifManifestConfig struct {
	Schema  int                             `json:"schema"`
	Plugins []OclifManifestPluginDefinition `json:"plugins"`
}

type OclifManifestPackageJson struct {
	OclifConfig OclifManifestConfig `json:"oclif"`
}

func GetOclifManifest(filepath string) (OclifManifestPackageJson, error) {
	var oclifManifestPackageJson OclifManifestPackageJson

	manifestByteData, err := ReadFileAsByte(filepath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	} else {
		json.Unmarshal(manifestByteData, &oclifManifestPackageJson)
	}

	return oclifManifestPackageJson, err
}
