package got_test

import (
	"testing"

	"github.com/gophero/got"
)

func TestWrap(t *testing.T) {
	tr := got.Wrap(t)
	tr.Case("wrapping testing.T")
	tr.Require(tr != nil, "wrapping should be success")
}
