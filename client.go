package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	connType = "tcp"
	host     = "0.0.0.0:23"
)

type Client struct {
	parser *Parser
}

func (c *Client) initParserPipe() {
	c.parser = &Parser{}
	c.parser.initParserPipe()
}

func (c *Client) commandProcessorPipe(command string, args []string) interface{} {
	result := c.parser.parseCommand(command, args)
	return result
}

func (c *Client) trimAllSpaces(args []string) []string {
	for i := range args {
		args[i] = strings.Trim(args[i], "\t\n")
	}

	return args
}

func (c *Client) activateVisualPipe() {
	var command string
	var args string
	var argList []string
	var err error
	in := bufio.NewReader(os.Stdin)
	fmt.Println("You may start typing some commands")
	for {
		fmt.Print("godis>")
		command, err = in.ReadString(' ')
		if err != nil {
			log.Println(err)
			return
		}
		args, err = in.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}
		command = strings.TrimSpace(command)
		argList = strings.Split(args, " ")
		argList = c.trimAllSpaces(argList)
		result := c.commandProcessorPipe(command, argList)
		fmt.Printf("%v\n", result)
	}
}

func (c *Client) runClientPipe(done chan bool) {
	log.Println("Running local client...")
	c.activateVisualPipe()
	done <- true
}
