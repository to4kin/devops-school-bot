package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var sep string = "|"

// Callback represents a telebot.Callback().Data() for inline messages
type Callback struct {
	// ID of element in database
	ID int64
	// Type of element
	Type string
	// Command for button
	Command string
	// Command for elements in list
	ListCommand string
}

// ToString converts Callback object to string with separator
//
// NOTE: default separator is "|", so the string will be ->
// ID|Type|Command|ListCommand
func (c *Callback) ToString() string {
	return fmt.Sprintf("%d%v%v%v%v%v%v", c.ID, sep, c.Type, sep, c.Command, sep, c.ListCommand)
}

// Unmarshal converts string to Callback object
//
// NOTE: string should be in a correct format ->
// ID|Type|Command|ListCommand
func (c *Callback) Unmarshal(s string) error {
	str := strings.Split(s, sep)

	if len(str) != 4 {
		return errors.New("incorrect string")
	}

	var err error
	c.ID, err = strconv.ParseInt(str[0], 10, 64)
	if err != nil {
		return err
	}

	c.Type = str[1]
	c.Command = str[2]
	c.ListCommand = str[3]

	return nil
}
