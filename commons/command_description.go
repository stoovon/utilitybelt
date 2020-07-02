package commons

import (
	"github.com/spf13/cobra"
)

// Essentially a Cobra commander command, but exporting the sub-commands so they can also traverse RPC.
type CommandDescription struct {
	Use string
	Short string
	Long string
	SubCommands []*CommandDescription
	Run func(cmd *cobra.Command, args []string)
	Args cobra.PositionalArgs
}

func (c *CommandDescription) AddCommand(command *CommandDescription) {
	if c.SubCommands == nil {
		c.SubCommands = []*CommandDescription{}
	}

	c.SubCommands = append(c.SubCommands, command)
}

func (c *CommandDescription) ToCobraCommand(path string, run func(path string) func(cmd *cobra.Command, args []string)) *cobra.Command {
	result := &cobra.Command {
		Use:   c.Use,
		Short: c.Short,
		Long: c.Long,
		Args: c.Args,
		Run: run(path),
	}

	if c.SubCommands != nil && len(c.SubCommands) > 0 {
		for _, subCommand := range c.SubCommands {
			result.AddCommand(subCommand.ToCobraCommand(path + "/" + subCommand.Use, run))
		}
	}

	return result
}