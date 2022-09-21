package utils

import (
	"github.com/fatih/color"
	"log"
)

func LogInfo(s1 string) {
	yellowString := color.YellowString("[info]")
	log.Println(yellowString + " " + s1)
	log.Println()
}
func LogSuc(s1 string) {
	greenString := color.GreenString("[main process]")
	log.Println(greenString + " " + s1)
	log.Println()
}
