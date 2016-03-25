package bot

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/micro/cli"
	"github.com/micro/micro/bot/command"
	"github.com/micro/micro/bot/input"
	_ "github.com/micro/micro/bot/input/slack"
)

type bot struct {
	inputs   []input.Input
	commands []command.Command
	exit     chan bool
}

var (
	commands = map[string]func(*cli.Context) command.Command{
		"hello":      command.Hello,
		"ping":       command.Ping,
		"list":       command.List,
		"get":        command.Get,
		"health":     command.Health,
		"query":      command.Query,
		"register":   command.Register,
		"deregister": command.Deregister,
	}
)

func help(commands []command.Command) command.Command {
	usage := "help"
	desc := "Displays help for all known commands"

	return command.NewCommand("help", usage, desc, func(args ...string) ([]byte, error) {
		var response []string
		for _, cmd := range commands {
			response = append(response, fmt.Sprintf("%s - %s", cmd.Usage(), cmd.Description()))
		}
		return []byte(strings.Join(response, "\n")), nil
	})
}

func newBot(inputs []input.Input, commands []command.Command) *bot {
	// generate help command
	commands = append(commands, help(commands))

	return &bot{
		inputs:   inputs,
		commands: commands,
		exit:     make(chan bool),
	}
}

func (b *bot) loop(io input.Input) {
	fmt.Println("[bot] starting", io.String(), "loop")

	for {
		select {
		case <-b.exit:
			return
		default:
			if err := b.run(io); err != nil {
				fmt.Println(err)
				time.Sleep(time.Second)
			}
		}
	}
}

func (b *bot) run(io input.Input) error {
	fmt.Println("[bot] connecting to", io.String())

	c, err := io.Connect()
	if err != nil {
		return err
	}

	for {
		select {
		case <-b.exit:
			return c.Close()
		default:
			var ev input.Event
			if err := c.Recv(&ev); err != nil {
				return err
			}
			// TODO: do something with this
			fmt.Println("received %+v", ev)
		}
	}
}

func (b *bot) start() {
	fmt.Println("[bot] starting")
	/*
		for _, io := range b.inputs {
			go b.loop(io)
		}
	*/

	// register commands
	for _, io := range b.inputs {
		for _, cmd := range b.commands {
			fmt.Printf("[bot] registering %s command with %s\n", cmd.Name(), io.String())
			if err := io.Process(cmd); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}
	}
}

func (b *bot) stop() {
	close(b.exit)
}

func run(ctx *cli.Context) {
	// Parse flags
	if len(ctx.String("inputs")) == 0 {
		fmt.Println("[bot] no inputs specified")
		os.Exit(1)
	}

	inputs := strings.Split(ctx.String("inputs"), ",")
	if len(inputs) == 0 {
		fmt.Println("[bot] no inputs specified")
		os.Exit(1)
	}

	var ios []input.Input
	var cmds []command.Command

	// create commands
	for _, cmd := range commands {
		cmds = append(cmds, cmd(ctx))
	}

	// Parse inputs
	for _, io := range inputs {
		i, ok := input.Inputs[io]
		if !ok {
			fmt.Printf("[bot] input %s not found\n", i)
			os.Exit(1)
		}
		ios = append(ios, i)
	}

	// Start inputs
	for _, io := range ios {
		fmt.Println("[bot] starting input", io.String())

		if err := io.Init(ctx); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := io.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Start bot
	b := newBot(ios, cmds)
	b.start()

	// Exit on kill signal
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	fmt.Println(<-ch)

	// Stop bot
	b.stop()

	// Stop inputs
	for _, io := range ios {
		fmt.Println("[bot] stopping input", io.String())
		if err := io.Stop(); err != nil {
			fmt.Println(err)
		}
	}
}

func Commands() []cli.Command {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:  "inputs",
			Usage: "Inputs to load on startup",
		},
	}

	// setup input flags
	for _, input := range input.Inputs {
		flags = append(flags, input.Flags()...)
	}

	return []cli.Command{
		{
			Name:   "bot",
			Usage:  "Run the micro bot",
			Flags:  flags,
			Action: run,
		},
	}
}
