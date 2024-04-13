package cmd

import "testing"

func TestCheckOutOfBounds(t *testing.T) {
	cases := []struct {
		input struct {
			muscle []string
			num    int
		}
		expected []string
	}{
		{
			input: struct{muscle []string; num int}{
				muscle: []string{"hamstrings", "glutes", "quadriceps", "calves"},
				num: 5,
			},
			expected: []string{"hamstrings", "glutes", "quadriceps", "calves", "hamstrings"},
		},
		{
			input: struct{muscle []string; num int}{
				muscle: []string{"hamstrings", "glutes", "quadriceps", "calves"},
				num: 3,
			},
			expected: []string{"hamstrings", "glutes", "quadriceps"},
		},
		{
			input: struct{muscle []string; num int}{
				muscle: []string{"hamstrings", "glutes", "quadriceps", "calves"},
				num: 4,
			},
			expected: []string{"hamstrings", "glutes", "quadriceps", "calves"},
		},
	}

	for _, c := range cases {
		actual := checkOutOfBounds(c.input.muscle, c.input.num)
		if len(actual) != len(c.expected) {
			t.Errorf("The lengths are not equal: %v vs %v,",
				len(actual),
				len(c.expected),
			)
			continue
		}
		for i := range actual {
			actualExercise := actual[i]
			expectedExercise := c.expected[i]
			if actualExercise != expectedExercise {
				t.Errorf("%v does not equal %v", 
					actualExercise, 
					expectedExercise,
				)
			}
		}
	}
}
