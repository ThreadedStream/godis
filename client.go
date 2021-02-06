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

type charptr = []byte
type char = byte

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

func (c *Client) prettifyStrings(args string) []string {
	argList := strings.Split(args, " ")
	//Trim spaces
	for i := range argList {
		argList[i] = strings.Trim(argList[i], "\t\n")
	}

	return argList
}

func (c *Client) prettifyDict(args string) []string {
	argList := strings.Split(args, "{")

	//Trim space in key
	argList[0] = strings.Trim(argList[0], "\t\n")
	argList[0] = strings.TrimSpace(argList[0])
	argList[1] = "{" + argList[1]

	return argList
}

func (c *Client) prettifyList(args string) []string {
	argList := strings.Split(args, "[")

	argList[0] = strings.Trim(argList[0], "\t\n")
	argList[0] = strings.TrimSpace(argList[0])
	argList[1] = "[" + argList[1]

	return argList
}

//key {'key' : 'value'}

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
		//Prettify only in case when arguments are strings
		if strings.IndexAny(args, "{}[]") == -1 {
			argList = c.prettifyStrings(args)
		} else if strings.Index(args, "{") != -1 {
			argList = c.prettifyDict(args)
		} else {
			argList = c.prettifyList(args)
		}

		result := c.commandProcessorPipe(command, argList)
		fmt.Printf("%v\n", result)
	}
}

func (c *Client) runClientPipe(done chan bool) {
	log.Println("Running local client...")
	c.activateVisualPipe()
	done <- true
}
