package lexer

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

var tokenFixtures = [...]string{
	`"abc"`,
	`"\\a"`,
	`"\1a"`,
	`"\7777"`,
	`"\8"`,
	`"\xff"`,
	`"\x00"`,
	`"""a
b"c"
"""`,
	`0`,
	`1.2`,
	`-12`,
	`-12.34`,
}

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
		assert.Equal(t, token.Name, fixture.Output.Type)
		if len(fixture.Output.Value) > 0 {
			assert.Equal(t, fixture.Output.Value, token.Value)
		} else if len(fixture.Output.Chars) > 0 {
			assert.Equal(t, string(fixture.Output.Chars), token.Value)
		}
		assert.Equal(t, fixture.Output.Position[0], token.Position.Start)
		assert.Equal(t, fixture.Output.Position[1], token.Position.End)
	}

	for _, fixture := range Keywords {
		token := NewLexer([]rune(fixture)).Next()
		assert.NotNil(t, token)
	}
}
