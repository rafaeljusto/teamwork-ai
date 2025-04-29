package teamwork

// Ref is a utility function that returns a pointer to the value of type T.
func Ref[T any](v T) *T {
	return &v
}
