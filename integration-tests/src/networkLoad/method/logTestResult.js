// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
var fs = require('fs');
var util = require('util');
  // Or 'w' to truncate the file every time the process starts.
var logStdout = process.stdout;

module.exports = function (){
  var logFile = fs.createWriteStream('./logs/ConcurrentTestResult.txt', { flags: 'a' });
  fs.writeFile('./logs/ConcurrentTestResult.txt', '', (err) => {
    if (err) throw err;
  });
  this.Report = function(){
    console.log = function () {
        logFile.write(util.format.apply(null, arguments) + '\n');
        logStdout.write(util.format.apply(null, arguments) + '\n');
      }
      console.error = console.log;
  }
}
