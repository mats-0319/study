package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
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

//go:generate GOOS=linux   GOARCH=amd64 go build -o ./sample/transmission
//go:generate GOOS=windows GOARCH=amd64 go build -o ./sample/transmission.exe
func main() {
	go start()

	waitCtrlC()
}

func start() {
	for {
		log.Println("> Enter Your Command. ('h' for help)")
		inputStr, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			ExecLog("Read input", err)
			continue
		}

		matched := regexp.MustCompile(`(\w+)`).FindString(strings.ToLower(inputStr))
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
		default:
			log.Println("> Unknown input: ", inputStr, "'h' for help")
		}

		fmt.Println()
	}
}

func initMessageFile() {
	err := os.WriteFile(plainTextFileName+defaultExtension, []byte("test plain text"), 0777)
	if err != nil {
		ExecLog("Initialize message file", err)
		return
	}

	ExecLog("Initialize message file")
}

func printHelp() {
	log.Print(`
> Options:
  - h: this help
  - g: generate public & private key into files ('./priv.key' & './PUB.KEY')
  - i: initialize message file './message.txt'
  - e: encrypt plain text from './message.xxx' and write cipher to './CIPHER.XXX'
  - d: decrypt cipher from './CIPHER.XXX' and write plain text to './message_decrypted.xxx'

note: encrypt/decrypt support automatic recognize 'file extension', in fact, 
when encrypt, we find first file which name matched 'message.[xxx]' and encrypt it into 'CIPHER.[XXX]';
when decrypt, we find first file which name matched 'CIPHER.[XXX]' and decrypt it into 'message_decrypted.[xxx]'
`)
}

func waitCtrlC() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)

	defer func() {
		signal.Stop(c)
		close(c)
	}()

	for range c {
		break
	}

	fmt.Println()
	log.Println("> Exit.")
}
