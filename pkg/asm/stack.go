package asm

type stack struct {
	vals []bool
}

func (s *stack) push(v bool) {
	s.vals = append(s.vals, v)
}

func (s *stack) len() int {
	return len(s.vals)
}

func (s *stack) top() bool {
	return s.vals[len(s.vals)-1]
}

func (s *stack) pop() bool {
	res := s.vals[len(s.vals)-1]
	s.vals = s.vals[0 : len(s.vals)-1]
	return res
}
