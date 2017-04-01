package scraper

import "testing"

var casesVisits = []struct {
	u   string
	out bool
}{
	{"/page/1", false},
	{"/page/1", true},
	{"/page/2", false},
	{"/page/2", true},
	{"", false},
	{"", true},
	{"javascript:alert(1)", false},
	{"/", false},
}

func TestVisitsMemory(t *testing.T) {
	v := NewMemoryVisits()
	defer func() {
		v.Close()
	}()
	for _, c := range casesVisits {
		ok := v.Visit(c.u)
		if ok != c.out {
			t.Fatalf("Visit(%s) = %v, want %v", c.u, ok, c.out)
		}
	}
	if err := v.ResetVisit("/page/1"); err != nil {
		t.Fatal(err)
	}
	ok := v.Visit("/page/1")
	if ok {
		t.Fatalf("ResetVisit('/page/1') = %v, want = false", ok)
	}
	if err := v.Drop(); err != nil {
		t.Fatal(err)
	}
	for _, c := range casesVisits {
		ok := v.Visit(c.u)
		if ok != c.out {
			t.Fatalf("Visit(%s) = %v, want %v", c.u, ok, c.out)
		}
	}
}
