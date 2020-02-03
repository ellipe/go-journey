package forms

// Defines an errors type that holds the error information for each field in the form
type errors map[string][]string

// Add : function that add error messages to a given field to the map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get : Retrieves first error message for a given field from the map.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
