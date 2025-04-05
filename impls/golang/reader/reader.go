package reader

import (
	"regexp"
)

func NewReader(tokens string) *Reader {
	return &Reader{
		token:    tokenize(tokens),
		position: 0,
	}
}

type Reader struct {
	token    []string
	position int
}

func (r *Reader) Read() (string, bool) {
	if r.position >= len(r.token) {
		return "", false
	}
	token := r.token[r.position]
	r.position++
	return token, true
}

func (r *Reader) Peek() (string, bool) {
	if r.position >= len(r.token) {
		return "", false
	}
	return r.token[r.position], true
}

func tokenize(input string) []string {
	re := regexp.MustCompile(`[\s,]*(~@|[\[\]{}()'` + "`" +
		`~^@]|"(?:\\.|[^\\"])*"?|;.*|[^\s\[\]{}('"` + "`" +
		`,;)]*)`)
	result := []string{}
	for _, group := range re.FindAllStringSubmatch(input, -1) {
		if group[1] == "" || group[1][0] == ';' {
			continue
		}
		result = append(result, group[1])
	}
	return result
}
