package parsers

import "github.com/timtadh/lex"
import "github.com/timtadh/expr-calculator/lexer"
import "github.com/timtadh/expr-calculator/ast"

type Parse func([]byte) (bool, ast.Node)
type Errors []Error
type Error string

func (self Error) String() string {
    return string(self)
}
func (self Error) Error() string {
    return string(self)
}

func Recursive(text []byte) (Errors, ast.Node) {
    success, tokens := lex.Lex(lexer.Patterns, text)

    type production func()(bool, ast.Node)
    var Expr production
    var Expr_ production
    var Term production
    var Term_ production
    var Factor production

    errors := make(Errors, 0)

    var _top lex.Token
    peak := func() (tok lex.Token, closed bool) {
        if _top == nil {
            if token, ok := <-tokens; !ok {
                return nil, true
            } else {
                _top = token
            }
        }
        return _top, false
    }

    consume := func() {
        _top = nil
    }

    Expr = func() (bool, ast.Node) {
        // Expr : Term Expr_
        n := ast.New(ast.Type("Expr"))
        ok1, r0 := Term()                              // Expr : Term . Expr_
        ok2, r1 := Expr_()                             // Expr : Term Expr_ .
        return ok1&&ok2, n.AddKid(r0).AddKid(r1)
    }

    Expr_ = func() (bool, ast.Node) {
        /* Expr_ : PLUS Term Expr_
         * Expr_ : DASH Term Expr_
         * Expr_ : e (the empty string) */
        n := ast.New(ast.Type("Expr_"))
        token, closed := peak()
        if closed {                                       // Expr_ : e .
            return true, n.AddKid(ast.New(ast.Type("e")))
        }
        a := token.ID()
        if a == lexer.Tokens["PLUS"] {                    // Expr_ : PLUS . Term Expr_
            consume()
            n.AddKid(ast.New(ast.Type("+")))
        } else if a == lexer.Tokens["DASH"] {             // Expr_ : DASH . Term Expr_
            consume()
            n.AddKid(ast.New(ast.Type("-")))
        } else {
            return true, n.AddKid(ast.New(ast.Type("e")))           // Expr_ : e .
        }
        ok1, r0 := Term()                              // Expr_ : (PLUS|DASH) Term . Expr_
        ok2, r1 := Expr_()                             // Expr_ : (PLUS|DASH) Term Expr_ .
        return ok1&&ok2, n.AddKid(r0).AddKid(r1)
    }

    Term = func() (bool, ast.Node) {
        // Term : Factor Term_
        n := ast.New(ast.Type("Term"))
        ok1, r0 := Factor()                            // Term : Factor . Term_
        ok2, r1 := Term_()                             // Term : Facttor Term_ .
        return ok1&&ok2, n.AddKid(r0).AddKid(r1)
    }

    Term_ = func() (bool, ast.Node) {
        /* Term_ : STAR Factor Term_
         * Term_ : SLASH Factor Term_
         * Term_ : e (the empty string) */
        n := ast.New(ast.Type("Term_"))
        token, closed := peak()
        if closed {                             // Term_ : e .
            return true, n.AddKid(ast.New(ast.Type("e")))
        }
        a := token.ID()
        if a == lexer.Tokens["STAR"] {                    // Term_ : STAR . Factor Term_
            consume()
            n.AddKid(ast.New(ast.Type("*")))
        } else if a == lexer.Tokens["SLASH"] {            // Term_ : SLASH . Factor Term_
            consume()
            n.AddKid(ast.New(ast.Type("/")))
        } else {
            return true, n.AddKid(ast.New(ast.Type("e")))        // Term_ : e .
        }
        ok1, r0 := Factor()                            // Term_ : (STAR|SLASH) Factor . Term_
        ok2, r1 := Term_()                             // Term_ : (STAR|SLASH) Factor Term_ .
        return ok1&&ok2, n.AddKid(r0).AddKid(r1)
    }

    Factor = func() (bool, ast.Node) {
        /* Factor : NUMBER
         * Factor : DASH NUMBER
         * Factor : LPAREN Expr RPAREN */
        n := ast.New(ast.Type("Factor"))
        token, closed := peak()
        if closed {
            errors = append(errors, "Expected a (NUMBER|DASH|LPAREN) got EOF")
            return false, ast.New(ast.Type("end of input"))
        }
        a := token.ID()
        if a == lexer.Tokens["NUMBER"] {                  // Factor : NUMBER .
            consume()
            label := token.Attribute().Value().(ast.EquaCompare)
            n.AddKid(ast.New(label))
        } else if a == lexer.Tokens["DASH"] {             // Factor : DASH . NUMBER
            consume()
            token, closed := peak()
            if closed {
                errors = append(errors, "expected a NUMBER got EOF")
                return false, ast.New(ast.Type("end of input"))
            } else if token.ID() == lexer.Tokens["NUMBER"] { // Factor : DASH NUMBER .
                consume()
                label := token.Attribute().Value().(ast.EquaCompare)
                n.AddKid(ast.New(ast.Type("-")))
                n.AddKid(ast.New(label))
            } else {
                errors = append(errors, "Expected a NUMBER")
                return false, n
            }
        } else if a == lexer.Tokens["LPAREN"] {      // Factor : LPAREN . Expr RPAREN
            consume()
            ok, r0 := Expr()                         // Factor : LPAREN Expr . RPAREN
            if !ok {
                return false, n
            }
            token, closed := peak()
            if closed {
                errors = append(errors, "Expected an RPAREN found EOF")
                return false, n
            } else if token.ID() == lexer.Tokens["RPAREN"] { // Factor : LPAREN Expr RPAREN .
                consume()
                n.AddKid(ast.New(ast.Type("("))).AddKid(r0).AddKid(ast.New(ast.Type(")")))
            } else {
                errors = append(errors, "Expected an RPAREN")
                return false, n
            }
        } else {
            errors = append(errors, "Expected a (NUMBER|DASH|LPAREN)")
            return false, n
        }
        return true, n
    }

    ok, root := Expr()
    if _, closed := peak(); !closed {
        errors = append(errors, "unconsumed input")
        for _ = range tokens {}
        ok = false
    }
    if !(<-success) {
        errors = append(errors, "lexing error (unconsumed input in the lexer)")
        ok = false
    }
    if !ok {
        return errors, nil
    }

    return nil, root
}

func LL1_to_AST(root ast.Node) ast.Node {

    var Expr func(ast.Node) ast.Node
    var Expr_ func(ast.Node) (bool, ast.Node, ast.Node)
    var Term func(ast.Node) ast.Node
    var Term_ func(ast.Node) (bool, ast.Node, ast.Node)
    var Factor func(ast.Node) ast.Node

    Expr = func(node ast.Node) ast.Node {
        left := Term(node.GetKid(0))
        if empty, op, right := Expr_(node.GetKid(1)); empty {
            return left
        } else {
            return op.AddKid(left).AddKid(right)
        }
        panic("unreachable")
    }

    Expr_ = func(node ast.Node) (bool, ast.Node, ast.Node) {
        if node.GetKid(0).Label().Equals(ast.Type("e")) {
            return true, nil, nil
        }
        myop := ast.New(node.GetKid(0).Label())
        left := Term(node.GetKid(1))
        if empty, op, right := Expr_(node.GetKid(2)); empty {
            return false, myop, left
        } else {
            return false, myop, op.AddKid(left).AddKid(right)
        }
        panic("unreachable")
    }

    Term = func(node ast.Node) ast.Node {
        left := Factor(node.GetKid(0))
        if empty, op, right := Term_(node.GetKid(1)); empty {
            return left
        } else {
            return op.AddKid(left).AddKid(right)
        }
        panic("unreachable")
    }

    Term_ = func(node ast.Node) (bool, ast.Node, ast.Node) {
        if node.GetKid(0).Label().Equals(ast.Type("e")) {
            return true, nil, nil
        }
        myop := ast.New(node.GetKid(0).Label())
        left := Factor(node.GetKid(1))
        if empty, op, right := Term_(node.GetKid(2)); empty {
            return false, myop, left
        } else {
            return false, myop, op.AddKid(left).AddKid(right)
        }
        panic("unreachable")
    }

    Factor = func(node ast.Node) ast.Node {
        if node.GetKid(0).Label().Equals(ast.Type("(")) {
            return Expr(node.GetKid(1))
        } else if node.GetKid(0).Label().Equals(ast.Type("-")) {
            return ast.New(lexer.Int(int(node.GetKid(1).Label().(lexer.Int)) * -1))
        }
        return ast.New(node.GetKid(0).Label())
    }

    return Expr(root)
}

