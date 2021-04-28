package argsParse

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

type DataType int

const (
	BOOLEAN DataType = iota
	STRING
	INTEGER
	FLOAT
	UINT
	ERROR
)

type Argument struct {
	Name      string
	Triggers  []string
	Defintion string
	Value     interface{}
	Type      DataType
}

var Arguments = []Argument{}
var DefintionsLoaded = false
var Parsed = false

func (a *Argument) TypeString() string {
	switch a.Type {
	case 0:
		return "BOOLEAN"
	case 1:
		return "STRING"
	case 2:
		return "INTEGER"
	case 3:
		return "FLOAT"
	case 4:
		return "UINT"
	case 5:
		return "ERROR"
	default:
		return "ERROR"
	}
}

func AddDefinition(name string, triggers []string, definition string,
	value interface{}, typ int) {
	arg := Argument{
		Name:      name,
		Triggers:  triggers,
		Defintion: definition,
		Value:     value,
		Type:      DataType(typ),
	}
	Arguments = append(Arguments, arg)
	DefintionsLoaded = true
}

//This function will take the path to the string definitions and load them into
//the global variable "Arguments" (array of argument)
func LoadDefinitions(path string) {
	argDefs := readFile(path) //get byte arr of the definitions file

	if err := json.Unmarshal(argDefs, &Arguments); err != nil {
		fmt.Println(err)
		log.Fatal(err)
	} else {
		DefintionsLoaded = true
	}
}

func readFile(path string) []byte {
	content, err := ioutil.ReadFile(path)

	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}
	return content
}

//This will output to the terminal all of the definitions
func helpTrigger() {
	if DefintionsLoaded {
		for _, d := range Arguments {
			var triggers string

			for _, s := range d.Triggers {
				if len(s) > 1 {
					triggers += fmt.Sprintf("[--%v]", s)
				} else {
					triggers += fmt.Sprintf("[-%v]", s)
				}
			}

			msg := fmt.Sprintf("Command: %v, Default Value: %v, Usage: %v\r\n",
				triggers, d.Value, d.Defintion)

			fmt.Println(msg)
		}
	}
}

func HasArg(name string) bool {
	for _, a := range Arguments {
		if a.Name == name {
			return true
		}
	}
	return false
}

func GetArg(name string) (Argument, bool) {
	for i := 0; i < len(Arguments); i++ {
		if Arguments[i].Name == name {
			return Arguments[i], true
		}
	}
	return Argument{}, false
}

func SetArgValue(name string, value interface{}) bool {
	for i := 0; i < len(Arguments); i++ {
		if Arguments[i].Name == name {
			Arguments[i].Value = value
			return true
		}
	}
	return false
}

func Parse() {
	// definitions error
	if !DefintionsLoaded {
		fmt.Println("Error: Definitions need to be added or loaded before parsing!")
		os.Exit(1)
	}

	//local var
	flags := parseFlags()
	flagCount := len(flags)

	for i := 0; i < flagCount; i++ {
		flag := flags[i]
		if flag == "--help" || flag == "-h" {
			helpTrigger()
			os.Exit(0)
		}

		isKey := hasKey(flag)
		hasNext := i+1 < flagCount
		var next string

		if hasNext {
			next = flags[i+1]
		}

		// //this assumes that any flag that does not have a value is setting a boolean true
		if isKey && hasNext && !hasKey(next) {
			cleanFlag := strings.ReplaceAll(flag, "-", "")
			name, hasArg := GetArgNameByTrigger(cleanFlag)

			if hasArg {
				SetArgValue(name, next)
			}
			i++
			continue
		} else if isKey && (hasKey(next) || !hasNext) {
			cleanFlag := strings.ReplaceAll(flag, "-", "")
			name, hasArg := GetArgNameByTrigger(cleanFlag)

			if hasArg {
				SetArgValue(name, true)
			}
			continue
		}
	}
	Parsed = true
}

func parseFlags() []string {
	return os.Args
}

func GetArgNameByTrigger(trigger string) (string, bool) {
	for _, a := range Arguments {
		for i := 0; i < len(a.Triggers); i++ {
			if a.Triggers[i] == trigger {
				return a.Name, true
			}
		}
	}
	return "", false
}

func hasKey(flag string) bool {
	length := len(flag)
	// Check for both non-verbose and verbose keys
	if length > 0 && flag[0] == '-' || length > 1 && flag[0:1] == "--" {
		return true
	}
	return false
}
