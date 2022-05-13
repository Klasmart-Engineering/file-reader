package errors

import "fmt"

type CSVRowError struct {
	RowNum int
	What   string
}

func (e CSVRowError) Error() string {
	return fmt.Sprintf("Row number %v failed to process: %v", e.RowNum, e.What)
}
