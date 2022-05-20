package reader

type Reader interface {
	Read() ([]string, error)
}
