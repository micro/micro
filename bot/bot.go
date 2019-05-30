// Package bot is a Hubot style bot that sits a microservice environment
package bot

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/micro/cli"
	"github.com/micro/go-micro"

	"github.com/micro/go-bot/command"
	"github.com/micro/go-bot/input"
	"github.com/micro/go-log"
	botc "github.com/micro/micro/internal/command/bot"

	proto "github.com/micro/go-bot/proto"

	// inputs
	_ "github.com/micro/go-bot/input/discord"
	_ "github.com/micro/go-bot/input/slack"
	_ "github.com/micro/go-bot/input/telegram"
)

type bot struct {
	exit    chan bool
	ctx     *cli.Context
	service micro.Service

	sync.RWMutex
	inputs   map[string]input.Input
	commands map[string]command.Command
	services map[string]string
}

var (
	// Default server name
	Name = "go.micro.bot"
	// Namespace for commands
	Namespace = "go.micro.bot"
	// map pattern:command
	commands = map[string]func(*cli.Context) command.Command{
		"^echo ":                             botc.Echo,
		"^time$":                             botc.Time,
		"^hello$":                            botc.Hello,
		"^ping$":                             botc.Ping,
		"^list ":                             botc.List,
		"^get ":                              botc.Get,
		"^health ":                           botc.Health,
		"^call ":                             botc.Call,
		"^register ":                         botc.Register,
		"^deregister ":                       botc.Deregister,
		"^(the )?three laws( of robotics)?$": botc.ThreeLaws,
	}
)

func help(commands map[string]command.Command, serviceCommands []string) command.Command {
	usage := "help"
	desc := "Displays help for all known commands"

	var cmds []command.Command

	for _, cmd := range commands {
		cmds = append(cmds, cmd)
	}

	sort.Sort(sortedCommands{cmds})

	return command.NewCommand("help", usage, desc, func(args ...string) ([]byte, error) {
		response := []string{"\n"}
		for _, cmd := range cmds {
			response = append(response, fmt.Sprintf("%s - %s", cmd.Usage(), cmd.Description()))
		}
		response = append(response, serviceCommands...)
		return []byte(strings.Join(response, "\n")), nil
	})
}

func newBot(ctx *cli.Context, inputs map[string]input.Input, commands map[string]command.Command, service micro.Service) *bot {
	commands["^help$"] = help(commands, nil)

	return &bot{
		ctx:      ctx,
		exit:     make(chan bool),
		service:  service,
		commands: commands,
		inputs:   inputs,
		services: make(map[string]string),
	}
}

func (b *bot) loop(io input.Input) {
	log.Logf("[bot][loop] starting %s", io.String())

	for {
		select {
		case <-b.exit:
			log.Logf("[bot][loop] exiting %s", io.String())
			return
		default:
			if err := b.run(io); err != nil {
				log.Logf("[bot][loop] error %v", err)
				time.Sleep(time.Second)
			}
		}
	}
}

func (b *bot) process(c input.Conn, ev input.Event) error {
	args := strings.Split(string(ev.Data), " ")
	if len(args) == 0 {
		return nil
	}

	b.RLock()
	defer b.RUnlock()

	// try built in command
	for pattern, cmd := range b.commands {
		// skip if it doesn't match
		if m, err := regexp.Match(pattern, ev.Data); err != nil || !m {
			continue
		}

		// matched, exec command
		rsp, err := cmd.Exec(args...)
		if err != nil {
			rsp = []byte("error executing cmd: " + err.Error())
		}

		// send response
		return c.Send(&input.Event{
			Meta: ev.Meta,
			From: ev.To,
			To:   ev.From,
			Type: input.TextEvent,
			Data: rsp,
		})
	}

	// no built in match
	// try service commands
	service := Namespace + "." + args[0]

	// is there a service for the command?
	if _, ok := b.services[service]; !ok {
		return nil
	}

	// make service request
	req := b.service.Client().NewRequest(service, "Command.Exec", &proto.ExecRequest{
		Args: args,
	})
	rsp := &proto.ExecResponse{}

	var response []byte

	// call service
	if err := b.service.Client().Call(context.Background(), req, rsp); err != nil {
		response = []byte("error executing cmd: " + err.Error())
	} else if len(rsp.Error) > 0 {
		response = []byte("error executing cmd: " + rsp.Error)
	} else {
		response = rsp.Result
	}

	// send response
	return c.Send(&input.Event{
		Meta: ev.Meta,
		From: ev.To,
		To:   ev.From,
		Type: input.TextEvent,
		Data: response,
	})
}

func (b *bot) run(io input.Input) error {
	log.Logf("[bot][loop] connecting to %s", io.String())

	c, err := io.Stream()
	if err != nil {
		return err
	}

	for {
		select {
		case <-b.exit:
			log.Logf("[bot][loop] closing %s", io.String())
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

			if len(recvEv.Data) == 0 {
				continue
			}

			if err := b.process(c, recvEv); err != nil {
				return err
			}
		}
	}
}

func (b *bot) start() error {
	log.Log("[bot] starting")

	// Start inputs
	for _, io := range b.inputs {
		log.Logf("[bot] starting input %s", io.String())

		if err := io.Init(b.ctx); err != nil {
			return err
		}

		if err := io.Start(); err != nil {
			return err
		}

		go b.loop(io)
	}

	// start watcher
	go b.watch()

	return nil
}

func (b *bot) stop() error {
	log.Log("[bot] stopping")
	close(b.exit)

	// Stop inputs
	for _, io := range b.inputs {
		log.Logf("[bot] stopping input %s", io.String())
		if err := io.Stop(); err != nil {
			log.Logf("[bot] %v", err)
		}
	}

	return nil
}

func (b *bot) watch() {
	commands := map[string]command.Command{}
	services := map[string]string{}

	// copy commands
	b.RLock()
	for k, v := range b.commands {
		commands[k] = v
	}
	b.RUnlock()

	// getHelp retries usage and description from bot service commands
	getHelp := func(service string) (string, error) {
		// is within namespace?
		if !strings.HasPrefix(service, Namespace) {
			return "", fmt.Errorf("%s not within namespace", service)
		}

		if p := strings.TrimPrefix(service, Namespace); len(p) == 0 {
			return "", fmt.Errorf("%s not a service", service)
		}

		// get command help
		req := b.service.Client().NewRequest(service, "Command.Help", &proto.HelpRequest{})
		rsp := &proto.HelpResponse{}

		err := b.service.Client().Call(context.Background(), req, rsp)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("%s - %s", rsp.Usage, rsp.Description), nil
	}

	serviceList, err := b.service.Client().Options().Registry.ListServices()
	if err != nil {
		// log error?
		return
	}

	var serviceCommands []string

	// create service commands
	for _, service := range serviceList {
		h, err := getHelp(service.Name)
		if err != nil {
			continue
		}
		services[service.Name] = h
		serviceCommands = append(serviceCommands, h)
	}

	b.Lock()
	b.commands["^help$"] = help(commands, serviceCommands)
	b.services = services
	b.Unlock()

	w, err := b.service.Client().Options().Registry.Watch()
	if err != nil {
		// log error?
		return
	}

	go func() {
		<-b.exit
		w.Stop()
	}()

	// watch for changes to services
	for {
		res, err := w.Next()
		if err != nil {
			return
		}

		if res.Action == "delete" {
			delete(services, res.Service.Name)
		} else {
			h, err := getHelp(res.Service.Name)
			if err != nil {
				continue
			}
			services[res.Service.Name] = h
		}

		var serviceCommands []string
		for _, v := range services {
			serviceCommands = append(serviceCommands, v)
		}

		b.Lock()
		b.commands["^help$"] = help(commands, serviceCommands)
		b.services = services
		b.Unlock()
	}
}

func run(ctx *cli.Context) {
	// Init plugins
	for _, p := range Plugins() {
		p.Init(ctx)
	}

	if len(ctx.GlobalString("server_name")) > 0 {
		Name = ctx.GlobalString("server_name")
	}

	if len(ctx.String("namespace")) > 0 {
		Namespace = ctx.String("namespace")
	}

	// Parse flags
	if len(ctx.String("inputs")) == 0 {
		log.Fatal("[bot] no inputs specified")
	}

	inputs := strings.Split(ctx.String("inputs"), ",")
	if len(inputs) == 0 {
		log.Fatal("[bot] no inputs specified")
	}

	ios := make(map[string]input.Input)
	cmds := make(map[string]command.Command)

	// create built in commands
	for pattern, cmd := range commands {
		cmds[pattern] = cmd(ctx)
	}

	// take other commands
	for pattern, cmd := range command.Commands {
		if c, ok := cmds[pattern]; ok {
			log.Logf("[bot] command %s already registered for pattern %s\n", c.String(), pattern)
			continue
		}
		// register command
		cmds[pattern] = cmd
	}

	// Parse inputs
	for _, io := range inputs {
		i, ok := input.Inputs[io]
		if !ok {
			log.Logf("[bot] input %s not found\n", i)
			os.Exit(1)
		}
		ios[io] = i
	}

	// setup service
	service := micro.NewService(
		micro.Name(Name),
		micro.RegisterTTL(
			time.Duration(ctx.GlobalInt("register_ttl"))*time.Second,
		),
		micro.RegisterInterval(
			time.Duration(ctx.GlobalInt("register_interval"))*time.Second,
		),
	)

	// Start bot
	b := newBot(ctx, ios, cmds, service)

	if err := b.start(); err != nil {
		log.Logf("error starting bot %v", err)
		os.Exit(1)
	}

	// Run server
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}

	// Stop bot
	if err := b.stop(); err != nil {
		log.Logf("error stopping bot %v", err)
	}
}

func Commands() []cli.Command {
	flags := []cli.Flag{
		cli.StringFlag{
			Name:   "inputs",
			Usage:  "Inputs to load on startup",
			EnvVar: "MICRO_BOT_INPUTS",
		},
		cli.StringFlag{
			Name:   "namespace",
			Usage:  "Set the namespace used by the bot to find commands e.g. com.example.bot",
			EnvVar: "MICRO_BOT_NAMESPACE",
		},
	}

	// setup input flags
	for _, input := range input.Inputs {
		flags = append(flags, input.Flags()...)
	}

	command := cli.Command{
		Name:   "bot",
		Usage:  "Run the chatops bot",
		Flags:  flags,
		Action: run,
	}

	for _, p := range Plugins() {
		if cmds := p.Commands(); len(cmds) > 0 {
			command.Subcommands = append(command.Subcommands, cmds...)
		}

		if flags := p.Flags(); len(flags) > 0 {
			command.Flags = append(command.Flags, flags...)
		}
	}

	return []cli.Command{command}
}
