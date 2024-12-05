package schema

type Parser interface {
	Parse(path string) (*Stack, error)
}
