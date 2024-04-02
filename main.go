package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
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
)

type Estado int

const (
	Inicio Estado = iota
	IdentificadorEst
	EnteroEst
	DecimalEst
	MayorMenorEst
	DiferenciaEst
	IgualEst
	OpMasMenosEst
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
	conteoLinea int
}

type Lexema struct {
	input  string
	pos    int
	estado Estado
	Tokens []Token
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
		input:  content,
		pos:    0,
		estado: Inicio,
		Tokens: []Token{},
		contLinea: 1,
	}
}

func TipoToken(tipo string) TokenType {
	// Aquí puedes agregar más lógica para clasificar el tipo de token
	for _, palRes := range palabrasRes{
		if tipo == palRes{
			return TokenPalabrasReservadas
		}
	}

	if _, err := strconv.Atoi(tipo); err == nil {
		return TokenNumEntero
	}

	if _, err := strconv.ParseFloat(tipo, 64); err == nil{
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

	}
	// Aquí podrías agregar más casos para otros tipos de tokens, como operadores, etc.
	return TokenIdentificador
}

func (lex *Lexema) AnalisisLex() []Token {
	var tokens []Token
	var lexAux string = " "
	entrada := lex.input
	lex.estado = Inicio

	//Lectura de la cadena realizada del archivo de texto
	for i := 0; i < len(entrada); i++ {
		char := rune(entrada[i]) //rune() -> convierte los valores tipo byte o caracter en un rune para ser usado en funciones unicode
		fmt.Println("Vuelta: ", i)
		

		//Estado en los que se encuntra la expresion
		switch lex.estado {
		//Estado inicial de la expresion
		case Inicio:
			fmt.Println("Entre al inicio")
			fmt.Println(lexAux + "cam")
			if unicode.IsSpace(char) {
				fmt.Println("Hola .,::")
				if(char == '\n'){
					lex.contLinea++
				}
				continue

			} else if unicode.IsLetter(char) || char == '_' {
				lex.estado = IdentificadorEst
				lexAux += string(char)
				fmt.Print("pr:" + lexAux + "\n")
			} else if unicode.IsDigit(char) {
				lex.estado = EnteroEst
				lexAux += string(char)
				fmt.Println("xd: " + lexAux)
			} else if char == '(' || char == ';' || char == ',' || char == ')' || char == '{' || char == '}'{
				fmt.Println("Entre a: " + string(char))
				lex.estado = Inicio
				lexAux += string(char)
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				lexAux = ""
			} else if char == '<' || char == '>' {
				fmt.Println("Entre a " + string(char))
				lex.estado = MayorMenorEst
				lexAux += string(char)
			} else if char == '!' {
				fmt.Println("Entre a " + string(char))
				lex.estado = DiferenciaEst
				lexAux += string(char)
				fmt.Println(lexAux + "::")
			} else if char == '=' {
				lex.estado = IgualEst
				lexAux += string(char)
			} else if char == '+' || char == '-' {
				lex.estado = OpMasMenosEst
				lexAux += string(char)
			} else if char == '*' || char == '%' || char == '^' {
				lex.estado = Inicio
				lexAux += string(char)
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				lexAux = ""
			} else if char == '/' {
				lex.estado = OpDivEst
				lexAux += string(char)
			} else {
				if i == len(lex.input) && char == '\n' {
					fmt.Println("Termino el analisis")
				} else {
					fmt.Println("el caracter " + string(char) + " no es valido")
					lex.estado = Inicio
				}
			}
			println("mf:" + lexAux)

			fmt.Println(tokens)
			fmt.Println("Sali del estado de inicio")
		//Estado en la que la expresion inico como letra o '_'
		case IdentificadorEst:
			fmt.Println("Entre estado id")
			if unicode.IsLetter(char) || unicode.IsDigit(char) || char == '_' {
				lex.estado = IdentificadorEst
				lexAux += string(char)
				fmt.Println("k: " + lexAux)
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				lexAux = ""
				i--
			}
		//Estado en la que la expresion inicio como numero
		case EnteroEst:
			fmt.Println("Entre al estado de entero")
			if unicode.IsDigit(char) {
				lex.estado = EnteroEst
				lexAux += string(char)
				fmt.Println("po: " + lexAux)

			} else if char == '.' {
				lex.estado = DecimalEst
				lexAux += string(char)

			} else {
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				fmt.Println(tokens)
				lexAux = ""
				lex.estado = Inicio
				i--
			}
		case DecimalEst:
			fmt.Println("Entre al estado de decimal")
			if unicode.IsDigit(char) {
				lex.estado = DecimalEst
				lexAux += string(char)
				fmt.Println("po2: " + lexAux)
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				fmt.Println(tokens)
				lexAux = ""
				i--
			}

			fmt.Println("Sali del decimal")
		case MayorMenorEst:
			if char == '=' {
				lex.estado = MayorMenorEst
				lexAux += string(char)
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				lexAux = ""
				i--
			}
		case DiferenciaEst:
			fmt.Println("Entre al estado de diferencia")
			if char == '=' {
				lex.estado = DiferenciaEst
				lexAux += string(char)
			} else {
				println("El caracter " + string(char) + " No puede ir solo")
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				fmt.Println(tokens)
				fmt.Println(lexAux + "ññ")
				lexAux = ""
				i--

			}
		case IgualEst:
			if char == '=' {
				lex.estado = IgualEst
				lexAux += string(char)
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				lexAux = ""
				i--
			}
		case OpMasMenosEst:
			if unicode.IsDigit(char) {
				lex.estado = EnteroEst
				lexAux += string(char)
			}else if char == '+'{
				lex.estado = OpMasMenosEst
			}else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				lexAux = ""
				i--
			}
		case OpDivEst:
			fmt.Println("Entre al op /")
			if char == '*' {
				lex.estado = ComentMultiIniEst
				lexAux += string(char)
			} else if char == '/' {
				lex.estado = ComentUniEst
				lexAux += string(char)
			} else {
				lex.estado = Inicio
				tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
				lexAux = ""
				i--
			}
			
		case ComentMultiIniEst:
			fmt.Println("Entre al coment uni")
			if unicode.IsDigit(char) || unicode.IsLetter(char) || unicode.IsSpace(char) || unicode.IsSymbol(char) {
				lex.estado = ComentMultiIniEst
				lexAux += string(char)
			} else if char == '*' {
				lex.estado = ComentMultiFinEst
				lexAux += string(char)
			}
		case ComentMultiFinEst:
			if char == '/' {
				lex.estado = Inicio
				fmt.Println("hola lol")
				lexAux = ""
				i--
			}
		case ComentUniEst:
			fmt.Println("wow::" + lexAux)
			if unicode.IsDigit(char) || unicode.IsLetter(char) || char == ' ' || unicode.IsSymbol(char) {
				lex.estado = ComentUniEst
				lexAux += string(char)
			} else {
				lex.estado = Inicio
				fmt.Println("wow" + lexAux)
				lexAux = ""
				i--
			}
			
		}

	}

	if lexAux != "" {
		tokens = append(tokens, Token{TipoToken(lexAux), lexAux, lex.contLinea})
		lexAux = ""
	}

	fmt.Print(tokens)

	return tokens

}

func main() {
	file, err := os.Open("texto.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	tokens := newLexema(file)

	for _, token := range tokens.AnalisisLex() {
		
		fmt.Printf("Type: %v, Value: %v, numero de linea: %v\n", token.Type, token.Valor, token.conteoLinea)

	}

}
