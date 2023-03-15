package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	inQuotes bool = false
)

func ParseVal(val interface{}) interface{} {
	if !inQuotes {
		s := val.(string)
		if s[0] <= '9' && s[0] >= '0' {
			if strings.Contains(s, ".") {
				if v, err := strconv.ParseFloat(s, 64); err == nil {
					return v
				}
			} else if v, err := strconv.ParseInt(s, 10, 64); err == nil {
				return v
			}
		} else if s == "false" {
			return false
		} else if s == "true" {
			return true
		} else if s == "null" {
			return nil
		}
		return s
	} else {
		return val
	}
}

func GenerateMap(s string) interface{} {
	var (
		st Stack[interface{}]
		m  interface{}
		w  string
		// current object
		cm interface{}
	)
	for i := range s {
		if s[i] == '{' || s[i] == '[' {
			var z interface{}
			if s[i] == '{' {
				z = make(map[string]interface{})
			} else {
				b := make([]interface{}, 0)
				z = &b
			}
			if i == 0 {
				m = z
			} else {
				st.Push(cm)
				switch cm.(type) {
				case map[string]interface{}:
					{
						cm.(map[string]interface{})[w] = z
					}
				case *[]interface{}:
					{
						*cm.(*[]interface{}) = append(*cm.(*[]interface{}), z)
					}
				}
			}
			cm = z
			w = ""
			i++
			for s[i] == ' ' {
				i++
			}
		} else if s[i] == ':' {
			if !inQuotes {
				for s[i] != ' ' {
					i++
				}
				if s[i] != '{' && s[i] != '[' {
					var val string = ""
					for s[i] != ',' {
						if s[i] == ']' || s[i] == '}' {
							break
						}
						if s[i] != '"' {
							val += string(s[i])
						}
						i++
					}
					i++
					switch cm.(type) {
					case map[string]interface{}:
						{
							cm.(map[string]interface{})[w] = ParseVal(strings.TrimSpace(val))
						}
					case *[]interface{}:
						{
							*cm.(*[]interface{}) = append(*cm.(*[]interface{}), ParseVal(val))
						}
					}
				} else {
					continue
				}
			}
		} else if (s[i] == '}' || s[i] == ']') && !st.Empty() {
			switch cm.(type) {
			case *[]interface{}:
				{
					// Adding Last Element to the array.
					if s[i] == ']' && len(w) != 0 {
						*cm.(*[]interface{}) = append(*cm.(*[]interface{}), ParseVal(w))
					}
				}
			}
			cm = st.Top()
			st.Pop()
			w = ""
		} else if s[i] == ' ' || s[i] == '\t' || s[i] == '"' || s[i] == '\n' {
			if s[i] == '"' {
				inQuotes = !inQuotes
			}
			continue
		} else if s[i] == ',' {
			switch cm.(type) {
			case *[]interface{}:
				{
					if w != "" {
						*cm.(*[]interface{}) = append(*cm.(*[]interface{}), ParseVal(w))
					}
				}
			}
			w = ""
		} else {
			w += string(s[i])
		}
	}
	return m
}

func PrettyPrint(m interface{}) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("No File Provided!")
		os.Exit(-1)
	} else {
		path := os.Args[1]
		content, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("No File Found at path : " + path)
			os.Exit(-1)
		}
		m := GenerateMap(string(content))
		PrettyPrint(m)
	}
}
