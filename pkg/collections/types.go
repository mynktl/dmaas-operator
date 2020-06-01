package collections

type Entry struct {
	apiversion string
	kind       string
	name       string
	namespace  string
}

// Collection provide interface to update/check collection
type Collection interface {
	// Include add new entry in include collection
	Include(entry Entry) Collection

	// Exclude add new entry in exclude collection
	Exclude(entry Entry) Collection

	// RemoveInclude remove given entry from include collection
	RemoveInclude(entry Entry) Collection

	// RemoveExcludeEntry remove given entry from exclude collection
	RemoveExclude(entry Entry) Collection

	// IsIncluded check if given entry can be part of include collection
	// It will check the entry against include/exclude collection
	// if entry matches with exclude collection then it will return false
	// if entry doesn't match with exclude and match with include then it return true else false
	// match will be verified in following steps:
	// - check for match against field apiversion
	// - check for match against field kind
	// - check for match against field namespace
	// - check for match against field name
	// if field from collection entry is empty then match will be considered successful
	// if match fails then futher steps will be skipped
	IsIncluded(entry Entry) bool

	// IsExcluded check if given entry can be part of exclude collection
	// It will check the entry against include/exclude collection
	// if entry matches with exclude collection then it will return true
	// if entry doesn't match with exclude and match with include then it return false else true
	// match will be verified in following steps:
	// - check for match against field apiversion
	// - check for match against field kind
	// - check for match against field namespace
	// - check for match against field name
	// if field from collection entry is empty then match will be considered successful
	// if match fails then futher steps will be skipped
	IsExcluded(entry Entry) bool
}

type collection struct {
	include map[Entry]bool
	exclude map[Entry]bool
}
