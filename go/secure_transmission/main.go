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
		Info("Enter Your Command ('h' for help)")

		if !scanner.Scan() {
			break
		}

		text := strings.TrimSpace(scanner.Text())
		matched := regexp.MustCompile(`(\w+)`).FindString(strings.ToLower(text))
		switch matched {
		case "h", "help":
			printHelp()
		case "g", "gen", "generate":
			generateKeypair()
		case "i", "init", "initialize":
			initMessageFile()
		case "e", "encrypt":
			encrypt()
		case "d", "decrypt":
			decrypt()
		case "exit":
			Info("Exit")
			break ALL
		default:
			Info("Unknown input: " + text + ", 'h' for help")
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
	Info(`Options:
  - h: this help
  - g: generate public & private key into files ('./priv.key' & './PUB.KEY')
  - i: initialize message file './message.txt'
  - e: encrypt plain text from './message.xxx' and write cipher to './CIPHER.XXX'
  - d: decrypt cipher from './CIPHER.XXX' and write plain text to './message_decrypted.xxx'
  - exit: exit program

note: encrypt/decrypt support automatic recognize 'file extension', in fact: 
when encrypt, we find first file which name matched 'message.[xxx]' and encrypt it into 'CIPHER.[XXX]';
when decrypt, we find first file which name matched 'CIPHER.[XXX]' and decrypt it into 'message_decrypted.[xxx]'`)
}
