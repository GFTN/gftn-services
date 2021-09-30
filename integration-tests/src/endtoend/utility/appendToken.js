// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const environment = require('../environment/env')

let istoken = environment.ENV_KEY_ISTOKEN

module.exports = (options, id, WLorQuote) => {
    if (environment.ENV_KEY_ISTOKEN == "true" || WLorQuote == true) {
        switch (id) {
            case "ENV_KEY_PARTICIPANT_1_ID":
                // options["headers"].Authorization = 'Bearer ' + process.env["ibmsingapore1"]
                options["headers"].Authorization = 'Bearer ' + process.env.PARTICIPANT_1_JWT_TOKEN
                break;
            case "ENV_KEY_PARTICIPANT_2_ID":
                options["headers"].Authorization = 'Bearer ' + process.env.PARTICIPANT_2_JWT_TOKEN
                break;
            case "ENV_KEY_ANCHOR_ID":
                options["headers"].Authorization = 'Bearer ' + process.env.ANCHOR_JWT_TOKEN
                break;
            default:
                break;
        }
        return options
    } else { return (options) }

}