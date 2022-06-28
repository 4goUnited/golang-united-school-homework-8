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
	Id string `json:"id"`
	Email string `json:"email"`
	Age int `json:"age"`
}


const filePerm = 0644

//Open file(or create if not exists), read all that in there and output to writer
func List(fileName string, writer io.Writer) error {
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

//Add item (if not exists). Main idea was not to write all readed slice to file, but only one record.
func Add(args Arguments, userInput User, writer io.Writer) error {
	if args["item"] == "" {
		return errors.New("-item flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR|os.O_CREATE, filePerm)
	if err != nil {
		return err
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return err
	}

	var users []User

	//If file just created, need to adjust prfx to the start of the file
	prfx := "["
	if stat.Size() > 0 {

	        bytes, err := io.ReadAll(file)
	        if err != nil {
			return err
	        }

	        if err := json.Unmarshal(bytes, &users); err != nil {
			fmt.Println("here")
			return err
	        }

		//Search if user is already exists
	        for _, u := range users {
			if u.Id == userInput.Id {
				if _, err := writer.Write([]byte(fmt.Sprintf("Item with id %v already exists", u.Id))); err != nil {
					return err
				}
				return nil
			}
	        }

		//If no user exists, added it to the end of file
		prfx = ",\n"
		if _, err := file.Seek(-1, 2); err != nil {
			return err
		}

	}

	if _, err := file.WriteString(prfx); err != nil {
		return err
	}

	enc := json.NewEncoder(file)
	if err := enc.Encode(userInput); err != nil {
		return err
	}

	if _, err := file.Seek(-1, 2); err != nil {
		return err
	}

	if _, err := file.WriteString("]"); err != nil {
		return err
	}

	return nil
}

//Simple linear search for the User by ID
func FindById(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}

	file, err := os.Open(args["fileName"])
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

	//Search for user and output it to writer 
	for _, u := range users {
		if u.Id == args["id"] {
			jsn, err := json.Marshal(u)
			if err != nil {
				return err
			}

			if _, err := writer.Write([]byte(jsn)); err != nil {
				return err
			}
			return nil
		}
	}

	//If we don'f find User with needed ID, just output "" to writer 
	file.Write([]byte(""))
	return nil
}

//If finded - remove user from byte slice and then write to truncated (to 0) file
func Remove(args Arguments, writer io.Writer) error {
	if args["id"] == "" {
		return errors.New("-id flag has to be specified")
	}

	file, err := os.OpenFile(args["fileName"], os.O_RDWR, filePerm)
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

	for i, u := range users {
		if u.Id == args["id"] {
			users = append(users[:i], users[i+1:]...)

			jsn, err := json.Marshal(users)
			if err != nil {
				return err
			}

			file.Truncate(0)
			file.Seek(0, 0)
			_, err = file.Write([]byte(jsn))
			if err != nil {
				return err
			}
			return nil
		}
	}

	if _, err := writer.Write([]byte(fmt.Sprintf("Item with id %v not found", args["id"]))); err != nil {
		return err
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

	if args["fileName"] == "" {
		return errors.New("-fileName flag has to be specified")
	}

	switch args["operation"] {
		case "":
			return errors.New("-operation flag has to be specified")
		case "list":
			err := List(args["fileName"], writer)
			return err
		case "add":
			err := Add(args, userInput, writer)
			return err
		case "findById":
			err := FindById(args, writer)
			return err
		case "remove":
			err := Remove(args, writer)
			return err
		default:
			return fmt.Errorf("Operation %s not allowed!", args["operation"])
		}
	return nil
}

func parseArgs() Arguments {
	var oFlag = flag.String("operation", "", "Choose \"add\",\"list\",\"findById\" or \"remove\" operation.")
	var idFlag = flag.String("id", "", "Enter ID \"id\" 1")
	var iFlag = flag.String("item", "", "Enter user `{\"id\": \"1\", \"email\":\"email@test.com\",\"age\": 23}`")
	var fFlag = flag.String("fileName", "", "Enter file \"users.json\"")

	flag.Parse()

	args := Arguments{
		"id": *idFlag,
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
