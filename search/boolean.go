package search

import (
	"bytes"
	"github.com/degemer/document-search-engine/index"
	"sort"
	"strings"
)

const (
	BOOL_OR  = "or"
	BOOL_AND = "and"
	BOOL_NOT = "not"
)

type BooleanSearch struct {
	StandardSearch
}

func (bs BooleanSearch) Search(request string) (docScores []index.DocScore) {
	docScores, _ = bs.booleanScores(booleanPrefix(request), 0)
	sort.Sort(index.ByScore(docScores))
	return
}

func (bs BooleanSearch) booleanScores(prefixed []string, pos int) ([]index.DocScore, int) {
	if pos == len(prefixed) {
		return []index.DocScore{}, pos
	}
	switch prefixed[pos] {
	case BOOL_AND:
		docScores1, pos := bs.booleanScores(prefixed, pos+1)
		docScores2, pos := bs.booleanScores(prefixed, pos)
		return Intersect(docScores1, docScores2), pos
	case BOOL_OR:
		docScores1, pos := bs.booleanScores(prefixed, pos+1)
		docScores2, pos := bs.booleanScores(prefixed, pos)
		return union(docScores1, docScores2), pos
	case BOOL_NOT:
		docScores, pos := bs.booleanScores(prefixed, pos+1)
		return not(docScores, bs.Index.GetAllIds()), pos
	}
	return bs.Index.Get(prefixed[pos]), pos + 1
}

func Intersect(docScores1 []index.DocScore, docScores2 []index.DocScore) (Intersected []index.DocScore) {
	ind1, ind2 := 0, 0
	for ind1 != len(docScores1) && ind2 != len(docScores2) {
		if docScores1[ind1].Id == docScores2[ind2].Id {
			Intersected = append(Intersected,
				index.DocScore{Id: docScores1[ind1].Id,
					Score: docScores1[ind1].Score + docScores2[ind2].Score})
			ind1 += 1
			ind2 += 1
		} else if docScores1[ind1].Id < docScores2[ind2].Id {
			ind1 += 1
		} else {
			ind2 += 1
		}
	}
	return
}

func union(docScores1 []index.DocScore, docScores2 []index.DocScore) (unioned []index.DocScore) {
	ind1, ind2 := 0, 0
	for ind1 != len(docScores1) && ind2 != len(docScores2) {
		if docScores1[ind1].Id == docScores2[ind2].Id {
			unioned = append(unioned,
				index.DocScore{Id: docScores1[ind1].Id,
					Score: docScores1[ind1].Score + docScores2[ind2].Score})
			ind1 += 1
			ind2 += 1
		} else if docScores1[ind1].Id < docScores2[ind2].Id {
			unioned = append(unioned, docScores1[ind1])
			ind1 += 1
		} else {
			unioned = append(unioned, docScores2[ind2])
			ind2 += 1
		}
	}
	if ind1 != len(docScores1) {
		unioned = append(unioned, docScores1[ind1:]...)
	} else if ind2 != len(docScores2) {
		unioned = append(unioned, docScores2[ind2:]...)
	}
	return
}

func not(docScores []index.DocScore, ids []int) (notDocScores []index.DocScore) {
	ind1, ind2 := 0, 0
	for ind1 != len(docScores) && ind2 != len(ids) {
		if docScores[ind1].Id == ids[ind2] {
			ind1 += 1
			ind2 += 1
		} else {
			notDocScores = append(notDocScores,
				index.DocScore{Id: ids[ind2], Score: 1.0 / float64(len(ids)-len(docScores))})
			ind2 += 1
		}
	}
	for ind2 != len(ids) {
		notDocScores = append(notDocScores,
			index.DocScore{Id: ids[ind2], Score: 1.0 / float64(len(ids)-len(docScores))})
		ind2 += 1
	}
	return
}

func booleanPrefix(request string) (prefixedRequest []string) {
	// No real postfix, since it's already reversed
	postfixedRequest := []string{}
	opStack := new(Stack)
	for _, val := range prepareBooleanRequest(request) {
		switch {
		case val != BOOL_NOT && val != BOOL_AND && val != BOOL_OR && val != "(" && val != ")":
			postfixedRequest = append(postfixedRequest, val)
		case val == "(":
			opStack.Push(val)
		case val == ")":
			for opStack.Top() != "(" {
				postfixedRequest = append(postfixedRequest, opStack.Pop())
			}
			opStack.Pop()
		case val == BOOL_NOT:
			for opStack.Top() == BOOL_NOT {
				postfixedRequest = append(postfixedRequest, opStack.Pop())
			}
			opStack.Push(val)
		case val == BOOL_AND:
			for opStack.Top() == BOOL_NOT || opStack.Top() == BOOL_AND {
				postfixedRequest = append(postfixedRequest, opStack.Pop())
			}
			opStack.Push(val)
		case val == BOOL_OR:
			for opStack.Top() == BOOL_NOT || opStack.Top() == BOOL_AND || opStack.Top() == BOOL_OR {
				postfixedRequest = append(postfixedRequest, opStack.Pop())
			}
			opStack.Push(val)
		}
	}
	for opStack.Top() != "" {
		postfixedRequest = append(postfixedRequest, opStack.Pop())
	}
	return reverse(postfixedRequest)
}

func prepareBooleanRequest(request string) []string {
	preparedRequest := new(bytes.Buffer)
	for _, val := range strings.ToLower(request) {
		if val == '(' {
			preparedRequest.WriteString(" ) ")
		} else if val == ')' {
			preparedRequest.WriteString(" ( ")
		} else {
			preparedRequest.WriteRune(val)
		}
	}
	return reverse(strings.Fields(preparedRequest.String()))
}

func reverse(r []string) []string {
	for i, j := 0, len(r)-1; i < len(r)/2; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return r
}
