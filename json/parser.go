package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func GenerateMap(s string) map[string]interface{} {
	var st Stack[map[string]interface{}]
	m := make(map[string]interface{})
	var w string
	current_map := m
	for i := range s {
		if s[i] == '{' {
			if i != 0 {
				st.Push(current_map)
				z := make(map[string]interface{})
				current_map[w] = z
				current_map = z
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
			if s[i] != '{' {
				var val string = ""
				for s[i] != ',' {
					if s[i] != '"' {
						val += string(s[i])
					}
					i++
				}
				i++
				current_map[w] = strings.TrimSpace(val)
			} else {
				continue
			}
		} else if s[i] == '}' && !st.Empty() {
			current_map = st.Top()
			st.Pop()
			w = ""
		} else if s[i] == ' ' || s[i] == '\t' || s[i] == '"' {
			continue
		} else if s[i] == ',' {
			w = ""
		} else {
			w += string(s[i])
		}
	}
	return m
}

func PrettyPrint(m map[string]interface{}) {
	b, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Print(string(b))
}

func main() {
	PrettyPrint(GenerateMap("{ aditya : { you : {lmao : \"what\", bier : {what : 5, who : 3, } , akldfj : 5, }, what : 3, } , bitch : 1,}"))
}
