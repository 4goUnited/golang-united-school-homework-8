package main

import (
	"encoding/json"
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

// Variable for User to Unmarshal only once at parseArgs()
var u User

func Perform(args Arguments, writer io.Writer) error {
	if args.operation == nil {
		return errors.New("-operation flag has to be specified")
	fmt.Println(u.Email)
	return nil
}

func parseArgs() Arguments {
	var oFlag = flag.String("operation", "", "Choose \"add\",\"list\",\"findById\" or \"remove\" operation.")
	var iFlag = flag.String("item", "", "Enter user `{\"id\": \"1\", \"email\":\"email@test.com\",\"age\": 23}`")
	var fFlag = flag.String("fileName", "", "Enter file \"users.json\"")

	flag.Parse()

	if err := json.Unmarshal([]byte(*iFlag), &u); err != nil {
		panic(err)
	}

	fmt.Printf("User.ID = %v, User.Email = %v, User.Age = %v\n", u.Id, u.Email, u.Age)

	args := Arguments{
		"id": u.Id,
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
