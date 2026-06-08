package utils

import (
	"errors"
	"fmt"
	"strings"
)

// Sentinel validation errors. Library consumers should match these with
// errors.Is to handle classes of bad input.
var (
	// ErrNoComponents is returned when no components are supplied.
	ErrNoComponents = errors.New("at least one component required")

	// ErrLanguageRequired is returned when the configured components include
	// any of LanguageRequiredFor but Language is empty.
	ErrLanguageRequired = errors.New("language required for non-collector components")

	// ErrManualInstrumentationFile is returned when js + manual-instrumentation
	// is configured but InstrumentationFile is empty.
	ErrManualInstrumentationFile = errors.New("manual-instrumentation requires an instrumentation file")
)

// UnsupportedLanguageError reports that the configured Language is not in
// SupportedLanguages. Use errors.As to extract the offending value.
type UnsupportedLanguageError struct{ Language string }

func (e *UnsupportedLanguageError) Error() string {
	return fmt.Sprintf("language %q not supported. Possible values: %s",
		e.Language, strings.Join(SupportedLanguages, ", "))
}

// UnsupportedComponentError reports that one of the configured Components
// is not in SupportedComponents.
type UnsupportedComponentError struct{ Component string }

func (e *UnsupportedComponentError) Error() string {
	return fmt.Sprintf("component %q not supported. Possible values: %s",
		e.Component, strings.Join(SupportedComponents, ", "))
}

// UnsupportedFormatError reports that the configured Format is not in
// SupportedFormats.
type UnsupportedFormatError struct{ Format string }

func (e *UnsupportedFormatError) Error() string {
	return fmt.Sprintf("format %q not supported. Possible values: %s",
		e.Format, strings.Join(SupportedFormats, ", "))
}

// InvalidListenError reports that the configured Listen address could not
// be parsed as a host:port pair. The underlying net.SplitHostPort error is
// available via errors.Unwrap.
type InvalidListenError struct {
	Listen string
	Err    error
}

func (e *InvalidListenError) Error() string {
	return fmt.Sprintf("listen address %q is not a valid host:port: %v", e.Listen, e.Err)
}

func (e *InvalidListenError) Unwrap() error { return e.Err }
