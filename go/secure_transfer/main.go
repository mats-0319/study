package main

import (
	"bufio"
	"os"
	"regexp"
	"strings"
)

const (
	privateKeyFilePath         = "./priv.key"
	publicKeyFilePath          = "./PUB.KEY"
	plainTextFileName          = "message"
	cipherFileName             = "CIPHER"
	plainTextDecryptedFileName = "message_decrypted"
	defaultExtension           = ".txt"
)

func main() {
	start()
}

func start() {
	scanner := bufio.NewScanner(os.Stdin)

ALL:
	for {
		Info("Enter Your Command ('h' for help) .")

		if !scanner.Scan() {
			break
		}

		text := strings.TrimSpace(scanner.Text())
		matched := regexp.MustCompile(`(\w+)`).FindString(strings.ToLower(text))
		switch matched {
		case "h", "help", "0":
			printHelp()
		case "g", "gen", "generate", "1":
			generateKeypair()
		case "i", "init", "initialize", "2":
			initMessageFile()
		case "e", "encrypt", "3":
			encrypt()
		case "d", "decrypt", "4":
			decrypt()
		case "exit", "9":
			Info("Exit.")
			break ALL
		default:
			Info("Unknown input: " + text + ", 'h' for help.")
		}
	}
}

func initMessageFile() {
	err := os.WriteFile(plainTextFileName+defaultExtension, []byte("Write your message here"), 0777)
	if err != nil {
		Error("Initialize message file", err)
		return
	}

	Success("Initialize message file")
}

func printHelp() {
	Info(`
Brief: (number)
  - [1]: generate key pair
  - [2]: initialize message file
  - [3]: encrypt
  - [4]: decrypt
  - [0]: help
  - [9]: exit

Options:
  - h help: this help
  - g: generate public & private key into files ('./priv.key' & './PUB.KEY')
  - i: initialize message file './message.txt'
  - e: encrypt plain text from './message.xxx' and write cipher to './CIPHER.XXX'
  - d: decrypt cipher from './CIPHER.XXX' and write plain text to './message_decrypted.xxx'
  - exit: exit program
`)
}
