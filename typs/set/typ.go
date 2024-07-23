package set

// Typ is a generic type for a set data structure.
// It can only be used on comparable data types though
type Typ[T comparable] map[T]string

// New creates a new set.
func New[T comparable]() Typ[T] {
	return make(map[T]string)
}

// Add adds an element to the set.
func (s Typ[T]) Add(element T, withVal ...string) {
	if len(withVal) > 0 {
		s[element] = withVal[0]
	} else {
		s[element] = ""
	}
}

// Contains checks if an element is in the set.
func (s Typ[T]) Contains(element T) bool {
	_, exists := s[element]
	return exists
}

// Value will return the element value
func (s Typ[T]) Value(element T) (string, bool) {
	if s.Contains(element) {
		return s[element], true
	} else {
		return "", false
	}
}

// Remove removes an element from the set.
func (s Typ[T]) Remove(element T) {
	delete(s, element)
}

// Size returns the number of elements in the set.
func (s Typ[T]) Size() int {
	return len(s)
}

func (s Typ[T]) Empty() {
	for k, _ := range s {
		delete(s, k)
	}
}

func (s Typ[T]) AsSlice() []T {
	retVal := []T{}
	for k, _ := range s {
		retVal = append(retVal, k)
	}
	return retVal
}

func (s Typ[T]) Union(with Typ[T]) Typ[T] {
	retVal := s
	for k, v := range with {
		_, ok := s[k]
		if !ok {
			retVal[k] = v
		}
	}
	return retVal
}

func (s Typ[T]) Subtract(another Typ[T]) Typ[T] {
	retVal := New[T]()
	for k, v := range s {
		_, ok := another[k]
		if !ok {
			retVal.Add(k, v)
		}
	}
	return retVal
}
