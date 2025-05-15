package utils

import (
	"errors"
	"gfly/app/utils"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestTransformList(t *testing.T) {
	tests := []struct {
		name          string
		input         []int
		transformerFn func(int) string
		expected      []string
	}{
		{
			name:          "Empty slice",
			input:         []int{},
			transformerFn: strconv.Itoa,
			expected:      []string{},
		},
		{
			name:          "Transform integers to strings",
			input:         []int{1, 2, 3},
			transformerFn: strconv.Itoa,
			expected:      []string{"1", "2", "3"},
		},
		{
			name:  "Transform with calculation",
			input: []int{1, 2, 3},
			transformerFn: func(i int) string {
				return strconv.Itoa(i * 2)
			},
			expected: []string{"2", "4", "6"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.TransformList(test.input, test.transformerFn)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestTransformMap(t *testing.T) {
	tests := []struct {
		name          string
		input         map[string]int
		transformerFn func(int) string
		expected      map[string]string
	}{
		{
			name:          "Empty map",
			input:         map[string]int{},
			transformerFn: strconv.Itoa,
			expected:      map[string]string{},
		},
		{
			name: "Transform integers to strings",
			input: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			transformerFn: strconv.Itoa,
			expected: map[string]string{
				"a": "1",
				"b": "2",
				"c": "3",
			},
		},
		{
			name: "Transform with calculation",
			input: map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
			},
			transformerFn: func(i int) string {
				return strconv.Itoa(i * 2)
			},
			expected: map[string]string{
				"a": "2",
				"b": "4",
				"c": "6",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.TransformMap(test.input, test.transformerFn)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestTransformListWithError(t *testing.T) {
	tests := []struct {
		name          string
		input         []int
		transformerFn func(int) (string, error)
		expected      []string
		expectedErrs  int
	}{
		{
			name:  "Empty slice",
			input: []int{},
			transformerFn: func(i int) (string, error) {
				return strconv.Itoa(i), nil
			},
			expected:     []string{},
			expectedErrs: 0,
		},
		{
			name:  "No errors",
			input: []int{1, 2, 3},
			transformerFn: func(i int) (string, error) {
				return strconv.Itoa(i), nil
			},
			expected:     []string{"1", "2", "3"},
			expectedErrs: 0,
		},
		{
			name:  "Some errors",
			input: []int{1, 2, 3, 4, 5},
			transformerFn: func(i int) (string, error) {
				if i%2 == 0 {
					return "", errors.New("even number error")
				}
				return strconv.Itoa(i), nil
			},
			expected:     []string{"1", "3", "5"},
			expectedErrs: 2,
		},
		{
			name:  "All errors",
			input: []int{1, 2, 3},
			transformerFn: func(i int) (string, error) {
				return "", errors.New("error")
			},
			expected:     []string{},
			expectedErrs: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, errs := utils.TransformListWithError(test.input, test.transformerFn)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected result %v, got %v", test.expected, result)
			}
			if len(errs) != test.expectedErrs {
				t.Errorf("Expected %d errors, got %d", test.expectedErrs, len(errs))
			}
		})
	}
}

func TestTransformConcurrent(t *testing.T) {
	tests := []struct {
		name          string
		input         []int
		transformerFn func(int) string
		numWorkers    int
		expected      []string
	}{
		{
			name:          "Empty slice",
			input:         []int{},
			transformerFn: strconv.Itoa,
			numWorkers:    2,
			expected:      []string{},
		},
		{
			name:          "Single worker",
			input:         []int{1, 2, 3},
			transformerFn: strconv.Itoa,
			numWorkers:    1,
			expected:      []string{"1", "2", "3"},
		},
		{
			name:          "Multiple workers",
			input:         []int{1, 2, 3, 4, 5, 6},
			transformerFn: strconv.Itoa,
			numWorkers:    3,
			expected:      []string{"1", "2", "3", "4", "5", "6"},
		},
		{
			name:          "More workers than items",
			input:         []int{1, 2, 3},
			transformerFn: strconv.Itoa,
			numWorkers:    5,
			expected:      []string{"1", "2", "3"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.TransformConcurrent(test.input, test.transformerFn, test.numWorkers)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestTransformConcurrentRaceCondition(t *testing.T) {
	// This test checks for race conditions by having transformers that access shared state
	input := make([]int, 1000)
	for i := range input {
		input[i] = i
	}

	var counter int
	var mu sync.Mutex

	transformerFn := func(i int) string {
		mu.Lock()
		counter++
		mu.Unlock()
		time.Sleep(time.Millisecond) // Introduce some delay to increase chance of race
		return strconv.Itoa(i)
	}

	result := utils.TransformConcurrent(input, transformerFn, 10)

	if len(result) != len(input) {
		t.Errorf("Expected result length %d, got %d", len(input), len(result))
	}

	if counter != len(input) {
		t.Errorf("Expected counter to be %d, got %d", len(input), counter)
	}
}

func TestTransformBatch(t *testing.T) {
	tests := []struct {
		name          string
		input         []int
		transformerFn func([]int) []string
		batchSize     int
		expected      []string
	}{
		{
			name:  "Empty slice",
			input: []int{},
			transformerFn: func(batch []int) []string {
				result := make([]string, len(batch))
				for i, v := range batch {
					result[i] = strconv.Itoa(v)
				}
				return result
			},
			batchSize: 2,
			expected:  []string{},
		},
		{
			name:  "Default batch size",
			input: []int{1, 2, 3},
			transformerFn: func(batch []int) []string {
				result := make([]string, len(batch))
				for i, v := range batch {
					result[i] = strconv.Itoa(v)
				}
				return result
			},
			batchSize: 0, // Should use default
			expected:  []string{"1", "2", "3"},
		},
		{
			name:  "Exact batch size",
			input: []int{1, 2, 3, 4},
			transformerFn: func(batch []int) []string {
				result := make([]string, len(batch))
				for i, v := range batch {
					result[i] = strconv.Itoa(v)
				}
				return result
			},
			batchSize: 2,
			expected:  []string{"1", "2", "3", "4"},
		},
		{
			name:  "Partial last batch",
			input: []int{1, 2, 3, 4, 5},
			transformerFn: func(batch []int) []string {
				result := make([]string, len(batch))
				for i, v := range batch {
					result[i] = strconv.Itoa(v)
				}
				return result
			},
			batchSize: 2,
			expected:  []string{"1", "2", "3", "4", "5"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := utils.TransformBatch(test.input, test.transformerFn, test.batchSize)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

// TestTransformBatchVerifyBatching verifies that the batching logic works correctly
func TestTransformBatchVerifyBatching(t *testing.T) {
	input := []int{1, 2, 3, 4, 5, 6, 7}
	batchSize := 3

	var batches [][]int

	transformerFn := func(batch []int) []string {
		batches = append(batches, batch)
		result := make([]string, len(batch))
		for i, v := range batch {
			result[i] = strconv.Itoa(v)
		}
		return result
	}

	result := utils.TransformBatch(input, transformerFn, batchSize)

	// Check result
	expected := []string{"1", "2", "3", "4", "5", "6", "7"}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected result %v, got %v", expected, result)
	}

	// Check batches
	expectedBatches := [][]int{
		{1, 2, 3},
		{4, 5, 6},
		{7},
	}

	if len(batches) != len(expectedBatches) {
		t.Errorf("Expected %d batches, got %d", len(expectedBatches), len(batches))
	}

	for i, batch := range batches {
		if !reflect.DeepEqual(batch, expectedBatches[i]) {
			t.Errorf("Expected batch %v, got %v", expectedBatches[i], batch)
		}
	}
}
