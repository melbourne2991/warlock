package cmd

import (
	"fmt"
	"github.com/fatih/color"
)

var successHighlight func(a ...interface{}) string
var cyanHighlight func(a ...interface{}) string

func init() {
	successHighlight = color.New(color.FgGreen, color.Bold).SprintFunc()
	cyanHighlight = color.New(color.FgCyan, color.Bold).SprintFunc()
}

func underline(msg string) string {
	underlineStrSlice := make([]rune, len(msg))

	for i := 0; i < len(msg); i++ {
		underlineStrSlice[i] = '='
	}

	underlineStr := string(underlineStrSlice)
	return fmt.Sprintf("%s\n%s", msg, underlineStr)
}
