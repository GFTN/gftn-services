// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
var PS = require("./awsParameter");
var SM = require("./awsSecret");
var Env = require("./env");
env_main();
function env_main() {
    process.env['HOME_DOMAIN_NAME'] = "p1.worldwire.io";
    process.env['SERVICE_NAME'] = "api-service";
    process.env['ENVIRONMENT_VERSION'] = "dev";
    process.env['ACCOUNT_STORAGE_LOCATION'] = "dev";
    process.env['ACCOUNT_SOURCE'] = "dev";
    Env.InitEnv();
}
function parameter_main() {
    var title = {
        environment: "dev",
        domain: "worldwire.io",
        service: "IBM",
        variable: "test9"
    };
    var oldContent = {
        value: "worldwire.io",
        description: "IBM"
    };
    /* await/async way(ensure that you add the async keyword at this function)
    let result: any
    try{
      result = await PS.getParameter(title)
    }catch(e){
      console.log(e)
    }
      console.log(result)
      console.log("hihih")
    */
    var newContent = {
        value: "newest worldwire.io",
        description: "IBM New"
    };
    PS.createParameter(title, oldContent)
        .then(function (res) {
        console.log(res);
        return PS.getParameter(title);
    })
        .then(function (res) {
        console.log(res);
        return PS.updateParameter(title, newContent);
    })
        .then(function (res) {
        console.log(res);
        return PS.getParameter(title);
    })
        .then(function (res) {
        console.log(res);
        return PS.removeParameter(title);
    })
        .then(function (res) {
        console.log(res);
        return PS.getParameter(title);
    })
        .catch(function (err) {
        console.log(err);
    });
}
function secret_main() {
    var title = {
        environment: "dev",
        domain: "worldwire.io",
        service: "IBM",
        variable: "test12"
    };
    var oldContent = {
        // key: "dev",
        // value: "worldwire.io",
        filePath: "/Users/your.user/go/src/github.ibm.com/gftn/world-wire-services/utility/aws/javascript/src/test.json",
        description: "IBM"
    };
    var newContent = {
        key: "this is",
        value: "new!",
        //filePath: "/Users/your.user/go/src/github.ibm.com/gftn/world-wire-services/utility/aws/javascript/src/test.json",
        description: "IBM"
    };
    SM.createSecret(title, oldContent)
        .then(function (res) {
        console.log(res);
        return SM.getSecret(title);
    })
        .then(function (res) {
        console.log(res);
        return SM.updateSecret(title, newContent);
    })
        .then(function (res) {
        console.log(res);
        return SM.getSecret(title);
    })
        .then(function (res) {
        console.log(res);
        return SM.removeSecret(title);
    })
        .then(function (res) {
        console.log(res);
        return SM.getSecret(title);
    })
        .catch(function (err) {
        console.log(err);
    });
}
//# sourceMappingURL=test.spec.js.map
