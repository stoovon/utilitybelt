package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/spf13/cobra"
	"github.com/stoovon/utilitybelt/commons"
)

type BasicHello struct {
	logger hclog.Logger
}

func (g *BasicHello) Describe() *commons.CommandDescription {
	description := &commons.CommandDescription{
		Use:   "basic",
		Short: "UtilityTool belt core command",
		Long: `Basic REST API commands.`,
	}

	var listCmd = &commons.CommandDescription{
		Use:   "list",
		Short: "List subcommand description",
	}


	var getCmd = &commons.CommandDescription{
		Use:   "get",
		Short: "Get subcommand description",
		Args: cobra.MinimumNArgs(1),
	}

	description.AddCommand(listCmd)
	description.AddCommand(getCmd)

	return description
}


type Todo struct {
	UserID    int    `json:"userId"`
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

func (g *BasicHello) listTodos() (*string, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos")
	if err != nil {
		return nil, err
	}

	var todos []Todo

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&todos)
	if err != nil {
		return nil, err
	}

	var b strings.Builder

	b.WriteString("ID   | Title\n")
	b.WriteString("=====+==================================\n")

	for i, todo := range todos {
		b.WriteString(fmt.Sprintf("%-4v | %v\n", todo.ID, todo.Title))

		if i > 10 {
			break
		}
	}

	result := b.String()

	return &result, nil
}

func (g *BasicHello) getTodo(args []string) (*string, error) {
	resp, err := http.Get("https://jsonplaceholder.typicode.com/todos/" + args[0])
	if err != nil {
		return nil, err
	}

	var todo Todo

	defer resp.Body.Close()
	err = json.NewDecoder(resp.Body).Decode(&todo)
	if err != nil {
		return nil, err
	}

	var b strings.Builder

	b.WriteString(fmt.Sprintf("Completed: %t", todo.Completed))

	result := b.String()

	return &result, nil
}

func (g *BasicHello) Execute(pathArgs *commons.PathArgs) string {
	switch pathArgs.Path {
	case "basic/list":
		value, err := g.listTodos()
		if err != nil {
			panic(err)
		}
		return *value
	case "basic/get":
		value, err := g.getTodo(pathArgs.Args)
		if err != nil {
			panic(err)
		}
		return *value
	}

	return "Usage \n" +
		" utilitybelt basic list \n" +
		" utilitybelt basic get n"
}

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "UTILITYBELT_PLUGIN",
	MagicCookieValue: "rosebud",
}

func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Output:     os.Stderr,
		JSONFormat: true,
	})

	basic := &BasicHello{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"basic": &commons.UtilityToolPlugin{Impl: basic},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}
