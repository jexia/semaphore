package specs

import (
	"sort"
	"testing"
)

func TestPropertyListSort(t *testing.T) {
	list := PropertyList{
		&Property{Name: "third", Position: 2},
		&Property{Name: "first", Position: 0},
		&Property{Name: "second", Position: 1},
	}

	sort.Sort(list)

	for index, item := range list {
		if int(item.Position) != index {
			t.Fatalf("unexpected property list order %d, expected %d", item.Position, index)
		}
	}
}

func TestPropertyListGet(t *testing.T) {
	list := PropertyList{
		&Property{Name: "first"},
		&Property{Name: "second"},
	}

	result := list.Get("second")
	if result == nil {
		t.Fatal("unexpected empty result when looking up second")
	}

	unexpected := list.Get("unexpected")
	if unexpected != nil {
		t.Fatal("unexpected lookup returned a unexpected property")
	}
}

func TestPropertyListGetNil(t *testing.T) {
	list := PropertyList{
		nil,
		&Property{Name: "first"},
		nil,
		&Property{Name: "second"},
		nil,
	}

	result := list.Get("second")
	if result == nil {
		t.Fatal("unexpected empty result when looking up second")
	}

	unexpected := list.Get("unexpected")
	if unexpected != nil {
		t.Fatal("unexpected lookup returned a unexpected property")
	}
}
