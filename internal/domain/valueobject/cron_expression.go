package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var (
	// ErrInvalidCronExpression is returned when an invalid cron expression is provided.
	ErrInvalidCronExpression = errors.New("invalid cron expression")
	// ErrEmptyCronExpression is returned when an empty cron expression is provided.
	ErrEmptyCronExpression = errors.New("cron expression cannot be empty")
	// cronRegex is the regex pattern for cron expressions.
	// Format: MIN HOUR DAY MONTH WEEKDAY (e.g., "0 * * * *", "*/5 * * * *")
	// Supports: minute (0-59), hour (0-23), day (1-31), month (1-12), weekday (0-6)
	cronRegex = regexp.MustCompile(`^(\S+)\s+(\S+)\s+(\S+)\s+(\S+)\s+(\S+)$`)
)

// validateCronPart validates a single part of a cron expression.
func validateCronPart(part string, min, max int) bool {
	if part == "*" {
		return true
	}
	if strings.HasPrefix(part, "*/") {
		// Step format: */n
		n, err := strconv.Atoi(part[2:])
		return err == nil && n > 0
	}
	if strings.Contains(part, "-") {
		// Range format: a-b
		rangeParts := strings.Split(part, "-")
		if len(rangeParts) != 2 {
			return false
		}
		start, err1 := strconv.Atoi(rangeParts[0])
		end, err2 := strconv.Atoi(rangeParts[1])
		if err1 != nil || err2 != nil {
			return false
		}
		return start >= min && start <= max && end >= min && end <= max && start <= end
	}
	if strings.Contains(part, ",") {
		// List format: a,b,c
		listParts := strings.Split(part, ",")
		for _, p := range listParts {
			num, err := strconv.Atoi(p)
			if err != nil || num < min || num > max {
				return false
			}
		}
		return true
	}
	// Single number
	num, err := strconv.Atoi(part)
	return err == nil && num >= min && num <= max
}

// CronExpression represents a cron expression for scheduling.
// It follows the standard cron format: MINUTE HOUR DAY MONTH WEEKDAY
type CronExpression string

// String returns the string representation of the cron expression.
func (c CronExpression) String() string {
	return string(c)
}

// IsEmpty returns true if the cron expression is empty.
func (c CronExpression) IsEmpty() bool {
	return string(c) == ""
}

// IsValid checks if the cron expression is valid (not empty and matches cron pattern).
func (c CronExpression) IsValid() bool {
	if c.IsEmpty() {
		return false
	}
	matches := cronRegex.FindStringSubmatch(string(c))
	if matches == nil {
		return false
	}
	// Validate each part with its range
	// matches[1] = minute (0-59)
	// matches[2] = hour (0-23)
	// matches[3] = day (1-31)
	// matches[4] = month (1-12)
	// matches[5] = weekday (0-6)
	return validateCronPart(matches[1], 0, 59) &&
		validateCronPart(matches[2], 0, 23) &&
		validateCronPart(matches[3], 1, 31) &&
		validateCronPart(matches[4], 1, 12) &&
		validateCronPart(matches[5], 0, 6)
}

// MarshalJSON implements json.Marshaler interface.
func (c CronExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(c))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (c *CronExpression) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyCronExpression
	}
	expr := CronExpression(str)
	if !expr.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidCronExpression, str)
	}
	*c = expr
	return nil
}

// NewCronExpression creates a new CronExpression from a string.
// Returns an error if the string is not a valid cron expression.
func NewCronExpression(expr string) (CronExpression, error) {
	if expr == "" {
		return "", ErrEmptyCronExpression
	}
	c := CronExpression(expr)
	if !c.IsValid() {
		return "", ErrInvalidCronExpression
	}
	return c, nil
}

// MustNewCronExpression creates a new CronExpression from a string.
// Panics if the string is not a valid cron expression.
func MustNewCronExpression(expr string) CronExpression {
	c, err := NewCronExpression(expr)
	if err != nil {
		panic(err)
	}
	return c
}
