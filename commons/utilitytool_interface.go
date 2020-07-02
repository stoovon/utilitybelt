package commons

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// UtilityTool is the interface that we're exposing as a plugin.
type UtilityTool interface {
	Describe() *CommandDescription
	Execute(pathArgs *PathArgs) string
}

type PathArgs struct {
	Path string
	Args []string
}

// Here is an implementation that talks over RPC
type UtilityToolRPC struct{ client *rpc.Client }

func (g *UtilityToolRPC) Describe() *CommandDescription {
	var resp *CommandDescription
	err := g.client.Call("Plugin.Describe", new(interface{}), &resp)
	if err != nil {
		panic(err)
	}

	return resp
}

func (g *UtilityToolRPC) Execute(pathArgs *PathArgs) string {
	var resp string
	err := g.client.Call("Plugin.Execute", pathArgs, &resp)
	if err != nil {
		// You usually want your interfaces to return errors. If they don't,
		// there isn't much other choice here.
		panic(err)
	}

	return resp
}

// Here is the RPC server that UtilityToolRPC talks to, conforming to
// the requirements of net/rpc
type UtilityToolRPCServer struct {
	// This is the real implementation
	Impl UtilityTool
}

func (s *UtilityToolRPCServer) Describe(args interface{}, resp **CommandDescription) error {
	*resp = s.Impl.Describe()
	return nil
}

func (s *UtilityToolRPCServer) Execute(pathArgs *PathArgs, resp *string) error {
	*resp = s.Impl.Execute(pathArgs)
	return nil
}

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a UtilityToolRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return UtilityToolRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type UtilityToolPlugin struct {
	// Impl Injection
	Impl UtilityTool
}

func (p *UtilityToolPlugin) Server(*plugin.MuxBroker) (interface{}, error) {
	return &UtilityToolRPCServer{Impl: p.Impl}, nil
}

func (UtilityToolPlugin) Client(b *plugin.MuxBroker, c *rpc.Client) (interface{}, error) {
	return &UtilityToolRPC{client: c}, nil
}