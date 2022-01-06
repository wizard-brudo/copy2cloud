package main

const (
	RED    = "\u001b[31m"
	YELLOW = "\u001b[33m"
	GREEN  = "\u001b[32m"
	RESET  = "\u001b[0m"
)

func setTextColor(text, color string) string {
	return color + text + RESET
}
