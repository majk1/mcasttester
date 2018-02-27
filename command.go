package main

import (
	"flag"
	"fmt"
)

type Command struct {
	name        string
	description string
	flags       *flag.FlagSet
	addr        string
	payload     string
	loopCount   int
	execFunc    func(addr string, payload string, loopCount int)
}

type Commands struct {
	items [] *Command
}

func (c *Commands) AddCommand(name string, description string) *Command {
	command := &Command{
		name:        name,
		description: description,
	}
	command.flags = flag.NewFlagSet(command.name, flag.ExitOnError)
	c.items = append(c.items, command)
	return command
}

func (c *Commands) GetByName(commandName string) *Command {
	for _, cmd := range c.items {
		if cmd.name == commandName {
			return cmd
		}
	}
	return nil
}

func (c *Commands) PrintUsage(appName string) {
	fmt.Printf("Usage: %s <command> [options]\n", appName)
	for _, cmd := range c.items {
		fmt.Printf("  %-10s - %s\n", cmd.name, cmd.description)
	}
}

func (c *Command) ParseArgs(args []string) {
	c.flags.Parse(args)
}

func (c *Command) Execute() {
	c.execFunc(c.addr, c.payload, c.loopCount)
}
