package scraper

import (
	"os"
	"testing"
)

func TestVisitsBoltDb(t *testing.T) {
	dbpath := "/tmp/visits.db"
	v, err := NewBoltDbVisits(dbpath)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		os.Remove(dbpath)
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
