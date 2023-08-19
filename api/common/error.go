package common

type Error struct {
	container string
}

func (s *Error) Error() string {
	return s.container
}
func (s *Error) Join(errs ...error) {
	if len(errs) > 0 {
		for _, err := range errs {
			if s.container[len(s.container)] != '\n' {
				s.container += "\n"
			}
			s.container += err.Error()
		}
	}
}
func (s *Error) Clear() {
	s.container = ""
}
