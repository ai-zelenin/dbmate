package dbmate

import "strings"

type SQLStatementSplitter struct {
	ignoreSeq   map[string]string
	ignoreFlags map[string]bool
}

func NewSQLStatementSplitter() *SQLStatementSplitter {
	return &SQLStatementSplitter{
		ignoreSeq: map[string]string{
			"'":      "'",
			"--":     "\n",
			"$body$": "$body$",
		},
		ignoreFlags: map[string]bool{},
	}
}

func (s *SQLStatementSplitter) Split(text string) []string {
	runes := []rune(text)
	splitIndexes := make([]int, 0)
	statements := make([]string, 0)

	for i, r := range runes {
		for begin, end := range s.ignoreSeq {
			if s.lookForward(runes, i, begin) {
				if !s.ignoreFlags[begin+end] {
					s.ignoreFlags[begin+end] = true
					continue
				}
			}
			if s.lookForward(runes, i, end) {
				if s.ignoreFlags[begin+end] {
					s.ignoreFlags[begin+end] = false
					continue
				}
			}
		}
		if r == ';' && !s.IsSkip() {
			splitIndexes = append(splitIndexes, i+1)
		}
	}
	for i, index := range splitIndexes {
		var statement string
		var from int
		var to int
		if i > 0 {
			from = splitIndexes[i-1]
		}
		to = index
		statement = string(runes[from:to])
		statement = strings.TrimFunc(statement, func(r rune) bool {
			return r == ' ' || r == '\n'
		})
		statements = append(statements, statement)
	}
	return statements
}

func (s *SQLStatementSplitter) IsSkip() bool {
	for _, b := range s.ignoreFlags {
		if b {
			return true
		}
	}
	return false
}

func (s *SQLStatementSplitter) lookForward(src []rune, index int, seq string) bool {
	if seq == "" {
		if index == len(seq)-1 {
			return true
		} else {
			return false
		}
	}
	rseq := []rune(seq)

	r := src[index]
	if r == rseq[0] {
		if len(rseq) == 1 {
			return true
		}
		for i, desiredRune := range rseq {
			next := index + i
			if next > len(src)-1 {
				return false
			}
			if src[next] != desiredRune {
				return false
			}
		}
		return true
	}
	return false
}
