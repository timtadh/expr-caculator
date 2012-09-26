package ast

import "fmt"
import "strings"

type Node interface {
    Label() interface{}
    Kids() int
    GetKid(int) Node
    AddKid(Node) Node
    String() string
    Dotty() string
}

type AstNode struct {
    label string
    children []Node
}

func NewAstNode(label string) *AstNode {
    return &AstNode{
      label:label,
      children:make([]Node, 0),
    }
}

func New(label string) Node {
    return NewAstNode(label)
}

func (self *AstNode) Label() interface{} {
    return self.label
}

func (self *AstNode) Kids() int {
    return len(self.children)
}

func (self *AstNode) GetKid(i int) Node {
    return self.children[i]
}

func (self *AstNode) AddKid(kid Node) Node {
    self.children = append(self.children, kid)
    return self
}

func (self *AstNode) String() string {
    return String(self)
}

func (self *AstNode) Dotty() string {
    return Dotty(self)
}

func String(self Node) string {
    s := fmt.Sprintf("%v:%v\n", self.Kids(), self.Label())
    for j := 0; j < self.Kids(); j++ {
        s += self.GetKid(j).String()
    }
    return s
}

func Dotty(self Node) string {
    internal := "%v [shape=rect, label=\"%v\"];"
    leaf := "%v [shape=rect, label=\"%v\" style=\"filled\" fillcolor=\"#dddddd\"];"
    edge := "%v -> %v;"
    nodes := make([]string, 0)
    edges := make([]string, 0)

    type elem struct {
        node Node
        i int
    }
    i := 0
    stack := make([]*elem, 0)
    stack = append(stack, &elem{self, i})
    i += 1
    for len(stack) > 0 {
        n := stack[0].node
        c := stack[0].i
        stack = stack[1:]
        name := fmt.Sprintf("n%v", c)
        if n.Kids() > 0 {
            nodes = append(nodes, fmt.Sprintf(internal, name, n.Label()))
        } else {
            nodes = append(nodes, fmt.Sprintf(leaf, name, n.Label()))
        }
        for j := 0; j < n.Kids(); j++ {
            kid := n.GetKid(j)
            edges = append(edges, fmt.Sprintf(edge, name, fmt.Sprintf("n%v", i)))
            stack = append(stack, &elem{kid, i})
            i += 1
        }
    }
    header := "digraph G {\n"
    footer := "\n}\n"
    return header + strings.Join(nodes, "\n") + strings.Join(edges, "\n") + footer
}

