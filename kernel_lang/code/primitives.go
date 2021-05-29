package code

type primitive_map map[string] func(*block) *block

var primitives = primitive_map{
	"add": binary_general(Add),
	"sub": binary_general(Sub),
	"mul": binary_general(Mul),
}

// strict in that sense that they can only be used in 
// the context of an application supplying all of the arguments

var strict_binary_primitives = primitive_map{
	"add": binary_strict(Add),
	"sub": binary_strict(Sub),
	"mul": binary_strict(Mul),
}

func binary_general(operation Op) func(*block) *block {
	return func(b *block) *block {
		b.emit(PushFun)
		b.immeadiate(8)
		b.emit(Grab)
		b.emit(Grab)
		b.emit(Access)
		b.immeadiate(0)
		b.emit(Access)
		b.immeadiate(1)
		b.emit(operation)
		b.emit(Return)
		return b
	}
}

func binary_strict(operation Op) func(*block) *block {
	return func(b *block) *block {
		b.emit(operation)
		return b
	}
}
