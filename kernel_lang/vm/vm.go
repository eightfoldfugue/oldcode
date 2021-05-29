package vm

import (
	"lang/obj"
	"lang/code"
	"log"
)


func Eval(c code.Code) obj.Val {
	vm := new(machine)
	vm.instr = c.Ops
	vm.globals = make([]obj.Val, c.SizeGlobals)
	vm.run()
	return vm.pop()
}

const stk_size int = 327668

type mark code.Op  // encodes the last ep and bottom of env
type ret  code.Op  // encodes the return point

type machine struct {
	sp, ip int
	ep int
	
	stk   [stk_size]obj.Val

	instr []code.Op

	globals []obj.Val
}

func (vm *machine) fetch() code.Op {
	i := vm.instr[vm.ip]
	vm.ip++
	return i
}

// could change this to have nice stack overflow errors
func (vm *machine) push(v obj.Val) {
	if vm.sp == len(vm.stk) {
		log.Fatalf("stack overflow, limit is %v", len(vm.stk))
	}
	vm.stk[vm.sp] = v
	vm.sp++
}

func (vm *machine) pop() obj.Val {
	vm.sp--
	v := vm.stk[vm.sp]
	vm.stk[vm.sp] = nil // clear for gc
	return v
}

func (vm *machine) apply(fn obj.Val, rp int) {
	vm.ep = vm.sp
	switch f := fn.(type) {
	case obj.Fun:
		vm.push(ret(rp))
		vm.ip = int(f)
	case obj.Clos:
		// push args in the correct order
		for i := len(f.Env)-1; i >= 0; i-- {
			vm.push(f.Env[i])
		}
		
		// then push ret point then go to ip
		vm.push(ret(rp))
		vm.ip = f.Cp
			
	default:
		log.Fatal("from apply blank is not a function line n")
	}
}

func (vm *machine) run() {
loop:
	for {
		switch vm.fetch() {
		case code.PushInt:
			num := vm.fetch()
			vm.push(num) // TODO annotate with type

		case code.Add, code.Sub, code.Mul:
			op := vm.instr[vm.ip-1]
			num1 := vm.pop().(code.Op)
			num2 := vm.pop().(code.Op)
			switch op {
			case code.Add:
				vm.push(num1 + num2)
			case code.Mul:
				vm.push(num1 * num2)
			case code.Sub:
				vm.push(num1 - num2)
			}

		case code.PushMark:
			ep_last := vm.ep
		    vm.push(mark(ep_last))

		case code.PushFun:
			ofs := vm.fetch()
			vm.push(obj.Fun(vm.ip))
			vm.ip += int(ofs)

		case code.Call:
			ret_point := vm.ip
			fn := vm.pop()
			vm.apply(fn, ret_point)
						
		case code.Grab:
			vm.ep--
			_, is_mark := vm.stk[vm.ep].(mark)
			if is_mark {
				ret_point := vm.pop().(ret)
				env_size := (vm.sp - 1) - vm.ep

				env := make([]obj.Val, env_size)

				for i, _ := range env {
					env[i] = vm.pop()
				}

				// ip-1 sets us back to the grab instr for next invocation
				clos := obj.Clos{vm.ip-1, env}

				mark := vm.pop().(mark)
				vm.ep = int(mark)
				
				vm.push(clos)
				
				vm.ip = int(ret_point)
				
			}
			
		case code.Access:
			at := int(vm.fetch())
			vm.push(vm.stk[vm.ep + at])
			
		case code.Return:
			ret_val := vm.pop()

			// either way, deactivate frame
			
			ret_point := vm.pop().(ret)
			env_size := vm.sp - vm.ep
			for i := 0; i < env_size; i++ {
				vm.pop()
			}
			// if we have a mark, we're done, else apply fn to results
			_, is_mark := vm.stk[vm.sp-1].(mark)
			if is_mark {
				mark := vm.pop().(mark)
				vm.ep = int(mark)
				vm.push(ret_val)
				vm.ip = int(ret_point)
 								
			} else {
				vm.apply(ret_val, int(ret_point))
			}

		case code.SetDef:
			index := vm.fetch()
			vm.globals[index] = vm.pop()
			
		case code.GetDef:
			index := vm.fetch()
			vm.push(vm.globals[index])
			
		case code.Halt:
			break loop
		}
	}
}

