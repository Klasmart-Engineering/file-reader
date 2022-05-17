package errors

import "fmt"

type CSVRowError struct {
	RowNum  int
	Message string
}

type InvalidRowError struct {
	Message string
}

func (e CSVRowError) Error() string {
	return fmt.Sprintf("Row number %v failed to process: %v", e.RowNum, e.Message)
}

func (e InvalidRowError) Error() string {
	return fmt.Sprintf("Error: %s", e.Message)
}
