package code

import (
	"lang/ast"
	"log"
)

func Compile(prog ast.Ast) graph {
	g := graph{}
	for !isNil(prog) {
		def, ok := prog.(ast.Def)
		if !ok {
			msg := "non defining form at top level"
			log.Fatal(msg)
		}

		_, alreadyDefined := g[def.Name] 
		if alreadyDefined {
			msg := "%v already declared in top level"
			log.Fatalf(msg, def.Name)
		}
		
		g[def.Name] = comp_expr(def.Expr, newblock(), *new(env))
		prog = def.Next
	}
	return g
}

func isNil(a ast.Ast) bool {
	_, ok := a.(ast.Nil)
	return ok
}



func comp_expr(t ast.Ast, blk *block, env env) *block {
	switch tree := t.(type) {
	case ast.Int64:
		num := Op(tree.Val)
		blk.emit(PushInt)
		blk.immeadiate(num)	

	case ast.Iden:
		index, is_local := env.lookup(tree.Val)
		inline, is_prim := primitives[tree.Val]
		
		switch {
		case is_local:
			blk.emit(Access)
			blk.immeadiate(Op(index))
		case is_prim:
			blk = inline(blk)
		default:
			// default is the globals case
			blk.emit(GetDef)
			loc := blk.label()
			blk.addRef(tree.Val, loc)
		}
		
	case ast.Lambda:
		blk.emit(PushFun)
		offset := blk.label()
		
		blk = comp_lambdas(tree, blk, env)

		blk.emit(Return)
		fun_len := blk.last() - offset
		blk.backpatch(offset, fun_len)		
		
	case ast.Apply:
		more_optimal := check_bin_prim(tree.Rator)
		if more_optimal != nil {
			e1 := tree.Rand
			e2 := tree.Rator.(ast.Apply).Rand
			blk = comp_expr(e1, blk, env)
			blk = comp_expr(e2, blk, env)
			blk = more_optimal(blk)
		} else {
			blk.emit(PushMark)
			blk = comp_applies(tree, blk, env)
			blk.emit(Call)
		}
	}
	
	return blk
}



func comp_applies(t ast.Ast, blk *block, env env) *block {
	switch tree := t.(type) {
	case ast.Apply:
		blk = comp_expr(tree.Rand, blk, env)
		blk = comp_applies(tree.Rator, blk, env)
	default:
		blk = comp_expr(t, blk, env)
	}
	return blk
}

// this function yearns for pattern matching syntax
func check_bin_prim(t ast.Ast) func(*block) *block {
	a, ok := t.(ast.Apply)
	if ok {
		i, ok := a.Rator.(ast.Iden)
		if ok {
			f, ok := strict_binary_primitives[i.Val]
			if ok {
				return f
			}
		}
	}
	return nil
}

func comp_lambdas(t ast.Ast, blk *block, env env) *block {
	switch tree := t.(type) {
	case ast.Lambda:
		blk.emit(Grab)
		blk = comp_lambdas(tree.Body, blk, env.extend(tree.Var))
	default:
	
		blk = comp_expr(t, blk, env)
	}
	return blk
}

// mapping from string names to de bruijn style numberings
type env []string

func (e env) extend(s string) env {
	return append(e, s)
}

func (e env) lookup(s string) (int, bool) {
	top := 0
	for i := len(e)-1; i >= 0; i-- {
		if e[i] == s {
			return top, true
		}
		top +=1
	}
	return -1, false
}

