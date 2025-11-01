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
	plainTextFilePath          = "./message.txt"
	plainTextDecryptedFilePath = "./message_decrypted.txt"
	cipherTextFilePath         = "./CIPHER.TXT"
)

//go:generate GOOS=linux   GOARCH=amd64 go build -o ./sample/transmission ./*.go
//go:generate GOOS=windows GOARCH=amd64 go build -o ./sample/transmission.exe ./*.go
func main() {
	go start()

	waitCtrlC()
}

func start() {
	for {
		log.Println("> Enter Your Command. ('h' for help)")
		inputStr, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			log.Println("> Read input failed, error: ", err)
			continue
		}

		matched := regexp.MustCompile(`(\w+)`).FindString(strings.ToLower(inputStr))
		switch matched {
		case "h", "help":
			printHelp()
		case "i":
			initialize()
		case "g":
			generateKeypair()
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

func initialize() {
	err := os.WriteFile(plainTextFilePath, []byte("test plain text"), 0777)
	if err != nil {
		log.Println("> Initialize message file failed, error: ", err)
		return
	}

	log.Println("> Initialize message file success.")
}

func printHelp() {
	log.Println("> Options:")
	log.Println("  - h: this help")
	log.Println(fmt.Sprintf("  - i: initialize message file ('%s')", plainTextFilePath))
	log.Println(fmt.Sprintf("  - g: generate public & private key into files ('%s' & '%s')",
		publicKeyFilePath, privateKeyFilePath))
	log.Println(fmt.Sprintf("  - e: encrypt plain text from '%s' and output cipher text to '%s'",
		plainTextFilePath, cipherTextFilePath))
	log.Println(fmt.Sprintf("  - d: decrypt cipher text from '%s' and output plain text to '%s'",
		cipherTextFilePath, plainTextDecryptedFilePath))
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

	log.Println("> Exit.")
}
