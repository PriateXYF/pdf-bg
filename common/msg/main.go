package msg

import (
	"github.com/fatih/color"
	"fmt"
	// "strings"
)

func Show(format string, args ...interface{}){
	fmt.Printf(format, args...)
}

func Info(format string, args ...interface{}){
    color.Blue("ℹ️  " + format, args...)
}

func Fail(format string, args ...interface{}){
    color.Red("🔴 " + format, args...)
}

func Warn(format string, args ...interface{}){
    color.Yellow("⚠️  " + format, args...)
}

func Good(format string, args ...interface{}){
    color.Green("✅ " + format, args...)
}