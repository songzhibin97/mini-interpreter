package parser

import (
	"fmt"
	"strconv"

	"github.com/songzhibin97/mini-interpreter/ast"
	"github.com/songzhibin97/mini-interpreter/lexer"
	"github.com/songzhibin97/mini-interpreter/token"
)

type prefixParserFunc func() ast.Expr
type infixParserFunc func(left ast.Expr) ast.Expr

type Parser struct {
	l         *lexer.Lexer
	curToken  *token.Token
	peekToken *token.Token
	errors    []string

	prefixParseHandler map[token.Type]prefixParserFunc
	infixParseHandler  map[token.Type]infixParserFunc
}

func (p *Parser) registerPrefix(t token.Type, fn prefixParserFunc) {
	if p.prefixParseHandler == nil {
		p.prefixParseHandler = make(map[token.Type]prefixParserFunc)
	}
	p.prefixParseHandler[t] = fn
}

func (p *Parser) registerInfix(t token.Type, fn infixParserFunc) {
	if p.infixParseHandler == nil {
		p.infixParseHandler = make(map[token.Type]infixParserFunc)
	}
	p.infixParseHandler[t] = fn
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{
		Stmts: []ast.Stmt{},
	}
	for ; p.curToken.Type != token.EOF; p.nextToken() {
		stmt := p.parseStmt()
		if stmt == nil {
			continue
		}
		program.Stmts = append(program.Stmts, stmt)
	}
	return program
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) assertionCurToken(t token.Type) bool {
	return p.curToken.Type == t
}

func (p *Parser) assertionPeekToken(t token.Type) bool {
	return p.peekToken.Type == t
}

func (p *Parser) assertionPeekTokenErr(t token.Type) {
	p.errors = append(p.errors, fmt.Sprintf("expected token %s, got %s", t, p.peekToken.Type))
}

func (p *Parser) forecastNextPeek(t token.Type) bool {
	if p.assertionPeekToken(t) {
		p.nextToken()
		return true
	}
	p.assertionPeekTokenErr(t)
	return false
}

// ============================================================================

func (p *Parser) parseExpr(precedence int) ast.Expr {
	prefix := p.prefixParseHandler[p.curToken.Type]
	if prefix == nil {
		p.errors = append(p.errors, fmt.Sprintf("no prefix parse function for %s found", p.curToken.Type))
		return nil
	}
	leftExpr := prefix()

	for precedence < p.peekToken.Type.Precedence() {
		infix := p.infixParseHandler[p.peekToken.Type]
		if infix == nil {
			return leftExpr
		}
		p.nextToken()
		leftExpr = infix(leftExpr)
	}
	return leftExpr
}

func (p *Parser) parseIdentifierExpr() ast.Expr {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parseIntegerExpr() ast.Expr {
	v, err := strconv.ParseInt(p.curToken.Value, 0, 64)
	if err != nil {
		p.errors = append(p.errors, fmt.Sprintf("could not parse %s as integer", p.curToken.Value))
		return nil
	}
	return &ast.Integer{Token: p.curToken, Value: v}
}

func (p *Parser) parseStringExpr() ast.Expr {
	return &ast.String{Token: p.curToken, Value: p.curToken.Value}
}

func (p *Parser) parsePrefixExpr() ast.Expr {
	expr := &ast.PrefixExpr{
		Token:    p.curToken,
		Operator: p.curToken.Value,
	}
	p.nextToken()

	expr.Right = p.parseExpr(token.UnaryPrec)
	return expr
}

func (p *Parser) parseBooleanExpr() ast.Expr {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.assertionCurToken(token.TRUE),
	}
}

func (p *Parser) parseGroupedExpr() ast.Expr {
	p.nextToken()

	expr := p.parseExpr(token.LowestPrec)

	if !p.forecastNextPeek(token.RPAREN) {
		return nil
	}
	return expr
}

func (p *Parser) parseIfExpr() ast.Expr {
	expr := &ast.IfExpr{Token: p.curToken}
	if !p.forecastNextPeek(token.LPAREN) {
		return nil
	}

	p.nextToken()
	expr.Condition = p.parseExpr(token.LowestPrec)

	if !p.forecastNextPeek(token.RPAREN) {
		return nil
	}

	if !p.forecastNextPeek(token.LBRACE) {
		return nil
	}

	expr.Consequence = p.parseBlockStmt()

	if p.assertionPeekToken(token.ELSE) {
		p.nextToken()

		if !p.forecastNextPeek(token.LBRACE) {
			return nil
		}
		expr.Alternative = p.parseBlockStmt()
	}
	return expr
}

func (p *Parser) parseFuncExpr() ast.Expr {
	f := &ast.FuncExpr{Token: p.curToken}

	if !p.forecastNextPeek(token.IDENT) {
		return nil
	}
	f.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Value,
	}

	if !p.forecastNextPeek(token.LPAREN) {
		return nil
	}

	f.Params = p.parseFuncParams()

	if !p.forecastNextPeek(token.LBRACE) {
		return nil
	}

	f.Body = p.parseBlockStmt()

	return f
}

func (p *Parser) parseArrayExpr() ast.Expr {
	return &ast.Array{Token: p.curToken, Elements: p.parseElements(token.RBRACK)}
}

func (p *Parser) parseMapExpr() ast.Expr {
	mp := &ast.Map{Token: p.curToken, Elements: make(map[ast.Expr]ast.Expr)}

	for !p.assertionPeekToken(token.RBRACE) {
		p.nextToken()
		key := p.parseExpr(token.LowestPrec)

		if !p.forecastNextPeek(token.COLON) {
			return nil
		}
		p.nextToken()

		value := p.parseExpr(token.LowestPrec)
		mp.Elements[key] = value

		if !p.assertionPeekToken(token.RBRACE) && !p.forecastNextPeek(token.COMMA) {
			return nil
		}
	}

	if !p.forecastNextPeek(token.RBRACE) {
		return nil
	}
	return mp
}

func (p *Parser) parseMacroExpr() ast.Expr {
	expr := &ast.Macro{Token: p.curToken}

	if !p.forecastNextPeek(token.IDENT) {
		return nil
	}

	expr.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Value,
	}

	if !p.forecastNextPeek(token.LPAREN) {
		return nil
	}

	expr.Params = p.parseFuncParams()

	if !p.forecastNextPeek(token.LBRACE) {
		return nil
	}

	expr.Body = p.parseBlockStmt()

	return expr
}

func (p *Parser) parseInfixExpr(left ast.Expr) ast.Expr {
	expr := &ast.InfixExpr{
		Token:    p.curToken,
		Operator: p.curToken.Value,
		Left:     left,
	}
	precedence := p.curToken.Type.Precedence()
	p.nextToken()
	expr.Right = p.parseExpr(precedence)
	return expr
}

func (p *Parser) parseCallExpr(left ast.Expr) ast.Expr {
	return &ast.CallExpr{Token: p.curToken, Func: left, Args: p.parseElements(token.RPAREN)}
}

func (p *Parser) parseIndexExpr(left ast.Expr) ast.Expr {
	expr := &ast.IndexExpr{Token: p.curToken, Left: left}
	p.nextToken()
	expr.Index = p.parseExpr(token.LowestPrec)

	if !p.forecastNextPeek(token.RBRACK) {
		return nil
	}
	return expr
}

func (p *Parser) parseFuncParams() []*ast.Identifier {
	var params []*ast.Identifier

	if p.assertionPeekToken(token.RPAREN) {
		p.nextToken()
		return params
	}
	p.nextToken()

	params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Value})

	for p.assertionPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		params = append(params, &ast.Identifier{Token: p.curToken, Value: p.curToken.Value})
	}

	if !p.forecastNextPeek(token.RPAREN) {
		return nil
	}
	return params
}

func (p *Parser) parseElements(end token.Type) []ast.Expr {
	var args []ast.Expr

	if p.assertionPeekToken(end) {
		p.nextToken()
		return args
	}

	p.nextToken()
	args = append(args, p.parseExpr(token.LowestPrec))

	for p.assertionPeekToken(token.COMMA) {
		p.nextToken()
		p.nextToken()
		args = append(args, p.parseExpr(token.LowestPrec))
	}

	if !p.forecastNextPeek(end) {
		return nil
	}

	return args
}

// ============================================================================

func (p *Parser) parseStmt() ast.Stmt {
	switch p.curToken.Type {
	case token.VAR:
		return p.parseVarStmt()
	case token.RETURN:
		return p.parseReturnStmt()
	default:
		return p.parseExprStmt()
	}
}

func (p *Parser) parseVarStmt() *ast.VarStmt {
	s := &ast.VarStmt{
		Token: p.curToken,
	}
	if !p.forecastNextPeek(token.IDENT) {
		return nil
	}

	s.Name = &ast.Identifier{
		Token: p.curToken,
		Value: p.curToken.Value,
	}
	if !p.forecastNextPeek(token.ASSIGN) {
		return nil
	}

	p.nextToken()

	s.Value = p.parseExpr(token.LowestPrec)

	return s
}

func (p *Parser) parseReturnStmt() *ast.ReturnStmt {
	s := &ast.ReturnStmt{
		Token: p.curToken,
	}
	p.nextToken()

	s.Value = p.parseExpr(token.LowestPrec)

	return s
}

func (p *Parser) parseExprStmt() *ast.ExprStmt {
	s := &ast.ExprStmt{
		Token: p.curToken,
		Expr:  p.parseExpr(token.LowestPrec),
	}

	return s
}

func (p *Parser) parseBlockStmt() *ast.BlockStmt {
	block := &ast.BlockStmt{Token: p.curToken}
	p.nextToken()
	for !p.assertionCurToken(token.RBRACE) && !p.assertionCurToken(token.EOF) {
		stmt := p.parseStmt()
		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		}
		p.nextToken()
	}
	return block
}

// ============================================================================

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}
	p.nextToken()
	p.nextToken()

	p.registerPrefix(token.IDENT, p.parseIdentifierExpr)
	p.registerPrefix(token.INT, p.parseIntegerExpr)
	p.registerPrefix(token.STRING, p.parseStringExpr)
	p.registerPrefix(token.SUB, p.parsePrefixExpr)
	p.registerPrefix(token.NOT, p.parsePrefixExpr)
	p.registerPrefix(token.TRUE, p.parseBooleanExpr)
	p.registerPrefix(token.FALSE, p.parseBooleanExpr)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpr)
	p.registerPrefix(token.IF, p.parseIfExpr)
	p.registerPrefix(token.FUNC, p.parseFuncExpr)
	p.registerPrefix(token.LBRACK, p.parseArrayExpr)
	p.registerPrefix(token.LBRACE, p.parseMapExpr)
	p.registerPrefix(token.MACRO, p.parseMacroExpr)

	p.registerInfix(token.ADD, p.parseInfixExpr)
	p.registerInfix(token.SUB, p.parseInfixExpr)
	p.registerInfix(token.QUO, p.parseInfixExpr)
	p.registerInfix(token.MUL, p.parseInfixExpr)
	p.registerInfix(token.EQL, p.parseInfixExpr)
	p.registerInfix(token.ASSIGN, p.parseInfixExpr)
	p.registerInfix(token.NEQ, p.parseInfixExpr)
	p.registerInfix(token.LSS, p.parseInfixExpr)
	p.registerInfix(token.GTR, p.parseInfixExpr)
	p.registerInfix(token.LPAREN, p.parseCallExpr)
	p.registerInfix(token.LBRACK, p.parseIndexExpr)

	return p
}
