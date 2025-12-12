package stringsx

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClip_EdgeCases(t *testing.T) {
	cases := []struct {
		name string
		s    string
		max  int
		want string
	}{

		{"clip to 5", "Hello, World!", 5, "Hello"},
		{"full string", "Hello", 10, "Hello"},
		{"empty string", "", 5, ""},

		{"max = 0", "Hello", 0, ""},
		{"max = len(s)", "Hello", 5, "Hello"},
		{"max < 0", "Hello", -1, ""},
		{"max < -5", "Hello", -100, ""},

		{"with spaces", "Hello World", 7, "Hello W"},
		{"unicode", "Привет", 2, "Пр"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Clip(c.s, c.max)
			require.Equal(t, c.want, got, "Clip(%q, %d) = %q", c.s, c.max, got)
		})
	}
}

func TestIsEmpty_Table(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want bool
	}{

		{"empty string", "", true},
		{"spaces", "   ", true},
		{"tabs", "\t\t", true},
		{"newlines", "\n\n", true},
		{"mixed whitespace", " \t\n ", true},

		{"single char", "a", false},
		{"word", "hello", false},
		{"spaces with text", "  hello  ", false},
		{"number", "123", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := IsEmpty(c.s)
			assert.Equal(t, c.want, got, "IsEmpty(%q)", c.s)
		})
	}
}

func TestReverse_Table(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want string
	}{
		{"simple", "Hello", "olleH"},
		{"empty", "", ""},
		{"single char", "a", "a"},
		{"palindrome", "racecar", "racecar"},
		{"with spaces", "a b c", "c b a"},
		{"unicode", "Привет", "тевирП"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Reverse(c.s)
			assert.Equal(t, c.want, got, "Reverse(%q) = %q", c.s, got)
		})
	}
}

func TestCountWords_Table(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want int
	}{
		{"single word", "hello", 1},
		{"two words", "hello world", 2},
		{"multiple words", "one two three four", 4},
		{"empty string", "", 0},
		{"spaces only", "   ", 0},
		{"extra spaces", "hello   world", 2},
		{"leading/trailing spaces", "  hello world  ", 2},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := CountWords(c.s)
			assert.Equal(t, c.want, got, "CountWords(%q) = %d", c.s, got)
		})
	}
}

func TestCapitalize_Table(t *testing.T) {
	cases := []struct {
		name string
		s    string
		want string
	}{
		{"lowercase", "hello", "Hello"},
		{"uppercase", "Hello", "Hello"},
		{"empty", "", ""},
		{"single char", "a", "A"},
		{"already capital", "Hello world", "Hello world"},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Capitalize(c.s)
			assert.Equal(t, c.want, got, "Capitalize(%q) = %q", c.s, got)
		})
	}
}

func BenchmarkClip(b *testing.B) {
	s := "Hello, World! This is a test string."
	for i := 0; i < b.N; i++ {
		_ = Clip(s, 10)
	}
}

func BenchmarkIsEmpty(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = IsEmpty("   ")
	}
}

func BenchmarkReverse(b *testing.B) {
	s := "Hello, World!"
	for i := 0; i < b.N; i++ {
		_ = Reverse(s)
	}
}

func BenchmarkCountWords(b *testing.B) {
	s := "one two three four five six seven eight"
	for i := 0; i < b.N; i++ {
		_ = CountWords(s)
	}
}
