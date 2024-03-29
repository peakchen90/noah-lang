package lexer

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

type Fixture struct {
	Input  string `json:"input"`
	Output struct {
		Type     string `json:"type"`
		Value    string `json:"value"`
		Chars    []rune `json:"chars"`
		Flag     string `json:"flag"`
		Position [2]int `json:"position"`
	} `json:"output"`
}

func validateToken() {
	for _, item := range tokenMetaTable {
		if &item == nil {
			panic("TokenMetaTable")
		}
	}
}

func TestLexer(t *testing.T) {
	validateToken()

	fixturesContent, err := os.ReadFile("testdata/fixtures.json")
	if err != nil {
		panic(err)
	}

	fixtures := make([]Fixture, 0, 20)
	err = json.Unmarshal(fixturesContent, &fixtures)
	if err != nil {
		panic(err)
	}

	for _, fixture := range fixtures {
		token := NewLexer([]rune(fixture.Input)).Next()
		assert.Equal(t, token.Name, fixture.Output.Type, "Token name")
		if len(fixture.Output.Value) > 0 {
			assert.Equal(t, fixture.Output.Value, token.Value, "Token value")
		} else if len(fixture.Output.Chars) > 0 {
			assert.Equal(t, string(fixture.Output.Chars), token.Value, "Token value")
		}
		if len(fixture.Output.Flag) > 0 {
			assert.Equal(t, string(fixture.Output.Flag), token.Flag, "Token Flag")
		}
		assert.Equal(t, fixture.Output.Position[0], token.Start, "Token position start")
		assert.Equal(t, fixture.Output.Position[1], token.End, "Token position end")
	}

	// 关键字
	for _, item := range keywords {
		token := NewLexer([]rune(item)).Next()
		if item == "is" {
			assert.Equal(t, TTIsOp, token.Type, "Token type")
		} else if item == "as" {
			assert.Equal(t, TTAsOp, token.Type, "Token type")
		} else {
			assert.Equal(t, TTKeyword, token.Type, "Token type")
		}
		assert.Equal(t, item, token.Value, "Token value")
		assert.Equal(t, 0, token.Start, "Token position start")
		assert.Equal(t, len(item), token.End, "Token position end")
	}

	// 内置常量
	for _, item := range builtInConstants {
		token := NewLexer([]rune(item)).Next()
		assert.Equal(t, TTConst, token.Type, "Token type")
		assert.Equal(t, item, token.Value, "Token value")
		assert.Equal(t, 0, token.Start, "Token position start")
		assert.Equal(t, len(item), token.End, "Token position end")
	}

	// 其他字符 token
	for i, item := range tokenMetaTable {
		if len(item.Text) == 0 {
			continue
		}
		lexer := NewLexer([]rune(item.Text))

		switch item.Type {
		case TTPrefixPlus, TTPrefixSub, TTPrefixInc, TTPrefixDec:
			lexer.allowExpr = true
		default:
			lexer.allowExpr = false
		}

		token := lexer.Next()
		assert.Equal(t, uint(i), uint(token.Type), "Table index")
		assert.Equal(t, item.Type, token.Type, "Token type")
		assert.Equal(t, 0, token.Start, "Token position start")
		assert.Equal(t, len(item.Text), token.End, "Token position end")
	}
}
