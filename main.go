package main

import (
	"fmt"

	"github.com/pmq24/kp/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
