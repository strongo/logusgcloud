package logusgcloud

import (
	"cloud.google.com/go/logging"
	"github.com/strongo/logus"
	"testing"
)

func Test_getGCSeverity(t *testing.T) {
	type args struct {
		severity logus.Severity
	}
	tests := []struct {
		name string
		args args
		want logging.Severity
	}{
		{"debug", args{logus.SeverityDebug}, logging.Debug},
		{"default", args{logus.SeverityDefault}, logging.Default},
		{"info", args{logus.SeverityInfo}, logging.Info},
		{"notice", args{logus.SeverityNotice}, logging.Notice},
		{"warning", args{logus.SeverityWarning}, logging.Warning},
		{"error", args{logus.SeverityError}, logging.Error},
		{"critical", args{logus.SeverityCritical}, logging.Critical},
		{"alert", args{logus.SeverityAlert}, logging.Alert},
		{"unknown", args{logus.Severity(100)}, logging.Default},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getGCSeverity(tt.args.severity); got != tt.want {
				t.Errorf("getGCSeverity() = %v, want %v", got, tt.want)
			}
		})
	}
}
