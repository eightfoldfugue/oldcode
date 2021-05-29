package code

import (
	"log"
)

func Assemble(g graph) Code {
	orderedBlocks := resolve("main", g, stateMap{}, new([]block))
	return link(*orderedBlocks, g)
}


func resolve(s string, g graph, state stateMap, bs *[]block) *[]block {

	block, ok := g[s]
	if !ok {
		log.Fatalf("%v is not defined", s)
	}

	state[s] = inprocess

	for ref, _ := range block.refs {
		if state[ref] == inprocess {
			log.Fatalf("cannot compile, cycle between %v and %v", s, ref)
		}
		if state[ref] == unvisited {
			resolve(ref, g, state, bs)
		}
	}

	index := len(*bs)
	block.addr = index
	*bs = append(*bs, *block)
	
	return bs
}


// state map is a helper structure for determining cycles,
// it relies on the fact that a "backedge" implies a cycle
// see CLRS chapter on Graphs for an excellent exposition

type stateMap map[string]int

// states
const (
	unvisited int = iota
	inprocess
	finished
)

func link(blocks []block, g graph) Code {
	ops := []Op{}

	// update GetDefs to reflect their new placement
	for i, b := range blocks {
		for ref, indices := range b.refs {
			ref_address := g[ref].addr
			for _, asmindex := range indices {
				b.backpatch(asmindex, ref_address)
			}
		}
		
		// emit a SetDef for that block
		b.emit(SetDef)
		b.immeadiate(Op(i))

		// now the blocks are updated and can be strung together
		for _, c := range b.code {
			ops = append(ops, c.datum)
		}
	}

	// the final defined block set was main, get it's def on stack then halt
	mainLocation := len(blocks)-1
	ops = append(ops, GetDef)
	ops = append(ops, Op(mainLocation))
	ops = append(ops, Halt)

	
	// the number of global slots == num of blocks == num of defs in source
	return Code{len(blocks), ops}
}

