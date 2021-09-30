// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
package cryptoservice

import (
	"net/http"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/GFTN/gftn-services/gftn-models/model"
	comn "github.com/GFTN/gftn-services/utility/common"
)

func TestRequestSigning(t *testing.T) {
	part := model.Participant{}
	URL := "http://localhost:8888"
	part.URLCallback = &URL

	csc := Client{
		HTTP: &http.Client{Timeout: time.Second * 10},
	}
	signedXdr, err := csc.RequestSigning(
		"AAAAACIKcSda2GY1UmuKYyRF2uvTJPI6uhi1tYQ/MzZktwAQAAAAyAANhhcAAAAeAAAAAQAAAAAAAAAAAAAAAGIFUsIAAAAAAAAAAgAAAAEAAAAAN56OXjAeFHiGiWIiJUAocnJK3tU6wm3JxUfakiTHSNkAAAABAAAAAF49wJinSedEzsd5aWgwQWSQs2akIOI9+A8HnQemh+B6AAAAAlNHRERPAAAAAAAAAAAAAABePcCYp0nnRM7HeWloMEFkkLNmpCDiPfgPB50HpofgegAAAAAAmJaAAAAAAQAAAABePcCYp0nnRM7HeWloMEFkkLNmpCDiPfgPB50HpofgegAAAAEAAAAAN56OXjAeFHiGiWIiJUAocnJK3tU6wm3JxUfakiTHSNkAAAACVEhCRE8AAAAAAAAAAAAAADeejl4wHhR4holiIiVAKHJySt7VOsJtycVH2pIkx0jZAAAAAABcQ4oAAAAAAAAAAA==",
		"test",
		"4NJ2Wxj7a+sGUGgGnapsXZr8nS82886Ged5TLXwztyJHgwr87qdt5hR4nftZdxnn8oG0CbUYymaLnCnTMa51DA==",
		comn.ISSUING,
		part,
	)
	Convey("Successful get caller identity", t, func() {
		So(err, ShouldBeNil)
		So(signedXdr, ShouldNotBeNil)
		LOGGER.Info(signedXdr)
	})
}
