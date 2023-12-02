package internal

import "fmt"

func valueHash(name, value string) string {
	return fmt.Sprintf("%s:%s", name, value)
}
