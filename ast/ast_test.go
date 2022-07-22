package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/songzhibin97/mini-interpreter/token"
)

func TestProgram_String(t *testing.T) {
	p := Program{Stmts: []Stmt{
		VarStmt{
			Token: &token.Token{
				Type:  token.VAR,
				Value: "var",
			},
			Name: &Identifier{
				Token: &token.Token{
					Type:  token.IDENT,
					Value: "test",
				},
				Value: "test",
			},
			Value: &Identifier{
				Token: &token.Token{
					Type:  token.IDENT,
					Value: "value",
				},
				Value: "value",
			},
		},
	}}
	assert.Equal(t, p.String(), "var test = value")
}
