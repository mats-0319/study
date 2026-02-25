package main

import (
	"bufio"
	"log"
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
		log.Println("> Enter Your Command. ('h' for help)")

		if !scanner.Scan() {
			break
		}

		text := strings.TrimSpace(scanner.Text())
		matched := regexp.MustCompile(`(\w+)`).FindString(strings.ToLower(text))
		switch matched {
		case "h", "help":
			printHelp()
		case "g":
			generateKeypair()
		case "i":
			initMessageFile()
		case "e":
			encrypt()
		case "d":
			decrypt()
		case "exit":
			log.Println("> Exit.")
			break ALL
		default:
			log.Println("> Unknown input: ", text, ", 'h' for help")
		}
	}
}

func initMessageFile() {
	err := os.WriteFile(plainTextFileName+defaultExtension, []byte("test plain text"), 0777)
	if err != nil {
		Log("Initialize message file", err)
		return
	}

	Log("Initialize message file")
}

func printHelp() {
	log.Print(`> Options:
  - h: this help
  - g: generate public & private key into files ('./priv.key' & './PUB.KEY')
  - i: initialize message file './message.txt'
  - e: encrypt plain text from './message.xxx' and write cipher to './CIPHER.XXX'
  - d: decrypt cipher from './CIPHER.XXX' and write plain text to './message_decrypted.xxx'
  - exit: exit program

note: encrypt/decrypt support automatic recognize 'file extension', in fact: 
when encrypt, we find first file which name matched 'message.[xxx]' and encrypt it into 'CIPHER.[XXX]';
when decrypt, we find first file which name matched 'CIPHER.[XXX]' and decrypt it into 'message_decrypted.[xxx]'
`)
}
