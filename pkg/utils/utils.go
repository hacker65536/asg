package utils

import "fmt"

type Color string

const (
	Reset             Color = "\x1b[0000m"
	Bright            Color = "\x1b[0001m"
	BlackText         Color = "\x1b[0030m"
	RedText           Color = "\x1b[0031m"
	GreenText         Color = "\x1b[0032m"
	YellowText        Color = "\x1b[0033m"
	BlueText          Color = "\x1b[0034m"
	MagentaText       Color = "\x1b[0035m"
	CyanText          Color = "\x1b[0036m"
	WhiteText         Color = "\x1b[0037m"
	DefaultText       Color = "\x1b[0039m"
	BrightRedText     Color = "\x1b[1;31m"
	BrightGreenText   Color = "\x1b[1;32m"
	BrightYellowText  Color = "\x1b[1;33m"
	BrightBlueText    Color = "\x1b[1;34m"
	BrightMagentaText Color = "\x1b[1;35m"
	BrightCyanText    Color = "\x1b[1;36m"
	BrightWhiteText   Color = "\x1b[1;37m"
)

func (c *Color) String() string {
	return fmt.Sprintf("%#v", c)
}

func Paint(color Color, value string) string {
	return fmt.Sprintf("%v%v%v", color, value, Reset)
}
func Red(value string) string {
	return fmt.Sprintf("%v%v%v", RedText, value, Reset)
}

func Normal(value string) string {
	return fmt.Sprintf("%v%v%v", DefaultText, value, Reset)
}

func Yellow(value string) string {
	return fmt.Sprintf("%v%v%v", YellowText, value, Reset)
}

func Green(value string) string {
	return fmt.Sprintf("%v%v%v", GreenText, value, Reset)
}
