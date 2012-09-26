package main

import "fmt"
import "os"
import "io/ioutil"

import "github.com/timtadh/lex"
import "github.com/timtadh/expr-calculator/lexer"
import "github.com/timtadh/expr-calculator/parsers"

func main() {
    fmt.Fprintln(os.Stderr, "hello", lexer.HELLO)
    if expr, err := ioutil.ReadAll(os.Stdin); err != nil {
        panic("could not read the stdin")
    } else {
        success, tokens := lex.Lex(lexer.Patterns, expr)
        for token := range tokens {
            fmt.Fprintln(os.Stderr, token.Name(), token.Attribute())
        }
        fmt.Fprintln(os.Stderr, <-success)
        if ok, ast := parsers.Recursive(expr); ok {
            fmt.Fprintln(os.Stderr, ast)
            fmt.Println(ast.Dotty())
        } else {
            fmt.Fprintln(os.Stderr, "parsing failed")
        }
    }
}

