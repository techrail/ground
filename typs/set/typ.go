package set

import (
	"fmt"
)

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

func main() {
	// Example usage of the Typ data type with integers.
	intSet := New[int]()
	intSet.Add(1)
	intSet.Add(2)
	intSet.Add(3)

	fmt.Println("Typ:", intSet)
	fmt.Println("Contains 2:", intSet.Contains(2))
	fmt.Println("Contains 4:", intSet.Contains(4))
	fmt.Println("Size:", intSet.Size())
}
