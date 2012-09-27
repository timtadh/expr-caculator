package il

import "fmt"

type Inst interface {
    Op() Op
    A() Argument
    B() Argument
    Result() Argument
    String() string
}

type Op interface {
    ID() int
    Name() string
    String() string
}

type Argument interface {
    String() string
}

type inst struct {
    op Op
    a Argument
    b Argument
    result Argument
}

type None bool
type Constant int
type Symbol string

type op int

var rops []string = []string{
    "ADD", "SUB", "MUL", "DIV", "IMM",
}

var ops map[string]op

func init() {
    ops = make(map[string]op)
    for i, _op := range rops {
        ops[_op] = op(i)
    }
}

func MakeNone() Argument {
    return None(false)
}

func MakeConst(i int) Argument {
    return Constant(i)
}

func MakeSymbol(name string) Argument {
    return Symbol(name)
}

func (self None) String() string {
    return ""
}

func (self Constant) String() string {
    return fmt.Sprintf("<const:%v>", int(self))
}

func (self Symbol) String() string {
    return fmt.Sprintf("<sym:%v>", string(self))
}

func GetOp(name string) Op {
    if val, has := ops[name]; has {
        return val
    }
    panic(fmt.Sprintf("Unknown op %v", name))
}

func (self op) ID() int {
    return int(self)
}

func (self op) Name() string {
    return rops[self]
}

func (self op) String() string {
    return rops[self]
}

func NewInst(op Op, a, b, result Argument) Inst {
    return &inst{op, a, b, result}
}

func (self *inst) Op() Op {
    return self.op
}

func (self *inst) A() Argument {
    return self.a
}

func (self *inst) B() Argument {
    return self.b
}

func (self *inst) Result() Argument {
    return self.result
}

func (self *inst) String() string {
    return fmt.Sprintf("%-5v %-10v %-10v -> %-10v", self.op, self.a, self.b, self.result)
}

