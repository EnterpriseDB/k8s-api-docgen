// Package log handles logging
package log

import (
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
)

// Log is the logger to be used inside this package
var Log logr.Logger

func init() {
	zapLogger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	Log = zapr.NewLogger(zapLogger)
}
