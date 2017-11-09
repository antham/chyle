package cmd

import (
	"github.com/fatih/color"
)

var failure = func(err error) {
	color.Red(err.Error())
}
