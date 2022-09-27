package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func GenerateMap(s string) interface{} {
	var st Stack[interface{}]
	var m interface{}
	if s[0] == '[' {
		b := make([]interface{}, 0)
		m = &b
	} else if s[0] == '{' {
		m = make(map[string]interface{})
	}
	var w string
	var cMap interface{} = m
	for i := range s {
		if s[i] == '{' || s[i] == '[' {
			if i != 0 {
				st.Push(cMap)
				var z interface{}
				if s[i] == '{' {
					z = make(map[string]interface{})
				} else {
					b := make([]interface{}, 0)
					z = &b
				}
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
				cMap = z
				w = ""
				i++
				for s[i] == ' ' {
					i++
				}
			} else {
				continue
			}
		} else if s[i] == ':' {
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
						cMap.(map[string]interface{})[w] = strings.TrimSpace(val)
					}
				case *[]interface{}:
					{
						*cMap.(*[]interface{}) = append(*cMap.(*[]interface{}), val)
					}
				}
			} else {
				continue
			}
		} else if (s[i] == '}' || s[i] == ']') && !st.Empty() {
			switch cMap.(type) {
			case *[]interface{}:
				{
					if s[i] == ']' {
						*cMap.(*[]interface{}) = append(*cMap.(*[]interface{}), w)
					}
				}
			}
			cMap = st.Top()
			st.Pop()
			w = ""
		} else if s[i] == ' ' || s[i] == '\t' || s[i] == '"' || s[i] == '\n' {
			continue
		} else if s[i] == ',' {
			switch cMap.(type) {
			case *[]interface{}:
				{
					if w != "" {
						*cMap.(*[]interface{}) = append(*cMap.(*[]interface{}), w)
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
