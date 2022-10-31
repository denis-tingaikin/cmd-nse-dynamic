package main

import (
	"bufio"
	"strings"
)

type Node struct {
	Metadata []string
	Body     []*Node
	Name     string
}

func Parse(source string) []*Node {
	var s = bufio.NewScanner(strings.NewReader(source))

	s.Split(bufio.ScanWords)

	var stack []*Node
	var current = new(Node)
	var metadata []string

	for s.Scan() {
		if s.Text() == "{" {
			stack = append(stack, current)
			current = current.Body[len(current.Body)-1]
			continue
		}
		if s.Text() == "}" {
			current = stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			continue
		}
		var name = s.Text()
		if strings.HasPrefix(name, "'") && strings.HasSuffix(name, "':") || strings.HasPrefix(name, "\"") && strings.HasSuffix(name, "\":") {
			metadata = append(metadata, name[1:len(name)-2])
			continue
		}

		current.Body = append(current.Body, &Node{
			Name:     name,
			Metadata: metadata,
		})
		metadata = nil

	}

	return current.Body
}
