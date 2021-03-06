package main

import "fmt"
import "os"
import "io/ioutil"

import "github.com/timtadh/expr-calculator/parsers"
import "github.com/timtadh/expr-calculator/eval"
import "github.com/timtadh/expr-calculator/il"

func main() {
    if expr, err := ioutil.ReadAll(os.Stdin); err != nil {
        panic("could not read the stdin")
    } else {
        if errors, parse_tree := parsers.Recursive(expr); errors == nil {
            ast := parsers.LL1_to_AST(parse_tree)
            fmt.Fprintln(os.Stderr, ast)
            fmt.Fprintln(os.Stderr, eval.Eval(ast))
            insts := il.IL(ast)
            for i, inst := range insts {
                fmt.Fprintln(os.Stderr, i, inst)
            }
            fmt.Fprintln(os.Stderr, eval.Interpret(insts))
            fmt.Println(parse_tree.Dotty())
        } else {
            fmt.Fprintln(os.Stderr, "parsing failed")
            for _, err := range errors {
                fmt.Fprintln(os.Stderr, ">", err)
            }
        }
    }
}

