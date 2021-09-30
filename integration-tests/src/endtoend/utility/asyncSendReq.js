// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const should = require('should');
const request = require('request')
const logRequest = require('./logRequest')
module.exports = (options) => {
    return new Promise((resolve, rejection) => {
        request(options, function(err, res, body) {
            if (err) {
                console.log(err);
                rejection(err)
            } else {
                logRequest(res, options)
                should(res.statusCode).be.exactly(200, "response data: " + JSON.stringify(res.body) + "\n request data : " + JSON.stringify(options))
                resolve(res)
            }
        });
    })
}