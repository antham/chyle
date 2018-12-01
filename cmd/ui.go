package cmd

import (
	"fmt"
	"log"

	"github.com/fatih/color"
)

func failure(err error) {
	c := color.New(color.FgRed)
	if _, ferr := c.Fprintf(writer, "%s\n", err.Error()); ferr != nil {
		log.Fatal(ferr)
	}
}

func printWithNewLine(str string) {
	if _, err := fmt.Fprintf(writer, "%s\n", str); err != nil {
		log.Fatal(err)
	}
}
