package logusgcloud

import (
	"cloud.google.com/go/logging"
	"github.com/strongo/logus"
)

func getGCSeverity(severity logus.Severity) logging.Severity {
	switch severity {
	case logus.SeverityDebug:
		return logging.Debug
	case logus.SeverityDefault:
		return logging.Default
	case logus.SeverityInfo:
		return logging.Info
	case logus.SeverityNotice:
		return logging.Notice
	case logus.SeverityWarning:
		return logging.Warning
	case logus.SeverityError:
		return logging.Error
	case logus.SeverityCritical:
		return logging.Critical
	case logus.SeverityAlert:
		return logging.Alert
	default:
		return logging.Default
	}
}
