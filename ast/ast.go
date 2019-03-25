package ast

type Program struct {
	Statements []Statement
}

func (i Identifier) String() string {
	return i.Value
}

type Scope struct {
	Outer   *Scope
	Objects map[string]Reference
}

func NewScope(outer *Scope) *Scope {
	return &Scope{
		Outer:   outer,
		Objects: make(map[string]Reference),
	}
}

func (s Scope) Lookup(ident string) Reference {
	return s.Objects[ident]
}
