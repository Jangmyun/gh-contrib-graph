package contrib

import "testing"

func TestLevel(t *testing.T) {
	cases := []struct {
		name  string
		count int
		max   int
		want  int
	}{
		{"zero count", 0, 10, 0},
		{"zero max", 5, 0, 0},
		{"quartile boundary 1", 25, 100, 1},
		{"quartile boundary 2", 50, 100, 2},
		{"quartile boundary 3", 75, 100, 3},
		{"max count", 100, 100, 4},
		{"just above q1", 26, 100, 2},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := Level(c.count, c.max)
			if got != c.want {
				t.Errorf("Level(%d, %d) = %d, want %d", c.count, c.max, got, c.want)
			}
		})
	}
}

func TestMaxCount(t *testing.T) {
	days := days(0, 3, 7, 2)
	if got := MaxCount(days); got != 7 {
		t.Errorf("MaxCount = %d, want 7", got)
	}
	if got := MaxCount(nil); got != 0 {
		t.Errorf("MaxCount(nil) = %d, want 0", got)
	}
}
