package sdk

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseVersionRange(t *testing.T) {
	tests := []struct {
		name        string
		giveRange   string
		giveVersion string
		want        bool
	}{
		{
			name:        "no upper limit",
			giveRange:   "[0.9.16,)",
			giveVersion: "1.5.16",
			want:        true,
		},
		{
			name:        "no lower limit",
			giveRange:   "(,1.0.0]",
			giveVersion: "0.9.16",
			want:        true,
		},
		{
			name:        "lower inclusive",
			giveRange:   "[1.0.0,1.1.0]",
			giveVersion: "1.0.0",
			want:        true,
		},
		{
			name:        "upper inclusive",
			giveRange:   "(1.0.0,1.1.0]",
			giveVersion: "1.1.0",
			want:        true,
		},
		{
			name:        "lower exclusive",
			giveRange:   "(1.0.0,1.1.0]",
			giveVersion: "1.0.0",
			want:        false,
		},
		{
			name:        "upper exclusive",
			giveRange:   "[1.0.0,1.1.0)",
			giveVersion: "1.1.0",
			want:        false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseVersionRange(tt.giveRange)
			if err != nil {
				t.Errorf("ParseVersionRange() error = %v", err)
				return
			}
			assert.Equal(t, tt.want, got.matches(FixVersion(tt.giveVersion)))
		})
	}
}
