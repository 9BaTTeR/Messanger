package Utility

type JSON interface {
	Parse(source string) error
	Compose() ([]byte, error)
}
