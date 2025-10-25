package scraper

import "testing"

func TestIsEarlyCareer(t *testing.T) {
	cases := []struct {
		title string
		want  bool
	}{
		{"Software Engineer", true},
		{"Junior Data Analyst", true},
		{"New Grad Backend Engineer", true},
		{"Senior Software Engineer", false},
		{"Staff Platform Engineer", false},
		{"Director of Engineering", false},
	}
	for _, c := range cases {
		if got := isEarlyCareer(c.title); got != c.want {
			t.Errorf("isEarlyCareer(%q) = %v, want %v", c.title, got, c.want)
		}
	}
}

func TestIsInUSA(t *testing.T) {
	cases := []struct {
		loc  string
		want bool
	}{
		{"San Francisco, CA, United States", true},
		{"Remote - US", true},
		{"Austin, TX, USA", true},
		{"Toronto, Canada", false},
		{"London, UK", false},
	}
	for _, c := range cases {
		if got := isInUSA(c.loc); got != c.want {
			t.Errorf("isInUSA(%q) = %v, want %v", c.loc, got, c.want)
		}
	}
}
