package main

import (
	"encoding/json"
	"io/ioutil"
	"sync"
)

type HooksChannelMessage struct {
	Name         string
	OriginModule string
}

type CommandsChannelMessage struct {
	Name         string
	OriginModule string
}

type Result struct {
	Hooks    map[string][]string `json:"hooks"`
	Commands map[string]string   `json:"commands"`
}

func main() {
	// TODO: Get this filename dynamically
	filename := "/Users/lshadler/.local/share/@appfabric/plugin-cli/package.json"

	hooks := make(map[string][]string)
	commands := make(map[string]string)

	hooksChannel := make(chan HooksChannelMessage, 100)
	commandsChannel := make(chan CommandsChannelMessage, 100)
	var waitGroup sync.WaitGroup

	oclifManifestPackageJson, err := GetOclifManifest(filename)

	if err != nil {
		panic(err)
	}

	for i := 0; i < len(oclifManifestPackageJson.OclifConfig.Plugins); i++ {
		plugin := oclifManifestPackageJson.OclifConfig.Plugins[i]

		LogPrintln("Plugin name: " + plugin.Name)
		LogPrintln("\tType: " + plugin.Type)
		LogPrintln("\tTag (if module): " + plugin.Tag)
		LogPrintln("\tRoot (if linked): " + plugin.Root)

		waitGroup.Add(1)
		if plugin.Type == "link" {
			go GetPluginInfo(plugin.Root, plugin.Type, hooksChannel, commandsChannel, &waitGroup)
		} else {
			go GetPluginInfo(plugin.Name, plugin.Type, hooksChannel, commandsChannel, &waitGroup)
		}
	}

	waitGroup.Add(1)
	go GetPluginInfo(".", "link", hooksChannel, commandsChannel, &waitGroup)

	go func() {
		for {
			hookMsg := <-hooksChannel
			hookList := hooks[hookMsg.Name]
			hookList = append(hookList, hookMsg.OriginModule)
			LogPrintln("Hook Recieved: " + hookMsg.Name + " | " + hookMsg.OriginModule)
			hooks[hookMsg.Name] = hookList
		}
	}()

	go func() {
		for {
			commandMsg := <-commandsChannel
			LogPrintln("Hook Recieved: " + commandMsg.Name + " | " + commandMsg.OriginModule)
			commands[commandMsg.Name] = commandMsg.OriginModule
		}
	}()

	waitGroup.Wait()

	LogPrintln(commands)
	LogPrintln(hooks)

	file, _ := json.MarshalIndent(Result{
		Hooks:    hooks,
		Commands: commands,
	}, "", " ")

	_ = ioutil.WriteFile("test.json", file, 0644)
}
