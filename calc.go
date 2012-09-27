package main

import "fmt"
import "os"
import "io/ioutil"

import "github.com/timtadh/expr-calculator/parsers"

func main() {
    if expr, err := ioutil.ReadAll(os.Stdin); err != nil {
        panic("could not read the stdin")
    } else {
        if errors, ast := parsers.Recursive(expr); errors == nil {
            fmt.Fprintln(os.Stderr, parsers.LL1_to_AST(ast))
            fmt.Println(ast.Dotty())
        } else {
            fmt.Fprintln(os.Stderr, "parsing failed")
            for _, err := range errors {
                fmt.Fprintln(os.Stderr, ">", err)
            }
        }
    }
}

