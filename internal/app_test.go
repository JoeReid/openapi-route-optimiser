package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessPath(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected string
		testErr  func(assert.TestingT, error, ...interface{}) bool
	}{
		{
			name:     "No braces",
			path:     "/api/v1/users",
			expected: "/api/v1/users",
			testErr:  assert.NoError,
		},
		{
			name:     "Single variable",
			path:     "/api/v1/users/{userId}",
			expected: "/api/v1/users/*",
			testErr:  assert.NoError,
		},
		{
			name:     "Multiple variables",
			path:     "/api/v1/users/{userId}/posts/{postId}",
			expected: "/api/v1/users/*/posts/*",
			testErr:  assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processPath(tt.path)
			tt.testErr(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProcessTags(t *testing.T) {
	tests := []struct {
		name     string
		tags     []string
		filter   string
		find     string
		replace  string
		expected []string
		testErr  func(assert.TestingT, error, ...interface{}) bool
	}{
		{
			name:     "No filter, no find or replace",
			tags:     []string{"tag1", "tag2", "tag3"},
			filter:   "",
			find:     "",
			replace:  "",
			expected: []string{"tag1", "tag2", "tag3"},
			testErr:  assert.NoError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := processTags(tt.tags, tt.filter, tt.find, tt.replace)
			tt.testErr(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
