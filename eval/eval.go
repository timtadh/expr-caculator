package eval

import "fmt"

import "github.com/timtadh/expr-calculator/lexer"
import "github.com/timtadh/expr-calculator/ast"
import "github.com/timtadh/expr-calculator/il"

func Eval(root ast.Node) int {

    var eval func(ast.Node) int
    eval = func(node ast.Node) int {
        if node.Kids() == 0 {
            return int(node.Label().(lexer.Int))
        }
        op := string(node.Label().(ast.Type))
        return map[string]func(int, int) int {
            "+":func(a, b int) int {
                return a + b
            },
            "-":func(a, b int) int {
                return a - b
            },
            "*":func(a, b int) int {
                return a * b
            },
            "/":func(a, b int) int {
                return a / b
            },
        }[op](eval(node.GetKid(0)), eval(node.GetKid(1)))
    }

    return eval(root)
}

func Interpret(program []il.Inst) int {

    symbols := make(map[string]int)

    resolve := func(A il.Argument) int {
        var a int
        switch v := A.(type) {
        case il.Symbol:
            a = symbols[string(v)]
        case il.Constant:
            a = int(v)
        case il.None:
            // do nothing
        default:
            panic(fmt.Sprintf("unknown type %v", v))
        }
        return a
    }

    for _, inst := range program {
        a := resolve(inst.A())
        b := resolve(inst.B())
        result := string(inst.Result().(il.Symbol))
        switch(inst.Op().ID()) {
        case il.GetOp("IMM").ID():
            symbols[result] = a
        case il.GetOp("ADD").ID():
            symbols[result] = a + b
        case il.GetOp("SUB").ID():
            symbols[result] = a - b
        case il.GetOp("MUL").ID():
            symbols[result] = a * b
        case il.GetOp("DIV").ID():
            symbols[result] = a / b
        default:
            panic(fmt.Sprintf("Unexpected op %v", inst.Op().Name()))
        }
    }

    return symbols[string(program[len(program)-1].Result().(il.Symbol))]
}

