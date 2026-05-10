package main

import (
	"log"

	"github.com/mats0319/study/go/gocts/printer"
	"github.com/mats0319/study/go/gocts/scanner"
)

func main() {
	log.Println("> Gocts: Start...")
	defer log.Println("> Gocts: Done.")

	scanner.TraversalDir()

	gen.GenerateConfigFile()
	gen.GenerateRequestFiles()
	gen.GenerateStructureFiles()
}
