// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package utility

import (
	"github.com/op/go-logging"
	"os"
)

func ExitOnErr(LOGGER *logging.Logger, err error, errorMsg string) {

	if err != nil {
		LOGGER.Errorf(errorMsg+":  %v", err)
		os.Exit(1)
	}

}
