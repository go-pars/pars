package pars

// Child will map to the i'th child of the result.
func Child(i int) Map {
	return func(result *Result) error {
		if result.Children == nil {
			panic("result does not have children")
		}
		*result = result.Children[i]
		return nil
	}
}
