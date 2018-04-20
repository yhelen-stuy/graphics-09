package main

import (
	"bytes"
	// "fmt"
)

type Stack struct {
	stack []*Matrix
}

func MakeStack() *Stack {
	s := make([]*Matrix, 0)
	m := MakeMatrix(4, 4)
	m.Ident()
	s = append(s, m)
	return &Stack{stack: s}
}

func (s *Stack) Push(m *Matrix) {
	s.stack = append(s.stack, m)
}

func (s *Stack) Pop() *Matrix {
	if s.isEmpty() {
		return nil
	}
	l := len(s.stack)
	popped := s.stack[l-1]
	s.stack = s.stack[:l-1]
	return popped
}

func (s *Stack) isEmpty() bool {
	return len(s.stack) == 0
}

func (s *Stack) Peek() *Matrix {
	if s.isEmpty() {
		return nil
	}
	l := len(s.stack)
	return s.stack[l-1]
}

func (s Stack) String() string {
	var buf bytes.Buffer
	for _, m := range s.stack {
		buf.WriteString(m.String())
		buf.WriteString("\n")
	}
	return buf.String()
}
