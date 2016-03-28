package bot

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/micro/bot/command"
	"github.com/micro/micro/bot/input"
	_ "github.com/micro/micro/bot/input/hipchat"
	_ "github.com/micro/micro/bot/input/slack"
)

type bot struct {
	inputs   map[string]input.Input
	commands map[string]command.Command
	exit     chan bool
	ctx      *cli.Context
}

var (
	// map pattern:command
	commands = map[string]func(*cli.Context) command.Command{
		"^echo ":                             command.Echo,
		"^time$":                             command.Time,
		"^hello$":                            command.Hello,
		"^ping$":                             command.Ping,
		"^list ":                             command.List,
		"^get ":                              command.Get,
		"^health ":                           command.Health,
		"^query ":                            command.Query,
		"^register ":                         command.Register,
		"^deregister ":                       command.Deregister,
		"^(the )?three laws( of robotics)?$": command.ThreeLaws,
	}
)

func help(commands map[string]command.Command) command.Command {
	usage := "help"
	desc := "Displays help for all known commands"

	var cmds []command.Command

	for _, cmd := range commands {
		cmds = append(cmds, cmd)
	}

	sort.Sort(sortedCommands{cmds})

	return command.NewCommand("help", usage, desc, func(args ...string) ([]byte, error) {
		response := []string{"\n"}
		for _, cmd := range commands {
			response = append(response, fmt.Sprintf("%s - %s", cmd.Usage(), cmd.Description()))
		}
		return []byte(strings.Join(response, "\n")), nil
	})
}

func newBot(ctx *cli.Context, inputs map[string]input.Input, commands map[string]command.Command) *bot {
	// generate help command
	commands["^help$"] = help(commands)

	return &bot{
		inputs:   inputs,
		commands: commands,
		exit:     make(chan bool),
		ctx:      ctx,
	}
}

func (b *bot) loop(io input.Input) {
	fmt.Println("[bot][loop] starting", io.String())

	for {
		select {
		case <-b.exit:
			fmt.Println("[bot][loop] exiting", io.String())
			return
		default:
			if err := b.run(io); err != nil {
				fmt.Println("[bot][loop] error", err)
				time.Sleep(time.Second)
			}
		}
	}
}

func (b *bot) run(io input.Input) error {
	fmt.Println("[bot][loop] connecting to", io.String())

	c, err := io.Stream()
	if err != nil {
		return err
	}

	for {
		select {
		case <-b.exit:
			fmt.Println("[bot][loop] closing", io.String())
			return c.Close()
		default:
			var recvEv input.Event
			// receive input
			if err := c.Recv(&recvEv); err != nil {
				return err
			}

			// only process TextEvent
			if recvEv.Type != input.TextEvent {
				continue
			}

			// process command
			for pattern, cmd := range b.commands {
				// skip if it doesn't match
				if m, err := regexp.Match(pattern, recvEv.Data); err != nil || !m {
					continue
				}

				// matched, exec command
				args := strings.Split(string(recvEv.Data), " ")
				rsp, err := cmd.Exec(args...)
				if err != nil {
					rsp = []byte("error executing cmd: " + err.Error())
				}

				// send response
				if err := c.Send(&input.Event{
					Meta: recvEv.Meta,
					From: recvEv.To,
					To:   recvEv.From,
					Type: input.TextEvent,
					Data: rsp,
				}); err != nil {
					return err
				}

				// done
				break
			}
		}
	}
}

func (b *bot) start() {
	fmt.Println("[bot] starting")

	// Start inputs
	for _, io := range b.inputs {
		fmt.Println("[bot] starting input", io.String())

		if err := io.Init(b.ctx); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		if err := io.Start(); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		go b.loop(io)
	}
}

func (b *bot) stop() {
	fmt.Println("[bot] stopping")
	close(b.exit)

	// Stop inputs
	for _, io := range b.inputs {
		fmt.Println("[bot] stopping input", io.String())
		if err := io.Stop(); err != nil {
			fmt.Println(err)
		}
	}
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

	ios := make(map[string]input.Input)
	cmds := make(map[string]command.Command)

	// create commands
	for pattern, cmd := range commands {
		cmds[pattern] = cmd(ctx)
	}

	// Parse inputs
	for _, io := range inputs {
		i, ok := input.Inputs[io]
		if !ok {
			fmt.Printf("[bot] input %s not found\n", i)
			os.Exit(1)
		}
		ios[io] = i
	}

	// Start bot
	b := newBot(ctx, ios, cmds)
	b.start()

	service := micro.NewService(
		micro.Name("go.micro.bot"),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop bot
	b.stop()

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
