package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	index := New()

	assert.Equal(t, map[string]*[]int{}, index.comics)

	input := map[string]*[]int{"word": {1, 2}}
	index.Update(input)

	assert.Equal(t, input, index.comics)
}

func TestGetComicsIdsByWords(t *testing.T) {
	index := New()

	_, err := index.GetComicsIdsByWords(t.Context(), []string{"word"})
	assert.Error(t, err, "no not found error provided in empty index")

	index.Update(map[string]*[]int{"word": {1, 2}})

	res, err := index.GetComicsIdsByWords(t.Context(), []string{"word"})
	require.NoError(t, err)

	assert.ElementsMatch(t, res, []int{1, 2})
}
