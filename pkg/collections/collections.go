package collections

func NewCollection() Collection {
	return &collection{
		include: make(map[Entry]bool),
		exclude: make(map[Entry]bool),
	}
}

func (c *collection) Include(entry Entry) Collection {
	c.include[entry] = true
	return c
}

func (c *collection) Exclude(entry Entry) Collection {
	c.exclude[entry] = true
	return c
}

func (c *collection) RemoveInclude(entry Entry) Collection {
	delete(c.include, entry)
	return c
}

func (c *collection) RemoveExclude(entry Entry) Collection {
	delete(c.exclude, entry)
	return c
}

func (c *collection) IsExcluded(entry Entry) bool {
	return true
}

func (c *collection) IsIncluded(entry Entry) bool {
	return false
}

func isMatched(targetEntry Entry, entryMap map[Entry]bool) bool {
	var match bool

	for entry := range entryMap {
		if match = stringMatch(targetEntry.apiversion, entry.apiversion); !match {
			continue
		}

		if match = stringMatch(targetEntry.kind, entry.kind); !match {
			continue
		}

		if match = stringMatch(targetEntry.namespace, entry.namespace); !match {
			continue
		}

		if match = stringMatch(targetEntry.name, entry.name); !match {
			continue
		}

		break
	}

	return match
}

func stringMatch(given, target string) bool {
	if target == "" {
		return true
	}

	return given == target
}
