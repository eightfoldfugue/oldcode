package obj

/* package obj contains runtime representation of objects
   located in its own dir to be importable by both ast
   and vm (which imports ast)
*/


// Val represents every runtime value in the language

type Val interface{}

// atoms
type Tag int
type Sym int
type Bool bool
type Int int64
type Char rune

// fns
type Fun int

type Clos struct {
	Cp  int
	Env []Val
}





// compound

// repreentation needs updating in light of new function based
// key convention
type Object struct {
	Tag  Tag
	Vals []Val
}

// Built in Tags

const (
	Tag_tag Tag = iota
	Sym_tag
	Bool_tag
	Int_tag
	Char_tag
)

// Object Tags
type tagTable struct {
	table map[string]Tag
	uniq  uint
}

func initTagTable() *tagTable {
	return nil
}

func (ot *tagTable) Lookup() {}
func (ot *tagTable) AddTag() {}
