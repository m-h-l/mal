package parser

import (
	"fmt"
	"mal/reader"
	"mal/types"
	"strings"
)

func readForm(r *reader.Reader) (types.MalType, bool) {
	t, _ := r.Peek()
	switch t {
	case "(":
		return readList("(", ")", types.List, r)
	case "[":
		return readList("[", "]", types.Vector, r)
	case "{":
		return readList("{", "}", types.Map, r)
	case "'":
		r.Read()
		f, ok := readForm(r)
		if ok {
			return types.NewMalList(types.List, []types.MalType{types.NewMalGenericAtom("quote"), f}), true
		} else {
			fmt.Println("EOF")
			return nil, false
		}
	case "`":
		r.Read()
		f, ok := readForm(r)
		if ok {
			return types.NewMalList(types.List, []types.MalType{types.NewMalGenericAtom("quasiquote"), f}), true
		} else {
			fmt.Println("EOF")
			return nil, false
		}
	case "~":
		r.Read()
		f, ok := readForm(r)
		if ok {
			return types.NewMalList(types.List, []types.MalType{types.NewMalGenericAtom("unquote"), f}), true
		} else {
			fmt.Println("EOF")
			return nil, false
		}
	case "@":
		r.Read()
		f, ok := readForm(r)
		if ok {
			return types.NewMalList(types.List, []types.MalType{types.NewMalGenericAtom("deref"), f}), true
		} else {
			fmt.Println("EOF")
			return nil, false
		}
	case "~@":
		r.Read()
		f, ok := readForm(r)
		if ok {
			return types.NewMalList(types.List, []types.MalType{types.NewMalGenericAtom("splice-unquote"), f}), true
		} else {
			fmt.Println("EOF")
			return nil, false
		}
	case "^":
		r.Read()
		meta, ok1 := readForm(r)
		obj, ok2 := readForm(r)
		if ok1 && ok2 {
			return types.NewMalList(types.List, []types.MalType{
				types.NewMalGenericAtom("with-meta"),
				obj,
				meta,
			}), true
		} else {
			fmt.Println("EOF in metadata")
			return nil, false
		}
	default:
		return readAtom(r)
	}
}

func readList(start string, end string, kind types.MalTypeId, r *reader.Reader) (*types.MalList, bool) {
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
			r.Read()
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

func readAtom(r *reader.Reader) (types.MalAtom, bool) {
	t, ok := r.Read()
	if !ok {
		return nil, false
	}

	if strings.HasPrefix(t, "\"") && strings.HasSuffix(t, "\"") {
		val := t[1 : len(t)-1]
		return types.NewMalString(val), true
	}

	if strings.HasPrefix(t, "\"") {
		fmt.Println("EOF")
		return nil, false
	}

	return types.NewMalGenericAtom(t), true
}

func Parse(r *reader.Reader) (types.MalType, bool) {
	return readForm(r)
}
