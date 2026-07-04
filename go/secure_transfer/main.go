package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"

	"github.com/mats0319/study/go/transmission/internal"
)

func main() {
	start()
}

func start() {
	scanner := bufio.NewScanner(os.Stdin)

ALL:
	for {
		internal.Info("Enter Your Command ('h' for help) .")

		if !scanner.Scan() {
			break
		}

		text := strings.TrimSpace(scanner.Text())
		matched := regexp.MustCompile(`(\w+)`).FindString(strings.ToLower(text))
		switch matched {
		case "h", "help":
			printHelp()
		case "g", "gen", "generate":
			internal.GenerateKeypair()
		case "i", "init", "initialize":
			internal.InitMessageFile()
		case "e", "encrypt":
			internal.Encrypt()
		case "d", "decrypt":
			internal.Decrypt()
		case "exit", "q":
			internal.Info("Exit.")
			break ALL
		default:
			internal.Info("Unknown input: \"" + text + "\", 'h' for help.")
		}
	}
}

func printHelp() {
	internal.Info(`Options:
  - h: this help
  - g: generate public & private key into files ('./priv.key' & './PUB.KEY')
  - i: initialize message file ('./message.txt')
  - e: encrypt plain text from './message.xxx' and write cipher to './CIPHER.XXX'
  - d: decrypt cipher from './CIPHER.XXX' and write plain text to './message_decrypted.xxx'
  - exit: exit program
`)
}
