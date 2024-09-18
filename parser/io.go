package parser

import (
	"github.com/yassinebenaid/nbs/ast"
	"github.com/yassinebenaid/nbs/token"
)

func (p *Parser) getIOParser(tt token.TokenType) func(*ast.Command) {
	switch tt {
	case token.GT:
		return p.parseStdoutRedirection
	default:
		return nil
	}
}

func (p *Parser) parseStdoutRedirection(cmd *ast.Command) {
	var r ast.Redirection
	r.Src = ast.FileDescriptor("1")
	r.Method = p.curr.Literal

	p.proceed()
	if p.curr.Type == token.BLANK {
		p.proceed()
	}

	r.Dst = p.parseSentence()

	cmd.Redirections = append(cmd.Redirections, r)
}