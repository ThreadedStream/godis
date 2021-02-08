package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type UserCache struct {
	username string
	//Some additional information may be added
}

type Client struct {
	parser Parser
	user   UserCache
}

func (c *Client) initParserPipe() {
	c.parser = Parser{}
	c.user.username = "anonymous"
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

func (c *Client) getCommandAndArgs(str string) (string, string) {
	currIndex := 0
	var command string

	//Extracting command
	for str[currIndex] != ' ' && str[currIndex] != '\n' {
		command += fmt.Sprintf("%c", str[currIndex])
		currIndex++
	}
	if str[currIndex] != '\n' {
		currIndex++
		return command, fmt.Sprintf("%s", str[currIndex:len(str)])
	}
	return command, ""

}

func (c *Client) activateVisualPipe() {
	var command string
	var args string
	var argList []string
	var err error
	in := bufio.NewReader(os.Stdin)
	fmt.Println("You may start typing some commands. Press q to quit")
	for {
		fmt.Print("godis>")
		args, err = in.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				fmt.Println(err.Error())
			}
		}

		command, args = c.getCommandAndArgs(args)
		command = strings.TrimSpace(command)
		//Prettify only in case when arguments are strings
		if strings.IndexAny(args, "{}[]") == -1 && args != "" {
			argList = c.prettifyStrings(args)
		} else if strings.Index(args, "{") != -1 {
			argList = c.prettifyDict(args)
		} else if strings.Index(args, "[") != -1 {
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
