package parsers

import "github.com/timtadh/lex"
import "github.com/timtadh/expr-calculator/lexer"
import "github.com/timtadh/expr-calculator/ast"

type Parse func([]byte) (bool, ast.Node)

type production func(int)(int, ast.Node)

func Recursive(text []byte) (bool, ast.Node) {
    var root ast.Node
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

    Expr = func(i int) (int, ast.Node) {
        // Expr : Term Expr_
        n := ast.New("Expr")
        var r0 ast.Node
        var r1 ast.Node
        i, r0 = Term(i)
        i, r1 = Expr_(i)
        return i, n.AddKid(r0).AddKid(r1)
    }

    Expr_ = func(i int) (int, ast.Node) {
        /* Expr_ : PLUS Term Expr_
         * Expr_ : DASH Term Expr_
         * Expr_ : e (the empty string) */
        n := ast.New("Expr_")
        if i >= len(tokens) {                   // Expr_ : e .
            return i, n.AddKid(ast.New("e"))
        }
        a := tokens[i].ID()
        if a == lexer.Tokens["PLUS"] {          // Expr_ : PLUS . Term Expr_
            i += 1
            n.AddKid(ast.New("+"))
        } else if a == lexer.Tokens["DASH"] {   // Expr_ : DASH . Term Expr_
            i += 1
            n.AddKid(ast.New("-"))
        } else {
            return i, n.AddKid(ast.New("e"))    // Expr_ : e .
        }
        var r0 ast.Node
        var r1 ast.Node
        i, r0 = Term(i)                        // Expr_ : (PLUS|DASH) Term . Expr_
        i, r1 = Expr_(i)                       // Expr_ : (PLUS|DASH) Term Expr_ .
        return i, n.AddKid(r0).AddKid(r1)
    }

    Term = func(i int) (int, ast.Node) {
        // Term : Factor Term_
        n := ast.New("Term")
        var r0 ast.Node
        var r1 ast.Node
        i, r0 = Factor(i)
        i, r1 = Term_(i)
        return i, n.AddKid(r0).AddKid(r1)
    }

    Term_ = func(i int) (int, ast.Node) {
        /* Term_ : STAR Factor Term_
         * Term_ : SLASH Factor Term_
         * Term_ : e (the empty string) */
        n := ast.New("Term_")
        if i >= len(tokens) {                   // Term_ : e .
            return i, n.AddKid(ast.New("e"))
        }
        a := tokens[i].ID()
        if a == lexer.Tokens["STAR"] {          // Term_ : STAR . Factor Term_
            i += 1
            n.AddKid(ast.New("*"))
        } else if a == lexer.Tokens["SLASH"] {  // Term_ : SLASH . Factor Term_
            i += 1
            n.AddKid(ast.New("/"))
        } else {
            return i, n.AddKid(ast.New("e"))    // Term_ : e .
        }
        var r0 ast.Node
        var r1 ast.Node
        i, r0 = Factor(i)                      // Term_ : (STAR|SLASH) Factor . Term_
        i, r1 = Term_(i)                       // Term_ : (STAR|SLASH) Factor Term_ .
        return i, n.AddKid(r0).AddKid(r1)
    }

    Factor = func(i int) (int, ast.Node) {
        /* Factor : NUMBER
         * Factor : DASH NUMBER
         * Factor : LPAREN Expr RPAREN */
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
                panic("Expected a NUMBER")
            }
        } else if a == lexer.Tokens["LPAREN"] { // Factor : LPAREN . Expr RPAREN
            i += 1
            var r0 ast.Node
            i, r0 = Expr(i)                    // Factor : LPAREN Expr . RPAREN
            a := tokens[i].ID()
            if a == lexer.Tokens["RPAREN"] {    // Factor : LPAREN Expr RPAREN .
                i += 1
                n.AddKid(ast.New("(")).AddKid(r0).AddKid(ast.New(")"))
            } else {
                panic("Expected an RPAREN")
            }
        } else {
            panic("Expected a (NUMBER|DASH|LPAREN)")
        }
        return i, n
    }

    _, root = Expr(0)
    return true, root
}

