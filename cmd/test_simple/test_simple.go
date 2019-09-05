package main

import (
	"fmt"

	"github.com/microamp/idgen/idgen"
)

func main() {
	var nodeID uint64 = 0
	idGen, err := idgen.NewIDGen(nodeID)
	if err != nil {
		panic(err)
	}

	ids := make(map[idgen.ID]bool)

	sampleSize := 1000000
	for i := 0; i < sampleSize; i++ {
		id := idGen.GenID()
		ids[id] = true
	}

	fmt.Printf("%d unique IDs generated\n", len(ids))
}
