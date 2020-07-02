package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"
	"github.com/mum4k/termdash"
	"github.com/mum4k/termdash/align"
	"github.com/mum4k/termdash/cell"
	"github.com/mum4k/termdash/container"
	"github.com/mum4k/termdash/linestyle"
	"github.com/mum4k/termdash/terminal/termbox"
	"github.com/mum4k/termdash/terminal/terminalapi"
	"github.com/mum4k/termdash/widgets/button"
	"github.com/mum4k/termdash/widgets/segmentdisplay"
	"github.com/stoovon/utilitybelt/commons"
)

type ButtonsHello struct {
	logger hclog.Logger
}

func (g *ButtonsHello) Describe() *commons.CommandDescription {
	return &commons.CommandDescription{
		Use:   "buttons",
		Short: "Add/Subtract manager",
		Long: `Demonstration of a simple plugin-based TUI binding to CobraCommander.`,
	}
}

func (g *ButtonsHello) Execute(args *commons.PathArgs) string {
	t, err := termbox.New()
	if err != nil {
		panic(err)
	}
	defer t.Close()

	ctx, cancel := context.WithCancel(context.Background())

	val := 0
	display, err := segmentdisplay.New()
	if err != nil {
		panic(err)
	}
	if err := display.Write([]*segmentdisplay.TextChunk{
		segmentdisplay.NewChunk(fmt.Sprintf("%d", val)),
	}); err != nil {
		panic(err)
	}

	addB, err := button.New("(a)dd", func() error {
		val++
		return display.Write([]*segmentdisplay.TextChunk{
			segmentdisplay.NewChunk(fmt.Sprintf("%d", val)),
		})
	},
		button.GlobalKey('a'),
		button.WidthFor("(s)ubtract"),
	)
	if err != nil {
		panic(err)
	}

	subB, err := button.New("(s)ubtract", func() error {
		val--
		return display.Write([]*segmentdisplay.TextChunk{
			segmentdisplay.NewChunk(fmt.Sprintf("%d", val)),
		})
	},
		button.FillColor(cell.ColorNumber(220)),
		button.GlobalKey('s'),
	)
	if err != nil {
		panic(err)
	}

	c, err := container.New(
		t,
		container.Border(linestyle.Light),
		container.BorderTitle("PRESS Q TO QUIT"),
		container.SplitHorizontal(
			container.Top(
				container.PlaceWidget(display),
			),
			container.Bottom(
				container.SplitVertical(
					container.Left(
						container.PlaceWidget(addB),
						container.AlignHorizontal(align.HorizontalRight),
					),
					container.Right(
						container.PlaceWidget(subB),
						container.AlignHorizontal(align.HorizontalLeft),
					),
				),
			),
			container.SplitPercent(60),
		),
	)
	if err != nil {
		panic(err)
	}

	quitter := func(k *terminalapi.Keyboard) {
		if k.Key == 'q' || k.Key == 'Q' {
			cancel()
		}
	}

	if err := termdash.Run(ctx, t, c, termdash.KeyboardSubscriber(quitter), termdash.RedrawInterval(100*time.Millisecond)); err != nil {
		panic(err)
	}

	return ""
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

	buttons := &ButtonsHello{
		logger: logger,
	}
	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"buttons": &commons.UtilityToolPlugin{Impl: buttons},
	}

	logger.Debug("message from plugin", "foo", "bar")

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
	})
}