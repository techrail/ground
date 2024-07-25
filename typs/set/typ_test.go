package set

import (
	"fmt"
	"testing"
)

func TestTyp_IntegersOnly(t *testing.T) {
	intSet := New[int]()
	intSet.Add(1)
	intSet.Add(1)
	intSet.Add(1)
	intSet.Add(2)
	intSet.Add(2)
	intSet.Add(3)
	intSet.Add(4)
	intSet.Add(5)

	expectedCount := 5
	count := intSet.Size()
	if count != expectedCount {
		t.Errorf("E#1PM2UA - Expected %v elements, found %v", expectedCount, count)
	}

	valToSearch := 2
	if !intSet.Contains(valToSearch) {
		t.Errorf("E#1PM30G - Expected to contain %v but not found", valToSearch)
	}

	valToSearch = 10
	if intSet.Contains(10) {
		t.Errorf("E#1PM3CI - Did not expect %v to be there but it was there.", valToSearch)
	}

	intSet.Remove(10)

	expectedCount = 5
	count = intSet.Size()
	if count != expectedCount {
		t.Errorf("E#1PM3ED - Expected %v elements, found %v", expectedCount, count)
	}

	intSet.Remove(1)
	expectedCount = 4
	count = intSet.Size()
	if count != expectedCount {
		t.Errorf("E#1PM3FW - Expected %v elements, found %v", expectedCount, count)
	}

	intSet.Empty()
	expectedCount = 0
	count = intSet.Size()
	if count != expectedCount {
		t.Errorf("E#1S9F7Z - Expected %v elements, found %v", expectedCount, count)
	} else {
		fmt.Printf("I#1S9F97 - As expected: %v\n", expectedCount)
	}

	intSet.Add(123)
	expectedCount = 1
	count = intSet.Size()
	if count != expectedCount {
		t.Errorf("E#1S9FBO - Expected %v elements, found %v", expectedCount, count)
	} else {
		fmt.Printf("I#1S9FBR - As expected: %v\n", expectedCount)
	}
}

func TestTyp_IntSetFromSlice(t *testing.T) {
	slice := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	fmt.Printf("slice: %v\n", slice)

	intSetFromSlice := NewFromSlice(slice)

	expectedCount := 10
	count := intSetFromSlice.Size()

	if count != expectedCount {
		t.Errorf("E#20LD0U - Expected %v elements, found %v", expectedCount, count)
	} else {
		fmt.Printf("I#20LD11 - As expected: %v\n", expectedCount)
	}

	valToSearch := 2
	if !intSetFromSlice.Contains(valToSearch) {
		t.Errorf("E#20LDC7 - Expected to contain %v but not found", valToSearch)
	}

	valToSearch = 20
	if intSetFromSlice.Contains(20) {
		t.Errorf("E#20LDCK- Did not expect %v to be there but it was there.", valToSearch)
	}

	intSetFromSlice.Remove(20)

	expectedCount = 10
	count = intSetFromSlice.Size()
	if count != expectedCount {
		t.Errorf("E#20LDCX - Expected %v elements, found %v", expectedCount, count)
	}

	intSetFromSlice.Remove(1)
	expectedCount = 9
	count = intSetFromSlice.Size()
	if count != expectedCount {
		t.Errorf("E#20LDD7 - Expected %v elements, found %v", expectedCount, count)
	}

	intSetFromSlice.Empty()
	expectedCount = 0
	count = intSetFromSlice.Size()
	if count != expectedCount {
		t.Errorf("E#20LDDJ - Expected %v elements, found %v", expectedCount, count)
	} else {
		fmt.Printf("I#20LDDX - As expected: %v\n", expectedCount)
	}

	intSetFromSlice.Add(123)
	expectedCount = 1
	count = intSetFromSlice.Size()
	if count != expectedCount {
		t.Errorf("E#20LDEH - Expected %v elements, found %v", expectedCount, count)
	} else {
		fmt.Printf("I#21LDEM - As expected: %v\n", expectedCount)
	}
}
