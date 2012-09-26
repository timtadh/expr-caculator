package lexer

import "fmt"
import "strconv"
import "github.com/timtadh/lex"

const (
    HELLO = iota
)

var rtokens []string = []string{
    "NUMBER", "PLUS", "DASH", "STAR", "SLASH", "LPAREN", "RPAREN" }
var tokens map[string]int

func init() {
    tokens = make(map[string]int)
    for i, token := range rtokens {
        tokens[token] = i
    }
}

type Token struct {
    id int
    attribute *Attr
}

type Attr struct {
    value int
    has_value bool
}

func (self *Attr) String() string {
    if self == nil {
        return "<nil>"
    } else if self.has_value {
        return fmt.Sprintf("<Attribute %v>", self.value)
    }
    return "<Attribute>"
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
        return true, &Token{tokens["NUMBER"], &Attr{int(i), true}}
    },
  },
  &lex.Pattern{
    "\\+",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{tokens["PLUS"], &Attr{0, false}}
    },
  },
  &lex.Pattern{
    "-",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{tokens["DASH"], &Attr{0, false}}
    },
  },
  &lex.Pattern{
    "\\*",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{tokens["STAR"], &Attr{0, false}}
    },
  },
  &lex.Pattern{
    "/",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{tokens["SLASH"], &Attr{0, false}}
    },
  },
  &lex.Pattern{
    "\\(",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{tokens["LPAREN"], &Attr{0, false}}
    },
  },
  &lex.Pattern{
    "\\)",
    func(match []byte)(bool, lex.Token) {
        return true, &Token{tokens["RPAREN"], &Attr{0, false}}
    },
  },
  &lex.Pattern{
    "\t+",
    func(match []byte)(bool, lex.Token) {
        return false, nil
    },
  },
  &lex.Pattern{
    "\n+",
    func(match []byte)(bool, lex.Token) {
        return false, nil
    },
  },
  &lex.Pattern{
    " +",
    func(match []byte)(bool, lex.Token) {
        return false, nil
    },
  },
}

