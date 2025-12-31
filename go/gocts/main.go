package main

import (
	"log"

	"github.com/mats9693/study/go/gocts/generate"
	"github.com/mats9693/study/go/gocts/parse"
)

func main() {
	log.Println("> Gocts: Start.")
	defer log.Println("> Gocts: Finish.")

	parse.TraversalDir()

	gen.GenerateConfigFile()
	gen.GenerateRequestFiles()
	gen.GenerateStructureFiles()
}
