package main

import (
	"log"

	"github.com/mats9693/study/go/gocts/generate_ts"
	"github.com/mats9693/study/go/gocts/parse"
)

func main() {
	log.Println("> Gocts: Start.")
	defer log.Println("> Gocts: Finish.")

	parse.TraversalDir()

	generate_ts.GenerateConfig()
	generate_ts.GenerateRequestFiles()
	generate_ts.GenerateStructureFiles()
}
