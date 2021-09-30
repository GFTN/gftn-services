// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
var fs = require('fs');
var generator = require('./generator')
let imageVersion = 'latest'
let network = 'wwcicdnet'
let newENVPath = process.env.CICD_PATH + '/file/worldwireServices/configMap.env'
let newDCPath = process.env.CICD_PATH + '/file/worldwireServices/docker-compose.yml'

async function createFile(newENVPath, newDCPath, version, network) {


    // Write docker-compose Infrastructure
    await generator.writeBasic(newDCPath)
    console.log(process.env.CICD_PATH + '/serviceObj/');

    let files = fs.readdirSync(process.env.CICD_PATH + '/serviceObj/');
    files.forEach(file => {
        let serviceObj = require('./serviceObj/' + file)

        // Grnetate Global services to docker-compose file
        // Grnetate participant services
        generator.writePServiceEnv(serviceObj.participantID, serviceObj.serviceObjArray, version, network, "./configMap.env", newDCPath)
        if (serviceObj.callbackObj != undefined) {
            // Grnetate  Client callback services for participant1, participant2 to docker-compose file
            generator.writeCallback(serviceObj.participantID, serviceObj.callbackObj, version, network, newDCPath)
        }
        if (serviceObj.rdoClientObj != undefined) {
            // Grnetate  RDO-client services for participant1, participant2 to docker-compose file
            generator.writeRDOClient(serviceObj.participantID, serviceObj.rdoClientObj, version, network, newDCPath)
        }

    });

}


// Running async function to execute automated generate docker-compose file
createFile(newENVPath, newDCPath, imageVersion, network)