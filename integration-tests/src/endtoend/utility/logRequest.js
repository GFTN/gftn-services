// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const environment = require('../environment/env')

const fs = require('fs');
let islog = environment.ENV_KEY_ISLOG

module.exports = (res, options) => {
    if (islog) {
        if (res.body) {
            console.log("request data : \n" + JSON.stringify(options) + "\nresponse data : \n" + JSON.stringify(res.body));
        } else {

            console.log("request data : \n" + JSON.stringify(options) + "\nresponse data : \n" + JSON.stringify(res));
        }

        let today = new Date();
        fs.appendFile('./file/requestLog.txt', '\n' + today + '\n' + "request data : \n" + JSON.stringify(options) + "\nresponse data : \n" + JSON.stringify(res.body) + '\n', function(err) {
            if (err) throw err;
            // console.log('Saved!');
        });
    }

}