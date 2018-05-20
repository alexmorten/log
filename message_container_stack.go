package log

//MessageContainerStack ...
type MessageContainerStack struct {
	Stack
}

//PeekMessageContainer peeks at the stack head
func (s *MessageContainerStack) PeekMessageContainer() MessageContainer {
	container := s.Peek()
	if container == nil {
		return nil
	}
	return container.(MessageContainer)
}

//PopMessageContainer pops the stack head and returns it
func (s *MessageContainerStack) PopMessageContainer() MessageContainer {
	container := s.Pop()
	if container == nil {
		return nil
	}

	return container.(MessageContainer)
}

//PutMessageContainer onto the Stack
func (s *MessageContainerStack) PutMessageContainer(containers ...MessageContainer) {
	for _, container := range containers {
		s.Put(container)
	}
}
