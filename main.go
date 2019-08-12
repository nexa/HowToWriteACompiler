package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"
)

type Token struct {
	kind  string // "intliteral", "punct"
	value string
}

var source []byte
var sourceIndex int = 0

func getChar() (byte, error) {
	if sourceIndex == len(source) {
		return 0, errors.New("EOF")
	}
	char := source[sourceIndex]
	sourceIndex++
	return char, nil
}

func ungetChar() {
	sourceIndex--
}

func tokenize() []*Token {
	var tokens []*Token
	fmt.Printf("# Tokens : ")

	for {
		char, err := getChar()
		if err != nil {
			break
		}
		var token *Token
		switch char {
		case ' ', '\t', '\n':
			continue
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			var number []byte = []byte{char}
			for {
				char, err := getChar()
				if err != nil {
					break
				}
				if '0' <= char && char <= '9' {
					number = append(number, char)
				} else {
					ungetChar()
					break
				}
			}
			token = &Token{
				kind:  "intliteral",
				value: string(number),
			}
		case ';':
			token = &Token{
				kind: "punct",
				value: string([]byte{char}),
			}
		default:
			panic(fmt.Sprintf("tokenizer: Invalid char: '%c'", char))
		}

		fmt.Printf(" '%s'", token.value)
		tokens = append(tokens, token)
	}

	fmt.Printf("\n")
	return tokens
}

var tokens []*Token
var tokenIndex int = 0

type Expr struct {
	kind   string // "intliteral"
	intval int    // for intliteral
}

func getToken() *Token {
	if tokenIndex >= len(tokens) {
		return nil
	}
	token := tokens[tokenIndex]
	tokenIndex++
	return token
}

func parseUnaryExpr() *Expr {
	token := getToken()

	intval, _ := strconv.Atoi(token.value)
	expr := &Expr{
		kind: "intliteral",
		intval: intval,
	}
	return expr
}

func parse() *Expr {
	expr := parseUnaryExpr()
	return expr
}

func generateExpr(expr *Expr) {
	fmt.Printf("  movq $%d, %%rax\n", expr.intval)
}

func generateCode(expr *Expr) {
	fmt.Printf("  .global main\n")
	fmt.Printf("main:\n")
	generateExpr(expr)
	fmt.Printf("  ret\n")
}

func main() {
	source, _ = ioutil.ReadFile("/dev/stdin")
	tokens = tokenize()
	expr := parse()
	generateCode(expr)
}
