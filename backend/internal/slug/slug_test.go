package slug

import "testing"

func TestMake(t *testing.T) {
	for _, tt := range []struct{ in, out string }{{"Hello, World!", "hello-world"}, {"  A  B ", "a-b"}, {"日本語", ""}} {
		if got := Make(tt.in); got != tt.out {
			t.Fatalf("Make(%q)=%q, want %q", tt.in, got, tt.out)
		}
	}
}
