package function

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListFunctionNames(t *testing.T) {
	type result struct {
		fnNames     []string
		numberOfFns int
	}
	testSuit := map[string]result{
		"function1,function2": {
			fnNames:     []string{"function1", "function2"},
			numberOfFns: 2,
		},
		"fn1,  fn2,  fn3 ": {
			fnNames:     []string{"fn1", "fn2", "fn3"},
			numberOfFns: 3,
		},
		"fn1_,fn2, fn3, fn4 ": {
			fnNames:     []string{"fn1_", "fn2", "fn3", "fn4"},
			numberOfFns: 4,
		},
		"fn1_,": {
			fnNames:     []string{"fn1_"},
			numberOfFns: 1,
		},
		"fn1": {
			fnNames:     []string{"fn1"},
			numberOfFns: 1,
		},
	}

	for test, expected := range testSuit {
		fnList := listFunctionNames(test)
		assert.Equal(t, len(fnList), expected.numberOfFns)
		assert.Equal(t, fnList, expected.fnNames)
	}
}
