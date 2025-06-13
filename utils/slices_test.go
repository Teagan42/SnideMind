package utils

import (
	"testing"
)

func TestIntersection_Ints(t *testing.T) {
	tests := []struct {
		name     string
		a        []int
		b        []int
		expected []int
	}{
		{"BothEmpty", []int{}, []int{}, []int{}},
		{"AEmpty", []int{}, []int{1, 2, 3}, []int{}},
		{"BEmpty", []int{1, 2, 3}, []int{}, []int{}},
		{"NoIntersection", []int{1, 2, 3}, []int{4, 5, 6}, []int{}},
		{"SomeIntersection", []int{1, 2, 3}, []int{2, 3, 4}, []int{2, 3}},
		{"AllIntersection", []int{1, 2, 3}, []int{1, 2, 3}, []int{1, 2, 3}},
		{"DuplicatesInA", []int{1, 2, 2, 3}, []int{2, 3}, []int{2, 2, 3}},
		{"DuplicatesInB", []int{1, 2, 3}, []int{2, 2, 3, 3}, []int{2, 3}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Intersection(tt.a, tt.b)
			if len(got) != len(tt.expected) {
				t.Errorf("Intersection(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
				return
			}
			for i := range got {
				if got[i] != tt.expected[i] {
					t.Errorf("Intersection(%v, %v) = %v, want %v", tt.a, tt.b, got, tt.expected)
					break
				}
			}
		})
	}
}

func TestIntersection_Strings(t *testing.T) {
	a := []string{"apple", "banana", "cherry"}
	b := []string{"banana", "date", "cherry"}
	expected := []string{"banana", "cherry"}

	got := Intersection(a, b)
	if len(got) != len(expected) {
		t.Errorf("Intersection(%v, %v) = %v, want %v", a, b, got, expected)
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("Intersection(%v, %v) = %v, want %v", a, b, got, expected)
		}
	}
}

func TestIntersection_Structs(t *testing.T) {
	type point struct{ X, Y int }
	a := []point{{1, 2}, {3, 4}, {5, 6}}
	b := []point{{3, 4}, {7, 8}, {1, 2}}
	expected := []point{{1, 2}, {3, 4}}

	got := Intersection(a, b)
	if len(got) != len(expected) {
		t.Errorf("Intersection(%v, %v) = %v, want %v", a, b, got, expected)
	}
	for i := range got {
		if got[i] != expected[i] {
			t.Errorf("Intersection(%v, %v) = %v, want %v", a, b, got, expected)
		}
	}
}
