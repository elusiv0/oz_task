package tools

import (
	"errors"
	"fmt"

	_ "github.com/99designs/gqlgen"
)

func some() {
	errors.Unwrap()
	fmt.Errorf()
}
