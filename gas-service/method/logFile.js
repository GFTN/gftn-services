// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
var fs = require('fs');
var util = require('util');
var environment = require('../environment/env')
    // Or 'w' to truncate the file every time the process starts.

var logFile = fs.createWriteStream(process.env[environment.ENV_KEY_SERVICE_LOG_FILE], {
    flags: 'a'
});
var logStdout = process.stdout;

module.exports = function() {

    console.log = function() {
        logFile.write(util.format.apply(null, arguments) + '\n');
        logStdout.write(util.format.apply(null, arguments) + '\n');
    }
    console.error = console.log;
}