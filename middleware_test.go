package fabyscore

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMiddlewareSort(t *testing.T) {
	m0 := middleware{sort: -255}
	m1 := middleware{sort: -5}
	m2 := middleware{sort: 0}
	m3 := middleware{sort: 100}
	m4 := middleware{sort: 255}

	middlewares := middlewares{
		m2, m1, m3, m0, m4,
	}

	sort.Sort(middlewares)

	assert.Equal(t, m0, middlewares[0])
	assert.Equal(t, m1, middlewares[1])
	assert.Equal(t, m2, middlewares[2])
	assert.Equal(t, m3, middlewares[3])
	assert.Equal(t, m4, middlewares[4])
}
