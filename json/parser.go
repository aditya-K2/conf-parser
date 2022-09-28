package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	var st Stack[interface{}]
	var m interface{}
	var w string
	var cMap interface{}
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
				st.Push(cMap)
				switch cMap.(type) {
				case map[string]interface{}:
					{
						cMap.(map[string]interface{})[w] = z
					}
				case *[]interface{}:
					{
						*cMap.(*[]interface{}) = append(*cMap.(*[]interface{}), z)
					}
				}
			}
			cMap = z
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
					switch cMap.(type) {
					case map[string]interface{}:
						{
							cMap.(map[string]interface{})[w] = ParseVal(strings.TrimSpace(val))
						}
					case *[]interface{}:
						{
							*cMap.(*[]interface{}) = append(*cMap.(*[]interface{}), ParseVal(val))
						}
					}
				} else {
					continue
				}
			}
		} else if (s[i] == '}' || s[i] == ']') && !st.Empty() {
			switch cMap.(type) {
			case *[]interface{}:
				{
					// Adding Last Element to the array.
					if s[i] == ']' {
						*cMap.(*[]interface{}) = append(*cMap.(*[]interface{}), ParseVal(w))
					}
				}
			}
			cMap = st.Top()
			st.Pop()
			w = ""
		} else if s[i] == ' ' || s[i] == '\t' || s[i] == '"' || s[i] == '\n' {
			if s[i] == '"' {
				inQuotes = !inQuotes
			}
			continue
		} else if s[i] == ',' {
			switch cMap.(type) {
			case *[]interface{}:
				{
					if w != "" {
						*cMap.(*[]interface{}) = append(*cMap.(*[]interface{}), ParseVal(w))
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
	var path string = ""
	if len(os.Args) < 2 {
		fmt.Println("No File Provided!")
		os.Exit(-1)
	} else {
		path = os.Args[1]
	}
	var s string
	if path != "" {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			fmt.Println("No File Found at path : " + path)
			os.Exit(-1)
		}
		s = string(content)
		m := GenerateMap(s)
		PrettyPrint(m)
	} else {
		fmt.Println("Path is Empty!")
	}
}
