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
		assert.Equal(t, fixture.Output.Position[0], token.Start, "Token position start")
		assert.Equal(t, fixture.Output.Position[1], token.End, "Token position end")
	}

	// 关键字
	for _, item := range Keywords {
		token := NewLexer([]rune(item)).Next()
		assert.Equal(t, TTKeyword, token.Type, "Token type")
		assert.Equal(t, item, token.Value, "Token value")
		assert.Equal(t, 0, token.Start, "Token position start")
		assert.Equal(t, len(item), token.End, "Token position end")
	}

	// 内置常量
	for _, item := range Constants {
		token := NewLexer([]rune(item)).Next()
		assert.Equal(t, TTConst, token.Type, "Token type")
		assert.Equal(t, item, token.Value, "Token value")
		assert.Equal(t, 0, token.Start, "Token position start")
		assert.Equal(t, len(item), token.End, "Token position end")
	}
}
