package spy

import (
	"testing"
)

func TestParsePage(t *testing.T) {
	page := ParsePage("list_23_171.html")
	t.Logf("find page %d", page)
}