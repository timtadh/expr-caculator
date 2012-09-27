package eval

import "github.com/timtadh/expr-calculator/lexer"
import "github.com/timtadh/expr-calculator/ast"

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

