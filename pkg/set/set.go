package set

type void struct{}

var exists void

type set struct {
	val map[string]void
}

type Set interface {
	Add(str string)
	Remove(str string)
	Contains(str string) bool
	Values() []string
}

// Initialize a new set
func New() Set {
	return &set{
		val: make(map[string]void),
	}
}

// Add new value to set
func (s *set) Add(str string) {
	s.val[str] = exists
}

// Remove value from set
func (s *set) Remove(str string) {
	delete(s.val, str)
}

// Check if set contains a value
func (s *set) Contains(str string) bool {
	_, ok := s.val[str]
	return ok
}

// Return all the values from set as a slice of string
func (s *set) Values() []string {
	values := make([]string, len(s.val))

	i := 0
	for v := range s.val {
		values[i] = v
		i++
	}

	return values
}
