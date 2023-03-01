package exceptions

import "fmt"

type LinkAlreadyExistsError struct {
	Value string
}

func (lae *LinkAlreadyExistsError) Error() string {
	return fmt.Sprintf("record for already exists: %s", lae.Value)
}
