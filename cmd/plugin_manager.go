package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"github.com/stoovon/utilitybelt/commons"
)

type PluginManager struct {
	// pluginMap is the map of plugins we can dispense.
	pluginMap map[string]plugin.Plugin

	// Separate execs so we don't reuse
	describeMap map[string]*exec.Cmd
	executeMap  map[string]*exec.Cmd
}

func (p *PluginManager) init() []*cobra.Command {
	availablePlugins, err := plugin.Discover("*", "./compiled/plugins")
	if err != nil {
		log.Fatal(err)
	}

	p.pluginMap = map[string]plugin.Plugin{}
	p.describeMap = map[string]*exec.Cmd{}
	p.executeMap = map[string]*exec.Cmd{}

	var commands []*cobra.Command

	for _, availablePlugin := range availablePlugins {

		pluginName := filepath.Base(availablePlugin)
		p.pluginMap[pluginName] = &commons.UtilityToolPlugin{}
		p.describeMap[pluginName] = exec.Command(availablePlugin)
		p.executeMap[pluginName] = exec.Command(availablePlugin)

		description, err := p.describe(pluginName)
		if err != nil {
			panic(err)
		}

		newCommand := description.ToCobraCommand(description.Use, func(path string) func(cmd *cobra.Command, args []string) {
			return func(cmd *cobra.Command, args []string) {
				retval, err := p.execute(pluginName, path, args)
				if err != nil {
					log.Fatal(err)
				}

				if len(*retval) > 0 {
					fmt.Println(*retval)
				}
			}
		})

		commands = append(commands, newCommand)
	}

	return commands
}

func (p *PluginManager) describe(name string) (result *commons.CommandDescription, err error) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         p.pluginMap,
		Cmd:             p.describeMap[name],
		Logger:          logger,
	})
	defer client.Kill()

	rpcClient, err := client.Client()

	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(name)
	if err != nil {
		return nil, err
	}

	// Execute the command over RPC, passing arguments directly or marshalled.
	tool := raw.(commons.UtilityTool)

	res := tool.Describe()

	return res, nil
}


func (p *PluginManager) execute(name string, path string, args []string) (result *string, err error) {
	logger := hclog.New(&hclog.LoggerOptions{
		Name:   "plugin",
		Output: os.Stdout,
		Level:  hclog.Info,
	})

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         p.pluginMap,
		Cmd:             p.executeMap[name],
		Logger:          logger,
	})
	defer client.Kill()

	rpcClient, err := client.Client()

	if err != nil {
		return nil, err
	}

	raw, err := rpcClient.Dispense(name)
	if err != nil {
		return nil, err
	}

	// Execute the command over RPC, passing arguments directly or marshalled.
	tool := raw.(commons.UtilityTool)

	res := tool.Execute(&commons.PathArgs{
		Path: path,
		Args: args,
	})

	return &res, nil
}

// Ensure compatibility for UX
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "UTILITYBELT_PLUGIN",
	MagicCookieValue: "rosebud",
}
