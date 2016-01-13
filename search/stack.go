package search

type Stack struct {
	top  *Element
	size int
}

type Element struct {
	value string
	next  *Element
}

func (s *Stack) Len() int {
	return s.size
}

func (s *Stack) Push(value string) {
	s.top = &Element{value, s.top}
	s.size++
}

func (s *Stack) Pop() (value string) {
	if s.size > 0 {
		value, s.top = s.top.value, s.top.next
		s.size--
	}
	return
}

func (s *Stack) Top() (value string) {
	if s.size > 0 {
		value = s.top.value
	}
	return
}
