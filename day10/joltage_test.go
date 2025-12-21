package main

import (
	"testing"
)

func Test_joltage_equals(t *testing.T) {
	data := []struct {
		name     string
		j1       Joltage
		j2       Joltage
		expected bool
	}{
		{"equal", NewJoltage([]int{1, 2}), NewJoltage([]int{1, 2}), true},
		{"not_equal", NewJoltage([]int{1, 2}), NewJoltage([]int{3, 4}), false},
		{"different_lengths", NewJoltage([]int{1}), NewJoltage([]int{1, 2}), false},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			result := d.j1.Equals(d.j2)
			if result != d.expected {
				t.Errorf("Expected %v, got %v", d.expected, result)
			}
		})
	}
}

func Test_transformJoltage(t *testing.T) {
	data := []struct {
		name     string
		button   Button
		j        Joltage
		expected Joltage
		errMsg   string
	}{
		{"transform", NewButton([]bool{true, false}), NewJoltage([]int{0, 0}), NewJoltage([]int{1, 0}), ""},
		{"transform_long", NewButton([]bool{true, false, true, false, true, false, true, false}), NewJoltage([]int{1, 1, 1, 1, 1, 1, 1, 1}), NewJoltage([]int{2, 1, 2, 1, 2, 1, 2, 1}), ""},
		{"transform_identity", NewButton([]bool{false, false, false}), NewJoltage([]int{1, 1, 1}), NewJoltage([]int{1, 1, 1}), ""},
		{"wrong_length", NewButton([]bool{false, false, false}), NewJoltage([]int{1}), Joltage{}, "Button cannot interact with joltage of different length. Button length: 3, joltage length: 1"},
	}
	for _, d := range data {
		t.Run(d.name, func(t *testing.T) {
			transformed, err := d.button.TransformJoltage(d.j)
			if !transformed.Equals(d.expected) {
				t.Errorf("Expected %v, got %v", d.expected, transformed)
			}
			var errMsg string
			if err != nil {
				errMsg = err.Error()
			}
			if errMsg != d.errMsg {
				t.Errorf("Expected %v, got %v", d.errMsg, errMsg)
			}

		})
	}
}

func Test_(t *testing.T) {
}
