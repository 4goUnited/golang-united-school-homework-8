package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
)

type Arguments map[string]string

type User struct {
	Id string
	Email string
	Age int
}


const filePerm = 0644

func List(fileName string, writer io.Writer) error {
	if fileName == "" {
		return errors.New("-fileName flag has to be specified")
	}

	file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	if _, err := writer.Write(bytes); err != nil {
		return err
	}

	return nil
}

func Add(args Arguments, inputId string, writer io.Writer) error {
	if args["item"] == "" {
		return errors.New("-item flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var users []User
	if err := json.Unmarshal(bytes, &users); err != nil {
		return err
	}

	for _, u := range users {
		if u.Id == inputId {
			if _, err := writer.Write([]byte(fmt.Sprintf("Item with id %v already exists", inputId))); err != nil {
				return err
			}
			return nil
		}
	}
	return nil
}

func Perform(args Arguments, writer io.Writer) error {
	// Variable for input item User
	var userInput User
	if args["item"] != "" {
		if err := json.Unmarshal([]byte(args["item"]), &userInput); err != nil {
			return err
		}
	}

	switch args["operation"] {
		case "":
			return errors.New("-operation flag has to be specified")
		case "list":
			err := List(args["fileName"], writer)
			return err
		case "add":
			err := Add(args, userInput.Id, writer)
			return err
		case "findById":
			//
		case "remove":
			//
		default:
			return fmt.Errorf("Operation %s not allowed!", args["operation"])
		}
	return nil
}

func parseArgs() Arguments {
	var oFlag = flag.String("operation", "", "Choose \"add\",\"list\",\"findById\" or \"remove\" operation.")
	var iFlag = flag.String("item", "", "Enter user `{\"id\": \"1\", \"email\":\"email@test.com\",\"age\": 23}`")
	var fFlag = flag.String("fileName", "", "Enter file \"users.json\"")

	flag.Parse()

	args := Arguments{
		"id": "",
		"operation": *oFlag,
		"item": *iFlag,
	        "fileName": *fFlag,
	}
	return args
}

func main() {
	err := Perform(parseArgs(), os.Stdout)
	if err != nil {
		panic(err)
	}
}
