// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package stellarconfig

import (
    "github.com/go-resty/resty"
    "errors"
    "github.com/BurntSushi/toml"
    "os"
    "strings"
    "github.com/GFTN/gftn-services/utility/global-environment"
)

type RestStellarConfigFinder struct {

    StellarConfigRemaps map[string]string
    StellarConfigScheme string

}

func CreateRestStellarConfigFinder() (RestStellarConfigFinder, error) {
    finder := RestStellarConfigFinder{}

    remaps := os.Getenv(global_environment.ENV_KEY_STELLAR_CONFIG_REMAP)
    scheme := os.Getenv(global_environment.ENV_KEY_STELLAR_CONFIG_SCHEME)

    if remaps != "" {
        parts := strings.Split(remaps, "->")
        finder.StellarConfigRemaps[parts[0]] = parts[1]
    }

    if scheme != "" {
        finder.StellarConfigScheme = scheme
    } else {
        finder.StellarConfigScheme = "https"
    }

    return finder, nil
}


func (finder RestStellarConfigFinder) GetStellarConfigForDomain(domain string) (StellarConfig, error) {

    configResponderUrl := finder.createConfigResponderUrl(domain)
    response, err := resty.R().Get(configResponderUrl)

    if err != nil {
        LOGGER.Errorf("Error while getting stellar config for domain (%v):  %v", domain, err)
        return StellarConfig{}, err
    }

    if response.StatusCode() != 200 {
        LOGGER.Warningf("Status code (%v) was not 200", response.StatusCode())
        return StellarConfig{}, errors.New("Invalid status code")

    }

    var stellarConfig StellarConfig
    err = toml.Unmarshal(response.Body(), &stellarConfig)

    if err != nil {
        LOGGER.Errorf("Error while unmarshalling stellar config for domain (%v):  %v", domain, err)
        return StellarConfig{}, err
    }


    return stellarConfig, nil

}


func (finder RestStellarConfigFinder) createConfigResponderUrl(domain string) (string) {

    configResponderUrl := finder.StellarConfigScheme + "://"
    if finder.StellarConfigRemaps[domain] != "" {
        configResponderUrl = configResponderUrl + finder.StellarConfigRemaps[domain] + "/.well_known/stellar.toml"
    } else {
        configResponderUrl = configResponderUrl + domain + "/.well_known/stellar.toml"
    }

    return configResponderUrl

}
