package utils

import (
	"log"
	"os"
	"regexp"
	"sync"
)

var (
	ErrorLogger   *log.Logger
	InfoLogger    *log.Logger
	CommandLogger *log.Logger
)

var (
	HyperlinkRegex = regexp.MustCompile(`\x7c(.*?)>$`)
	MentionRegex   = regexp.MustCompile(`<@(.*?)>`)
	DescRegex      = regexp.MustCompile(`"(.*?)"$`)
)

var BotLock = &sync.Mutex{}
var GameLock = &sync.Mutex{}

func init() {
	InfoLogger = log.New(os.Stdout, "[INFO]: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "[ERROR]: ", log.Ldate|log.Ltime|log.Lshortfile)
	CommandLogger = log.New(os.Stdout, "[COMNMAND EVENT]: ", 0)
}

func ExtractTxt(regex *regexp.Regexp, in string) string {
	str := regex.FindStringSubmatch(in)
	if len(str) != 2 {
		return ""
	}
	return str[1]
}
