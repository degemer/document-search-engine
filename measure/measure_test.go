package measure

import (
	"github.com/degemer/document-search-engine/index"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInExpected(t *testing.T) {
	expectedResults := []index.DocScore{index.DocScore{Id: 3}, index.DocScore{Id: 5}, index.DocScore{Id: 8}}
	assert.True(t, inExpected(index.DocScore{Id: 8}, expectedResults))
	assert.True(t, inExpected(index.DocScore{Id: 3}, expectedResults))
	assert.True(t, inExpected(index.DocScore{Id: 5}, expectedResults))

	assert.False(t, inExpected(index.DocScore{Id: 7}, expectedResults))
	assert.False(t, inExpected(index.DocScore{Id: 1}, expectedResults))
	assert.False(t, inExpected(index.DocScore{Id: 10}, expectedResults))
}

func TestPrecisionRappel(t *testing.T) {
	results := []index.DocScore{index.DocScore{Id: 3}, index.DocScore{Id: 5}, index.DocScore{Id: 8}}
	expectedResults := []index.DocScore{index.DocScore{Id: 4}, index.DocScore{Id: 3}}
	prec, rapp := precisionRappel(results, expectedResults)
	assert.Equal(t, prec, 1.0/3.0)
	assert.Equal(t, rapp, 1.0/2.0)

	expectedResults = []index.DocScore{}
	prec, rapp = precisionRappel(results, expectedResults)
	assert.Equal(t, prec, 0.0)
	assert.Equal(t, rapp, 1.0)

	results = []index.DocScore{}
	prec, rapp = precisionRappel(results, expectedResults)
	assert.Equal(t, prec, 1.0)
	assert.Equal(t, rapp, 1.0)

	expectedResults = []index.DocScore{index.DocScore{Id: 3}, index.DocScore{Id: 4}}
	prec, rapp = precisionRappel(results, expectedResults)
	assert.Equal(t, prec, 0.0)
	assert.Equal(t, rapp, 0.0)
}
