package main

import (
	"fmt"
	"os"

	sitter "github.com/tree-sitter/go-tree-sitter"
	tree_sitter_java "github.com/tree-sitter/tree-sitter-java/bindings/go"
)

type MatchGroup string

const (
	VARIABLE_NAME MatchGroup = "variable_name"
	FUNCTION_DECLARATION MatchGroup = "function_declaration"
	FUNCTION_CALL MatchGroup = "function_call"
)

var matchMap = map[uint]MatchGroup{
	0: FUNCTION_DECLARATION,
	1: VARIABLE_NAME,
	2: FUNCTION_CALL,
}

func main() {
	code, err := os.ReadFile("./SpelExpression.java")
	if err != nil {
		fmt.Println("Error in reading file, error:", err)
	}

	parser := sitter.NewParser()
	defer parser.Close()
	lang := sitter.NewLanguage(tree_sitter_java.Language())
	parser.SetLanguage(lang)

	tree := parser.Parse(code, nil)
	defer tree.Close()

	queryStr := fmt.Sprintf(`
    (method_declaration
      name: (identifier) @%s)
    (variable_declarator
      name: (identifier) @%s)
    (method_invocation
    name: (identifier) @%s)
`, FUNCTION_DECLARATION, VARIABLE_NAME, FUNCTION_CALL)
	query, _ := sitter.NewQuery(lang, queryStr)
	defer query.Close()

	qc := sitter.NewQueryCursor()
	defer qc.Close()

	captures := qc.Captures(query, tree.RootNode(), code)
	resultMap := map[MatchGroup][]string{}
	for match, index := captures.Next(); match != nil; match, index = captures.Next() {
		patternIndex := match.PatternIndex
		matchMapVal, _ := matchMap[patternIndex]
		_, ok := resultMap[matchMapVal]
		if !ok {
			resultMap[matchMapVal] = []string{}
		}
		resultMap[matchMapVal] = append(resultMap[matchMapVal], match.Captures[index].Node.Utf8Text(code))
	}
	for k, v := range resultMap {
		fmt.Printf("%s: %v\n", k, v)
	}
}
