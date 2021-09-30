// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package logconfig

import (
	"io"
	"log"
	"os"

	"github.com/op/go-logging"
)

var format = logging.MustStringFormatter(
	`%{color:bold}[%{level:.8s}]%{color:reset} %{color}%{time:2006-01-02T15:04:05Z07:00} %{shortfile} %{shortfunc}%{color:reset} ▶ %{id:03x} %{message}`,
)

/* sets up logging directing it to the given log file */
func SetupLogging(serviceLogs string, LOGGER *logging.Logger) (*os.File, error) {
	LOGGER.Infof("Log File: %s", serviceLogs)
	f, err := os.OpenFile(serviceLogs, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v: %v", err, serviceLogs)
		return nil, err
	}
	logWriter := io.MultiWriter(f, os.Stdout)

	backend1 := logging.NewLogBackend(logWriter, "", 0)
	backend2Formatter := logging.NewBackendFormatter(backend1, format)
	logging.SetBackend(backend2Formatter)
	return f, nil
}
