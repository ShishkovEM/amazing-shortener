package exceptions

import "fmt"

type LinkAlreadyExistsError struct {
	Value string
}

func (lae *LinkAlreadyExistsError) Error() string {
	return fmt.Sprintf("record for %s already exists", lae.Value)
}
