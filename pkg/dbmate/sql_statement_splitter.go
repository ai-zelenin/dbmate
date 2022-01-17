package dbmate

import (
	"regexp"
	"sort"
	"strings"
)

type IgnoreToken struct {
	BeginSeq   string
	EndSeq     string
	Pattern    *regexp.Regexp
	Flag       bool
	BeginIndex int
	EndIndex   int
}

type SQLStatementSplitter struct {
	autoIgnorePatterns []*regexp.Regexp
	manualSplitPattern *regexp.Regexp
}

func NewSQLStatementSplitter() *SQLStatementSplitter {
	return &SQLStatementSplitter{
		autoIgnorePatterns: []*regexp.Regexp{
			regexp.MustCompile(`--.*\n|'.*'`),
			regexp.MustCompile(`(?sU:\$[a-zA-Z_]*\$.*\$[a-zA-Z_]*\$)`),
		},
		manualSplitPattern: regexp.MustCompile(`--\s?-{2,}`),
	}
}

func (s *SQLStatementSplitter) SplitAuto(text string) []string {
	var ignores = make([]IgnoreToken, 0)
	for _, pattern := range s.autoIgnorePatterns {
		matches := pattern.FindAllStringIndex(text, -1)
		for _, match := range matches {
			ignore := IgnoreToken{
				BeginIndex: match[0],
				EndIndex:   match[1],
			}
			ignores = append(ignores, ignore)
		}
	}
	sort.Slice(ignores, func(i, j int) bool {
		return ignores[i].BeginIndex < ignores[j].BeginIndex
	})

	splitIndexes := make([]int, 0)

	i := 0
	for {
		if i >= len(text) {
			break
		}
		r := text[i]
		skip := false
		for _, ignore := range ignores {
			if i >= ignore.BeginIndex && i < ignore.EndIndex {
				skip = true
				break
			}
		}
		if r == ';' && !skip {
			splitIndexes = append(splitIndexes, i+1)
		}
		i++
	}

	return s.splitByIndexes(text, splitIndexes)
}

func (s *SQLStatementSplitter) SplitManual(text string) []string {
	return s.manualSplitPattern.Split(text, -1)
}

func (s *SQLStatementSplitter) splitByIndexes(text string, splitIndexes []int) []string {
	statements := make([]string, 0)
	for i, index := range splitIndexes {
		var statement string
		var from int
		var to int
		if i > 0 {
			from = splitIndexes[i-1]
		}
		to = index
		statement = text[from:to]
		statement = strings.TrimFunc(statement, func(r rune) bool {
			return r == ' ' || r == '\n'
		})
		statements = append(statements, statement)
	}
	return statements
}
