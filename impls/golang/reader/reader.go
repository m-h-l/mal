package reader

import (
	"fmt"
	"mal/types"
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

func GetAST(r *Reader) (types.MalType, bool) {
	return readForm(r)
}

func readForm(r *Reader) (types.MalType, bool) {
	t, _ := r.Peek()
	switch t {
	case "(":
		return readList("(", ")", types.List, r)
	case "[":
		return readList("[", "]", types.Vector, r)
	case "{":
		return readList("{", "}", types.Map, r)
	default:
		return readAtom(r)
	}
}

func readList(start string, end string, kind types.MalTypeId, r *Reader) (*types.MalList, bool) {
	if first, ok := r.Read(); first != start || !ok {
		panic(fmt.Sprint("readList: Expected '", start, "'"))
	}
	if next, _ := r.Peek(); next == end {
		r.Read()
		return types.NewMalList(kind, []types.MalType{}), true
	}

	items := []types.MalType{}

	for {
		t, ok := r.Peek()
		if t == end {
			break
		}
		if !ok {
			fmt.Println("EOF")
			return nil, false
		}
		form, ok := readForm(r)
		if ok {
			items = append(items, form)
		} else {
			return nil, false
		}

	}

	return types.NewMalList(kind, items), true
}

func readAtom(r *Reader) (*types.MalAtom, bool) {
	t, ok := r.Read()
	if !ok {
		return nil, false
	}
	return types.NewMalAtom(t), true
}
