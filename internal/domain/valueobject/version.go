package valueobject

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
)

var (
	// ErrInvalidVersion is returned when an invalid version is provided.
	ErrInvalidVersion = errors.New("invalid version")
	// ErrEmptyVersion is returned when an empty version is provided.
	ErrEmptyVersion = errors.New("version cannot be empty")
	// versionRegex is the regex pattern for semantic versioning (semver).
	// Format: MAJOR.MINOR.PATCH (e.g., "1.0.0", "2.3.1")
	versionRegex = regexp.MustCompile(`^\d+\.\d+\.\d+$`)
)

// Version represents a semantic version.
// It follows the MAJOR.MINOR.PATCH format (e.g., "1.0.0").
type Version string

// String returns the string representation of the version.
func (v Version) String() string {
	return string(v)
}

// IsEmpty returns true if the version is empty.
func (v Version) IsEmpty() bool {
	return string(v) == ""
}

// IsValid checks if the version is valid (not empty and matches semver pattern).
func (v Version) IsValid() bool {
	if v.IsEmpty() {
		return false
	}
	return versionRegex.MatchString(string(v))
}

// MarshalJSON implements json.Marshaler interface.
func (v Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(v))
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (v *Version) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	if str == "" {
		return ErrEmptyVersion
	}
	version := Version(str)
	if !version.IsValid() {
		return fmt.Errorf("%w: %s", ErrInvalidVersion, str)
	}
	*v = version
	return nil
}

// NewVersion creates a new Version from a string.
// Returns an error if the string is not a valid version.
func NewVersion(version string) (Version, error) {
	if version == "" {
		return "", ErrEmptyVersion
	}
	v := Version(version)
	if !v.IsValid() {
		return "", ErrInvalidVersion
	}
	return v, nil
}

// MustNewVersion creates a new Version from a string.
// Panics if the string is not a valid version.
func MustNewVersion(version string) Version {
	v, err := NewVersion(version)
	if err != nil {
		panic(err)
	}
	return v
}
