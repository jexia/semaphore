package http

// UniqueStringItems collects strings to be returned as a list of unique items.
// Note that it is not suitable for concurrent usage and does not guartantee the
// order of items.
type UniqueStringItems map[string]struct{}

// Add item to the list.
func (usi UniqueStringItems) Add(item string) {
	usi[item] = struct{}{}
}

// Get the list of unique items.
func (usi UniqueStringItems) Get() []string {
	list := make([]string, 0, len(usi))

	for key := range usi {
		list = append(list, key)
	}

	return list
}
