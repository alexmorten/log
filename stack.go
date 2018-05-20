package log

//Node of the Stack
type Node struct {
	data interface{}
	next *Node
}

//Stack with the ability to peek without poping of the stack
type Stack struct {
	head *Node
}

//Put a datum on the stack
func (s *Stack) Put(data ...interface{}) {
	for _, datum := range data {
		n := pools.Nodes.Get().(*Node)
		n.next = s.head
		n.data = datum
		s.head = n
	}
}

//Peek at stack head
func (s *Stack) Peek() interface{} {
	if s.head == nil {
		return nil
	}
	return s.head.data
}

//Pop head of the stack and return it
func (s *Stack) Pop() interface{} {
	if s.head == nil {
		return nil
	}
	n := s.head
	s.head = n.next

	data := n.data
	pools.Nodes.Put(n)
	return data
}

//Empty checks if the head of the stack points to nil
func (s *Stack) Empty() bool {
	return s.head == nil
}

//Flip the order of the stack
func (s *Stack) Flip() {
	if s.head == nil {
		return
	}
	var last *Node
	current := s.head
	next := current.next

	for {
		current.next = last
		last = current
		current = next
		if current == nil {
			s.head = last
			break
		}
		next = current.next
	}
}
