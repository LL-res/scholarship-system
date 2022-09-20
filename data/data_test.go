package data

import (
	"fmt"
	"testing"
)

func TestNewStudents(t *testing.T) {
	NewStudents()
	for k, v := range DataBase {
		fmt.Println(k, *v)
	}
}
