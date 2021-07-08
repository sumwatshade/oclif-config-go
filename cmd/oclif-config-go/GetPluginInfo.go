package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"strings"
	"sync"
)

type OclifPluginConfig struct {
	DevPlugins []string               `json:"devPlugins"`
	Plugins    []string               `json:"plugins"`
	Hooks      map[string]interface{} `json:"hooks"`
	Commands   string                 `json:"commands"`
}

type OclifCliPackageJson struct {
	Config  OclifPluginConfig `json:"oclif"`
	Name    string            `json:"name"`
	Version string            `json:"version"`
}

func GetPackageJson(moduleRoot string) (OclifCliPackageJson, error) {
	var oclifCliPackageJson OclifCliPackageJson

	packageJsonByteData, err := ReadFileAsByte(moduleRoot + "/package.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		panic(err)
	} else {
		json.Unmarshal(packageJsonByteData, &oclifCliPackageJson)
	}

	return oclifCliPackageJson, err
}

func GetPluginInfo(PluginName string, Type string, hooksChannel chan HooksChannelMessage, commandsChannel chan CommandsChannelMessage, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	moduleRoot := "../../node_modules/" + PluginName
	if Type == "link" {
		moduleRoot = PluginName
	}

	pluginInfo, err := GetPackageJson(moduleRoot)

	if err != nil {
		panic(err)
	}

	implementedHooks := GetAllKeys(pluginInfo.Config.Hooks)

	var exposedCommands []string

	if len(pluginInfo.Config.Commands) > 0 {
		files, err := ioutil.ReadDir(moduleRoot + "/" + pluginInfo.Config.Commands)
		if err != nil {
			log.Fatal(err)
		}
		for _, f := range files {
			if strings.Contains(f.Name(), ".js") {
				commandName := strings.ReplaceAll(f.Name(), ".js", "")
				exposedCommands = append(exposedCommands, commandName)
			}
		}
	}

	LogPrintln("Plugin name: " + pluginInfo.Name)
	LogPrintln("\tHooks: " + strings.Join(implementedHooks, ", "))
	LogPrintln("\tCommands Exposed: " + strings.Join(exposedCommands, ", "))

	for i := range exposedCommands {
		commandsChannel <- CommandsChannelMessage{
			OriginModule: pluginInfo.Name,
			Name:         exposedCommands[i],
		}
	}

	for i := range implementedHooks {
		hooksChannel <- HooksChannelMessage{
			OriginModule: pluginInfo.Name,
			Name:         implementedHooks[i],
		}
	}

	numAdditionalPlugins := len(pluginInfo.Config.Plugins)
	waitGroup.Add(numAdditionalPlugins)
	for i := 0; i < len(pluginInfo.Config.Plugins); i++ {
		go GetPluginInfo(pluginInfo.Config.Plugins[i], "module", hooksChannel, commandsChannel, waitGroup)
	}
}
