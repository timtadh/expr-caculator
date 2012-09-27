package il

import "fmt"

import "github.com/timtadh/expr-calculator/lexer"
import "github.com/timtadh/expr-calculator/ast"

func IL(root ast.Node) []Inst {
    sym_count := 0
    sym := func() Argument {
        name := fmt.Sprintf("t%v", sym_count)
        sym_count += 1
        return MakeSymbol(name)
    }

    var gen func(ast.Node) ([]Inst, Argument)
    gen = func(node ast.Node) ([]Inst, Argument) {
        if node.Kids() == 0 {
            c := MakeConst(int(node.Label().(lexer.Int)))
            i := NewInst(GetOp("IMM"), c, MakeNone(), sym())
            return []Inst{i}, i.Result()
        }
        left, l_val := gen(node.GetKid(0))
        right, r_val := gen(node.GetKid(1))
        insts := make([]Inst, len(left) + len(right), len(left) + len(right) + 1)
        copy(insts[:len(left)], left)
        copy(insts[len(left):], right)
        op := map[string]Op {
            "+":GetOp("ADD"),
            "-":GetOp("SUB"),
            "*":GetOp("MUL"),
            "/":GetOp("DIV"),
        }[string(node.Label().(ast.Type))]
        insts = append(insts, NewInst(op, l_val, r_val, sym()))
        return insts, insts[len(insts)-1].Result()
    }

    insts, _ := gen(root)
    return insts
}

