package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type TokenType int

const (
	TokenEOF TokenType = iota
	TokenNumEntero
	TokenNumDecimal
	TokenIdentificador
	TokenPalabrasReservadas
	TokenParentesisDer
	TokenParentesisIzq
	TokenLlaveIzq
	TokenLlaveDer
	TokenPuntoComa
	TokenComa
	TokenSuma
	TokenResta
	TokenMultiplicacion
	TokenDivision
	TokenResiduo
	TokenPotencia
	TokenAsignacion
	TokenOpCompMayor
	TokenOpCompMayorIgual
	TokenOpCompMenor
	TokenOpCompMenorIgual
	TokenOpIgual
	TokenNegacion
	TokenOpDiferencia
	TokenOpAsigUnitMas
	TokenOpAsigUnitMenos
	TokenError = 404
)

type Estado int

const (
	Inicio Estado = iota
	IdentificadorEst
	EnteroEst
	DecimalEst
	DecimalEst2
	MayorMenorEst
	DiferenciaEst
	IgualEst
	OpMasEst
	OpMenosEst
	OpAsingUnitMasEst
	OpAsingUnitMenosEst
	OpMultiExpResEst
	OpDivEst
	ComentMultiIniEst
	ComentMultiFinEst
	ComentUniEst
)

var palabrasRes = []string{"if", "else", "while", "for", "and", "or", "int", "float", "string", "while", "switch", "cin", "cout"}

type Token struct {
	Type  TokenType
	Valor string
}

type Errores struct {
	Valor     string
	contColum int
	contLine  int
}

type Lexema struct {
	input     string
	contColum int
	estado    Estado
	Tokens    []Token
	contLinea int
}

func newLexema(input io.Reader) *Lexema {
	scanner := bufio.NewScanner(input)
	var content string
	for scanner.Scan() {
		content += scanner.Text() + "\n"
	}

	fmt.Println(content)

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	return &Lexema{
		input:     content,
		contColum: 0,
		estado:    Inicio,
		Tokens:    []Token{},
		contLinea: 1,
	}
}

func TipoToken(tipo string) TokenType {
	// Aquí puedes agregar más lógica para clasificar el tipo de token
	for _, palRes := range palabrasRes {
		if tipo == palRes {
			return TokenPalabrasReservadas
		}
	}

	if _, err := strconv.Atoi(tipo); err == nil {
		return TokenNumEntero
	}

	if _, err := strconv.ParseFloat(tipo, 64); err == nil {
		return TokenNumDecimal
	}

	switch tipo {
	case "+":
		return TokenSuma
	case "-":
		return TokenResta
	case "*":
		return TokenMultiplicacion
	case "/":
		return TokenDivision
	case "%":
		return TokenResiduo
	case "^":
		return TokenPotencia
	case "(":
		return TokenParentesisIzq
	case ")":
		return TokenParentesisDer
	case "{":
		return TokenLlaveIzq
	case "}":
		return TokenLlaveDer
	case ";":
		return TokenPuntoComa
	case ",":
		return TokenComa
	case "<":
		return TokenOpCompMenor
	case "<=":
		return TokenOpCompMenorIgual
	case ">":
		return TokenOpCompMayor
	case ">=":
		return TokenOpCompMayorIgual
	case "==":
		return TokenOpIgual
	case "=":
		return TokenAsignacion
	case "!":
		return TokenNegacion
	case "!=":
		return TokenOpDiferencia
	case "++":
		return TokenOpAsigUnitMas
	case "--":
		return TokenOpAsigUnitMenos
	}
	// Aquí podrías agregar más casos para otros tipos de tokens, como operadores, etc.
	return TokenIdentificador
}

func (lex *Lexema) AnalisisLex() ([]Token, []Errores) {
	var tokens []Token
	var errores []Errores
	var lexAux string = " "
	entrada := lex.input
	lex.estado = Inicio
	//Lectura de la cadena realizada del archivo de texto
	for i := 0; i < len(entrada); i++ {
		char := rune(entrada[i])

		lex.contColum++

		//Estado en los que se encuntra la expresion
		switch lex.estado {
		//Estado inicial de la expresion
		case Inicio:
			if unicode.IsSpace(char) {

				if char == '\n' {
					lex.contLinea++
					lex.contColum = 0
				}
				continue

			} else if unicode.IsLetter(char) || char == '_' {
				lex.estado = IdentificadorEst
				lexAux += string(char)
			} else if char == '(' || char == ';' || char == ',' || char == ')' || char == '{' || char == '}' {
				lex.estado = Inicio
				lexAux += string(char)
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
			} else if char == '<' || char == '>' {
				lex.estado = MayorMenorEst
				lexAux += string(char)
			} else if char == '!' {
				lex.estado = DiferenciaEst
				lexAux += string(char)
			} else if char == '=' {
				lex.estado = IgualEst
				lexAux += string(char)
			} else if unicode.IsDigit(char) {
				lex.estado = EnteroEst
				lexAux += string(char)
			} else if char == '+' {
				lex.estado = OpMasEst
				lexAux += string(char)
			} else if char == '-' {
				lex.estado = OpMenosEst
				lexAux += string(char)
			} else if char == '*' || char == '%' || char == '^' {
				lex.estado = Inicio
				lexAux += string(char)
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
			} else if char == '/' {
				lex.estado = OpDivEst
				lexAux += string(char)
			} else {
				if i == len(lex.input) && char == '\n' {
				} else {
					lex.estado = Inicio
				}
				lexAux = string(char)
				errores = append(errores, Errores{lexAux, lex.contColum, lex.contLinea})
				lexAux = ""
			}
		case IdentificadorEst:
			if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' {
				lex.estado = IdentificadorEst
				lexAux += string(char)
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			}
		//Estado en la que la expresion inicio como numero
		case EnteroEst:
			if unicode.IsDigit(char) {
				lex.estado = EnteroEst
				lexAux += string(char)
			} else if char == '.' {
				lex.estado = DecimalEst
				lexAux += string(char)
			} else {
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				lex.estado = Inicio
				i--
				lex.contColum--
			}
		case DecimalEst:
			if unicode.IsDigit(char) {
				lex.estado = DecimalEst2
				lexAux += string(char)
			} else {
				lex.estado = Inicio
				lexAux = strings.TrimRight(lexAux, ".")
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				errores = append(errores, Errores{".", lex.contColum, lex.contLinea})
				lexAux = ""
				i--
				lex.contColum--
			}
		case DecimalEst2:
			if unicode.IsDigit(char) {
				lex.estado = DecimalEst2
				lexAux += string(char)
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			}

		case MayorMenorEst:
			if char == '=' {
				lexAux += string(char)
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			}
		case DiferenciaEst:
			if char == '=' {
				lexAux += string(char)
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--

			}
		case IgualEst:
			if char == '=' {
				lexAux += string(char)
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			}
		case OpMasEst:
			if unicode.IsDigit(char) {
				lex.estado = EnteroEst
				lexAux += string(char)
			} else if char == '+' {
				lexAux += string(char)
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			}
		case OpMenosEst:
			if unicode.IsDigit(char) {
				lex.estado = EnteroEst
				lexAux += string(char)
			} else if char == '-' {
				lexAux += string(char)
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			}
		case OpDivEst:
			if char == '*' {
				lex.estado = ComentMultiIniEst
				lexAux = ""
			} else if char == '/' {
				lex.estado = ComentUniEst
				lexAux = ""
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
				lexAux = ""
				i--
				lex.contColum--
			}
			

		case ComentMultiIniEst:
			if unicode.IsDigit(char) || unicode.IsLetter(char) || unicode.IsSpace(char) || unicode.IsSymbol(char) {
				if char == '\n' {
					lex.contLinea++
				}
				lex.estado = ComentMultiIniEst
			} else if char == '*' {
				lex.estado = ComentMultiFinEst
			}
		case ComentMultiFinEst:
			if char == '*'{
				lex.estado = ComentMultiFinEst
			}else if char == '/' {
				lex.estado = Inicio
			}else{
				if char == '\n' {
					lex.contLinea++
				}
				lex.estado = ComentMultiIniEst
			}
		case ComentUniEst:
			if char == '\n'{
				lex.estado = Inicio
				lexAux = ""
				i--
				lex.contColum--
			} else{
				lex.estado = ComentUniEst
			}
		}
	}

	if lexAux != "" {
		tokens = append(tokens, Token{TipoToken(lexAux), lexAux})
		lexAux = ""
	}
	return tokens, errores

}

func main() {
	file, err := os.Open("texto.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lex := newLexema(file)

	tokens, errores := lex.AnalisisLex()

	for _, token := range tokens {
		fmt.Printf("Type: %v, Value: %v \n", token.Type, token.Valor)
	}

	// Iterar sobre los errores
	for _, err := range errores {
		fmt.Printf("Error: %v,  Linea: %v, Columna: %v, \n", err.Valor,  err.contLine, err.contColum)
	}

}
