package lexer

import "fmt"
import "strconv"
import "github.com/timtadh/lex"
import "github.com/timtadh/expr-calculator/ast"

var rtokens []string = []string{
    "NUMBER", "PLUS", "DASH", "STAR", "SLASH", "LPAREN", "RPAREN" }
var Tokens map[string]int

func init() {
    Tokens = make(map[string]int)
    for i, token := range rtokens {
        Tokens[token] = i
    }
}

type Token struct {
    id int
    attribute *Attr
}

type Attr struct {
    value ast.EquaCompare
    has_value bool
}

type Int int

func (self Int) Equals(other ast.EquaCompare) bool {
    if o, ok := other.(Int); ok {
        return int(self) == int(o)
    }
    return false
}

func (self *Attr) String() string {
    if self == nil {
        return "<nil>"
    } else if self.has_value {
        return fmt.Sprintf("<Attribute %v>", self.value)
    }
    return "<Attribute>"
}

func (self *Attr) Value() interface{} {
    return self.value
}

func (self *Attr) HasValue() bool {
    return self.has_value
}

func (self *Token) Name() string {
    return rtokens[self.id]
}

func (self *Token) ID() int {
    return self.id
}

func (self *Token) Attribute() lex.Attribute {
    return self.attribute
}

var Patterns []*lex.Pattern = []*lex.Pattern{
  &lex.Pattern{
    "[0-9]+",
    func(match []byte)(bool, lex.Token) {
        i, err := strconv.ParseInt(string(match), 10, 64); 
        if err != nil {
            panic(err)
        }
        return true, &Token{Tokens["NUMBER"], &Attr{Int(i), true}}
    },
  },
  &lex.Pattern{
    "\\+",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{Tokens["PLUS"], &Attr{Int(0), false}}
    },
  },
  &lex.Pattern{
    "-",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{Tokens["DASH"], &Attr{Int(0), false}}
    },
  },
  &lex.Pattern{
    "\\*",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{Tokens["STAR"], &Attr{Int(0), false}}
    },
  },
  &lex.Pattern{
    "/",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{Tokens["SLASH"], &Attr{Int(0), false}}
    },
  },
  &lex.Pattern{
    "\\(",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{Tokens["LPAREN"], &Attr{Int(0), false}}
    },
  },
  &lex.Pattern{
    "\\)",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{Tokens["RPAREN"], &Attr{Int(0), false}}
    },
  },
  &lex.Pattern{
    "(\t|\n| )+",
    func(match []byte)(bool, lex.Token) {
        return false, nil
    },
  },
}

