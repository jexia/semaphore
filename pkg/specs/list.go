package specs

// PropertyList represents a list of properties
type PropertyList []*Property

func (list PropertyList) Len() int           { return len(list) }
func (list PropertyList) Swap(i, j int)      { list[i], list[j] = list[j], list[i] }
func (list PropertyList) Less(i, j int) bool { return list[i].Position < list[j].Position }

// Get attempts to return a property inside the given list with the given name
func (list PropertyList) Get(key string) *Property {
	for _, item := range list {
		if item == nil {
			continue
		}

		if item.Name == key {
			return item
		}
	}

	return nil
}
