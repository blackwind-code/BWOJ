package main

import (
	"io"
	"os"

	nested "github.com/antonfisher/nested-logrus-formatter"
	"github.com/blackwind-code/BWOJ/pkg/util"
	"github.com/sirupsen/logrus"
)

func init_logrus(log_fp, log_level, log_fmt string) {
	// Set Logrus level
	log_level_iota, e := logrus.ParseLevel(log_level)
	if e != nil {
		logrus.WithFields(logrus.Fields{
			"error": e,
		}).Panicf("Wrong argument: %v", log_level)
		os.Exit(1)
	}
	logrus.SetLevel(log_level_iota)

	// Set Logrus output
	var writer io.Writer
	// Default value is command line only
	writer = os.Stdout
	if log_fp != "" {
		log_fp = util.PathSanitize(log_fp)
		log_f, e := os.OpenFile(log_fp, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if e != nil {
			logrus.WithFields(logrus.Fields{
				"error": e,
			}).Fatalf("Unable to create log file: %v", log_fp)
		}
		writer = io.MultiWriter(log_f, os.Stdout)
	}
	logrus.SetOutput(writer)

	switch log_fmt {
	case "text":
		logrus.SetFormatter(&nested.Formatter{
			FieldsOrder:      []string{"error"},
			TimestampFormat:  "2006-01-02 15:04:05.000",
			HideKeys:         false,
			NoColors:         false,
			NoFieldsColors:   false,
			NoFieldsSpace:    false,
			ShowFullLevel:    false,
			NoUppercaseLevel: false,
			TrimMessages:     true,
			CallerFirst:      true,
		})
	case "json":
		logrus.SetFormatter(&logrus.JSONFormatter{
			TimestampFormat:   "2006-01-02 15:04:05.000",
			DisableTimestamp:  false,
			DisableHTMLEscape: true,
			PrettyPrint:       false,
		})
	default:
		logrus.Panicf("Wrong log format argument: %v", log_fmt)
		os.Exit(1)
	}
}
