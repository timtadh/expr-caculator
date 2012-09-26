package parsers

import "fmt"
import "os"

import "github.com/timtadh/lex"
import "github.com/timtadh/expr-calculator/lexer"
import "github.com/timtadh/expr-calculator/ast"

type Parse func([]byte) (bool, ast.Node)

type production func(int)(bool, int, ast.Node)

func Recursive(text []byte) (bool, ast.Node) {
    success, tokchan := lex.Lex(lexer.Patterns, text)
    tokens := make([]lex.Token, 0)
    for token := range tokchan {
        tokens = append(tokens, token)
    }
    if !(<-success) {
        return false, nil
    }

    var Expr production
    var Expr_ production
    var Term production
    var Term_ production
    var Factor production

    Expr = func(i int) (bool, int, ast.Node) {
        // Expr : Term Expr_
        n := ast.New("Expr")
        var ok1, ok2 bool
        var r0,r1 ast.Node
        ok1, i, r0 = Term(i)
        ok2, i, r1 = Expr_(i)
        return ok1&&ok2, i, n.AddKid(r0).AddKid(r1)
    }

    Expr_ = func(i int) (bool, int, ast.Node) {
        /* Expr_ : PLUS Term Expr_
         * Expr_ : DASH Term Expr_
         * Expr_ : e (the empty string) */
        n := ast.New("Expr_")
        if i >= len(tokens) {                   // Expr_ : e .
            return true, i, n.AddKid(ast.New("e"))
        }
        a := tokens[i].ID()
        if a == lexer.Tokens["PLUS"] {          // Expr_ : PLUS . Term Expr_
            i += 1
            n.AddKid(ast.New("+"))
        } else if a == lexer.Tokens["DASH"] {   // Expr_ : DASH . Term Expr_
            i += 1
            n.AddKid(ast.New("-"))
        } else {
            return true, i, n.AddKid(ast.New("e"))    // Expr_ : e .
        }
        var ok1, ok2 bool
        var r0,r1 ast.Node
        ok1, i, r0 = Term(i)                        // Expr_ : (PLUS|DASH) Term . Expr_
        ok2, i, r1 = Expr_(i)                       // Expr_ : (PLUS|DASH) Term Expr_ .
        return ok1&&ok2, i, n.AddKid(r0).AddKid(r1)
    }

    Term = func(i int) (bool, int, ast.Node) {
        // Term : Factor Term_
        n := ast.New("Term")
        var ok1, ok2 bool
        var r0,r1 ast.Node
        ok1, i, r0 = Factor(i)
        ok2, i, r1 = Term_(i)
        return ok1&&ok2, i, n.AddKid(r0).AddKid(r1)
    }

    Term_ = func(i int) (bool, int, ast.Node) {
        /* Term_ : STAR Factor Term_
         * Term_ : SLASH Factor Term_
         * Term_ : e (the empty string) */
        n := ast.New("Term_")
        if i >= len(tokens) {                   // Term_ : e .
            return true, i, n.AddKid(ast.New("e"))
        }
        a := tokens[i].ID()
        if a == lexer.Tokens["STAR"] {          // Term_ : STAR . Factor Term_
            i += 1
            n.AddKid(ast.New("*"))
        } else if a == lexer.Tokens["SLASH"] {  // Term_ : SLASH . Factor Term_
            i += 1
            n.AddKid(ast.New("/"))
        } else {
            return true, i, n.AddKid(ast.New("e"))    // Term_ : e .
        }
        var ok1, ok2 bool
        var r0,r1 ast.Node
        ok1, i, r0 = Factor(i)                      // Term_ : (STAR|SLASH) Factor . Term_
        ok2, i, r1 = Term_(i)                       // Term_ : (STAR|SLASH) Factor Term_ .
        return ok1&&ok2, i, n.AddKid(r0).AddKid(r1)
    }

    Factor = func(i int) (bool, int, ast.Node) {
        /* Factor : NUMBER
         * Factor : DASH NUMBER
         * Factor : LPAREN Expr RPAREN */
        if i >= len(tokens) {
            fmt.Fprintf(os.Stderr, "Expected a (NUMBER|DASH|LPAREN)")
            return false, i, ast.New("end of input")
        }
        n := ast.New("Factor")
        a := tokens[i].ID()
        if a == lexer.Tokens["NUMBER"] {        // Factor : NUMBER .
            label := tokens[i].Attribute().(*lexer.Attr).StringValue()
            i += 1
            n.AddKid(ast.New(label))
        } else if a == lexer.Tokens["DASH"] {   // Factor : DASH . NUMBER
            i += 1
            a := tokens[i].ID()
            if a == lexer.Tokens["NUMBER"] {    // Factor : DASH NUMBER .
                label := tokens[i].Attribute().(*lexer.Attr).StringValue()
                i += 1
                n.AddKid(ast.New("-"))
                n.AddKid(ast.New(label))
            } else {
                fmt.Fprintln(os.Stderr, "Expected a NUMBER")
                return false, i, n
            }
        } else if a == lexer.Tokens["LPAREN"] { // Factor : LPAREN . Expr RPAREN
            i += 1
            var ok bool
            var r0 ast.Node
            ok, i, r0 = Expr(i)                    // Factor : LPAREN Expr . RPAREN
            if !ok {
                return false, i, n
            }
            if i < len(tokens) &&
               tokens[i].ID() == lexer.Tokens["RPAREN"] {    // Factor : LPAREN Expr RPAREN .
                i += 1
                n.AddKid(ast.New("(")).AddKid(r0).AddKid(ast.New(")"))
            } else {
                fmt.Fprintln(os.Stderr, "Expected an RPAREN")
                return false, i, n
            }
        } else {
            fmt.Fprintln(os.Stderr, "Expected a (NUMBER|DASH|LPAREN)")
            return false, i, n
        }
        return true, i, n
    }

    ok, i, root := Expr(0)
    if !ok || i != len(tokens) {
        return false, nil
    }
    return true, root
}

