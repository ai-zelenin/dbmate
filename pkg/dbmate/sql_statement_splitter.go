package dbmate

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type IgnoreToken struct {
	Flag       bool
	BeginIndex int
	EndIndex   int
}

type SQLStatementSplitter struct {
	autoIgnorePatterns []*regexp.Regexp
	autoSplitPattern   *regexp.Regexp
	manualSplitPattern *regexp.Regexp
}

func NewSQLStatementSplitter() *SQLStatementSplitter {
	return &SQLStatementSplitter{
		autoIgnorePatterns: []*regexp.Regexp{
			regexp.MustCompile(`--.*\n`),
			regexp.MustCompile(`(?sU:'.*')`),
			regexp.MustCompile(`(?sU:\$[a-zA-Z_]*\$.*\$[a-zA-Z_]*\$)`),
		},
		autoSplitPattern:   regexp.MustCompile(`;`),
		manualSplitPattern: regexp.MustCompile(`--\s?-{2,}`),
	}
}

func (s *SQLStatementSplitter) SplitAuto(text string) []string {
	var ignoreMatrix = make([][]IgnoreToken, 0)
	t := time.Now()
	for _, pattern := range s.autoIgnorePatterns {
		matches := pattern.FindAllStringIndex(text, -1)
		ignores := make([]IgnoreToken, 0)
		for _, match := range matches {
			ignore := IgnoreToken{
				BeginIndex: match[0],
				EndIndex:   match[1],
			}
			ignores = append(ignores, ignore)
		}
		ignoreMatrix = append(ignoreMatrix, ignores)
		fmt.Printf("%s %s \n", pattern, time.Now().Sub(t))
	}

	splitIndexes := make([]int, 0)
	probableSplitIndex := s.autoSplitPattern.FindAllStringIndex(text, -1)
	for _, bn := range probableSplitIndex {
		i := bn[0]
		skip := false
		for _, ignores := range ignoreMatrix {
			if skip {
				break
			}
			for _, ignore := range ignores {
				if i >= ignore.BeginIndex && i < ignore.EndIndex {
					skip = true
					break
				}
			}
		}
		if !skip {
			splitIndexes = append(splitIndexes, i+1)
		}
	}
	//i := 0
	//for {
	//	if i >= len(text) {
	//		break
	//	}
	//	r := text[i]
	//	skip := false
	//	for _, ignores := range ignoreMatrix {
	//		if skip {
	//			break
	//		}
	//		for _, ignore := range ignores {
	//			if i >= ignore.BeginIndex && i < ignore.EndIndex {
	//				skip = true
	//				break
	//			}
	//		}
	//
	//	}
	//	if r == ';' && !skip {
	//		splitIndexes = append(splitIndexes, i+1)
	//	}
	//	i++
	//}
	fmt.Printf("run %s \n", time.Now().Sub(t))
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
