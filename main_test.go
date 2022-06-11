package main

import (
	"fmt"
	"testing"
)

func TestMain(t *testing.T) {
	fmt.Println(fill_path("/tmp/a.txt", "."))
}
