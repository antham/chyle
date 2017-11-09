package cmd

import (
	"fmt"

	"github.com/fatih/color"
)

var failure = func(err error) {
	color.Red(err.Error())
}

var println = func(data string) {
	fmt.Println(data)
}
