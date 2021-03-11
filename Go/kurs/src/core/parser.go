package core

import (
  "fmt"
  "strings"
)

const (
  NORMAL = iota
  SYMBOL = iota
  NAME = iota
  TEXT = iota
  TEXT2 = iota
  ESCAPE = iota
  NUMBER = iota
  EOF = iota
)

type Token struct {
  Typ         int
  StringValue string
  CharValue   int32
}

var nextTokenId int

func ResetTokenId() {
  nextTokenId = 0
}

func NextToken(tokens []Token) Token {
  if nextTokenId >= len(tokens) {
    return EOFTOKEN
  }
  t := tokens[nextTokenId]
  nextTokenId++
  return t
}

func IsSymbol(t Token, c int32) bool {
  return t.Typ == SYMBOL && t.CharValue == c
}

func IsName(t Token, s string) bool {
  return t.Typ == NAME && t.StringValue == s
}

func IsText(t Token, s string) bool {
  return t.Typ == TEXT && t.StringValue == s
}

func (t Token) String() string {
  return fmt.Sprintf("%v %v %v", t.Typ, t.StringValue, string(t.CharValue))
}

var EOFTOKEN = Token{EOF, "", 0}

func Parse(data []byte) []Token {
  mode := NORMAL
  var text strings.Builder
  tokens := make([]Token, 100000)
  for _, b := range string(data) {
    repeat := true
    for repeat {
      switch mode {
      case NORMAL:
        if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= 'а' && b <= 'я') || (b >= 'А' && b <= 'Я') {
          mode = NAME
          text.WriteRune(b)
        } else if (b >= '0' && b <= '9') || b == '-' {
          mode = NUMBER
          text.WriteRune(b)
        } else if b == '"' {
          mode = TEXT
        } else if b == '\'' {
          mode = TEXT2
        } else if b != ' ' && b != '\t' && b != '\r' && b != '\n' {
          //fmt.Printf("SYMBOL %v\n", string(b))
          tokens = append(tokens, Token{SYMBOL, "", b})
          text.Reset()
        }
        repeat = false
      case NUMBER:
        if (b >= '0' && b <= '9') || b == '.' {
          text.WriteRune(b)
          repeat = false
        } else {
          tokens = append(tokens, Token{NUMBER, text.String(), 0})
          text.Reset()
          mode = NORMAL
        }
      case NAME:
        if (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z')  || (b >= 'а' && b <= 'я') || (b >= 'А' && b <= 'Я') ||
          (b >= '0' && b <= '9') || b == '_' || b == '-' {
          text.WriteRune(b)
          repeat = false
        } else {
          //fmt.Printf("NAME %v\n", text.String())
          tokens = append(tokens, Token{NAME, text.String(), 0})
          text.Reset()
          mode = NORMAL
        }
      case TEXT:
        if b == '\\' {
          mode = ESCAPE
        } else if b != '"' {
          text.WriteRune(b)
        } else {
          //fmt.Printf("TEXT %v\n", text.String())
          tokens = append(tokens, Token{TEXT, text.String(), 0})
          text.Reset()
          mode = NORMAL
        }
        repeat = false
      case TEXT2:
        if b != '\'' {
          text.WriteRune(b)
        } else {
          //fmt.Printf("TEXT2 %v\n", text.String())
          tokens = append(tokens, Token{TEXT2, text.String(), 0})
          text.Reset()
          mode = NORMAL
        }
        repeat = false
      case ESCAPE:
        switch b {
        case 'n':
          text.WriteByte('\n')
        case 'r':
          text.WriteByte('\r')
        case 't':
          text.WriteByte('\t')
        default:
          text.WriteRune(b)
        }
        mode = TEXT
        repeat = false
      }
    }
  }
  return tokens
}
