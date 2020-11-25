package util

import "time"

// String returns a pointer to the passed string s.
func String(s string) *string { return &s }

// Int returns a pointer to the passed integer i.
func Int(i int) *int { return &i }

// Int64 returns a pointer to the passed integer i.
func Int64(i int64) *int64 { return &i }

// Bool returns a pointer to the passed bool b.
func Bool(b bool) *bool { return &b }

// Duration returns a pointer to the passed time.Duration d.
func Duration(d time.Duration) *time.Duration { return &d }

// IntSlice returns a pointer to the passed integer slice is
func IntSlice(is []int) *[]int { return &is }

// StringSlice returns a pointer to the passed string slice ss
func StringSlice(ss []string) *[]string { return &ss }
