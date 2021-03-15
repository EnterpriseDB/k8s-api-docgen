/*
Copyright 2021 EnterpriseDB Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

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
