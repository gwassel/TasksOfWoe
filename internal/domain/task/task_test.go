package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTaskStatus(t *testing.T) {
	tests := []struct {
		name     string
		task     Task
		expected taskStatus
	}{
		{
			name: "incomplete task",
			task: Task{
				Completed: false,
				InWork:    false,
			},
			expected: Incomplete,
		},
		{
			name: "working task",
			task: Task{
				Completed: false,
				InWork:    true,
			},
			expected: Working,
		},
		{
			name: "completed task",
			task: Task{
				Completed: true,
				InWork:    false,
			},
			expected: Completed,
		},
		{
			name: "completed and in work (completion takes priority)",
			task: Task{
				Completed: true,
				InWork:    true,
			},
			expected: Completed,
		},
		{
			name:     "default task",
			task:     Task{},
			expected: Incomplete,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.task.Status()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFormatDateForTask(t *testing.T) {
	tests := []struct {
		name      string
		date      string
		expectErr bool
		contains  []string
	}{
		{
			name:      "valid date",
			date:      "2024-01-15T10:30:45.123456789Z",
			expectErr: false,
			contains:  []string{"2024-01-15", "Moscow"},
		},
		{
			name:      "valid date with timezone offset",
			date:      "2024-01-15T10:30:45.123456789+03:00",
			expectErr: false,
			contains:  []string{"2024-01-15", "10:30:45", "Moscow"},
		},
		{
			name:      "empty date",
			date:      "",
			expectErr: true,
		},
		{
			name:      "invalid date format",
			date:      "2024-01-15",
			expectErr: true,
		},
		{
			name:      "invalid date string",
			date:      "not a date",
			expectErr: true,
		},
		{
			name:      "date with milliseconds",
			date:      "2024-12-31T23:59:59.999Z",
			expectErr: false,
			contains:  []string{"Moscow"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := FormatDateForTask(tt.date)

			if tt.expectErr {
				assert.Error(t, err)
				assert.Empty(t, result)
			} else {
				assert.NoError(t, err)
				for _, expected := range tt.contains {
					assert.Contains(t, result, expected)
				}
			}
		})
	}
}

func TestTaskStatusToString(t *testing.T) {
	tests := []struct {
		name     string
		status   taskStatus
		expected string
	}{
		{
			name:     "incomplete",
			status:   Incomplete,
			expected: "Incomplete",
		},
		{
			name:     "working",
			status:   Working,
			expected: "Working",
		},
		{
			name:     "completed",
			status:   Completed,
			expected: "Completed",
		},
		{
			name:     "invalid status",
			status:   taskStatus(99),
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.status.ToString()
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestTask_Precedence(t *testing.T) {
	t.Run("completed takes precedence over working", func(t *testing.T) {
		task := Task{
			Completed: true,
			InWork:    true,
		}
		status := task.Status()
		assert.Equal(t, Completed, status)
	})

	t.Run("working takes precedence over incomplete", func(t *testing.T) {
		task := Task{
			Completed: false,
			InWork:    true,
		}
		status := task.Status()
		assert.Equal(t, Working, status)
	})

	t.Run("incomplete when neither working nor completed", func(t *testing.T) {
		task := Task{
			Completed: false,
			InWork:    false,
		}
		status := task.Status()
		assert.Equal(t, Incomplete, status)
	})
}

func TestFormatDateForTask_MoscowTimezone(t *testing.T) {
	utcTime := time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	dateUTC := utcTime.Format(time.RFC3339Nano)

	result, err := FormatDateForTask(dateUTC)
	assert.NoError(t, err)

	assert.Contains(t, result, "Moscow")
}
