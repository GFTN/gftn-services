// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package stellarconfig

import "errors"


var (
    MOCK_DOMAIN = "mock.rdfi.gftn.io"
)

type MockStellarConfigFinder struct {



}



func (MockStellarConfigFinder) GetStellarConfigForDomain(domain string) (StellarConfig, error) {

    config := StellarConfig{}
    if domain == MOCK_DOMAIN {
        config.FederationServiceURL = "mock.federation.rdfi.gftn.io"
        return config, nil
    }

    return config, errors.New("Federation Server Could Not Be Found")

}
