package set

import "testing"

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
}
