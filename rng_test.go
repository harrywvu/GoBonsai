package main

import (
	"reflect"
	"testing"
)

func TestPythonRandomFloat64(t *testing.T) {
	r := newPythonRandom(1)
	want := []float64{
		0.13436424411240122,
		0.8474337369372327,
		0.763774618976614,
	}

	for i, expected := range want {
		if got := r.Float64(); got != expected {
			t.Fatalf("Float64 call %d = %.17g, want %.17g", i, got, expected)
		}
	}
}

func TestPythonRandomIntN(t *testing.T) {
	r := newPythonRandom(1)
	limits := []int{4, 4, 5, 181}
	want := []int{1, 0, 2, 30}

	for i, limit := range limits {
		if got := r.IntN(limit); got != want[i] {
			t.Fatalf("IntN call %d = %d, want %d", i, got, want[i])
		}
	}
}

func TestPythonRandomNormal(t *testing.T) {
	r := newPythonRandom(1)
	want := []float64{
		0.6074558576437062,
		-0.01422544551078489,
		1.2309072291166607,
	}

	for i, expected := range want {
		if got := r.NormFloat64(); got != expected {
			t.Fatalf("NormFloat64 call %d = %.17g, want %.17g", i, got, expected)
		}
	}
}

func TestGenerateBranchNumsMatchesPython(t *testing.T) {
	rng.Seed(1)
	got := generateBranchNums(fibNums(5), 5)
	want := [][]int{{1}, {2}, {1, 2}, {2, 2, 1}, {1, 1, 2, 2, 2}, {2, 1, 2, 2, 2, 2, 1, 1}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("generateBranchNums() = %v, want %v", got, want)
	}
}

func TestSeedOneInitialFibStateMatchesPython(t *testing.T) {
	rng.Seed(1)
	rng.IntN(4)
	gotBranches := generateBranchNums(fibNums(5), 5)
	wantBranches := [][]int{{1}, {2}, {1, 2}, {1, 2, 2}, {2, 2, 2, 1, 1}, {2, 2, 2, 1, 1, 1, 2, 2}}
	if !reflect.DeepEqual(gotBranches, wantBranches) {
		t.Fatalf("generateBranchNums() = %v, want %v", gotBranches, wantBranches)
	}
	angle := rng.NormFloat64() * angleStdDev

	if angle != 0.06393670624066468 {
		t.Fatalf("initial angle = %.17g", angle)
	}
}
