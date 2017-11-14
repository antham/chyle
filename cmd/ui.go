package cmd

import (
	"fmt"

	"github.com/fatih/color"
)

var failure = func(err error) {
	color.Red(err.Error())
}

var printWithNewLine = func(str string) {
	fmt.Println(str)
}
