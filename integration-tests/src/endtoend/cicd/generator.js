// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
var fs = require('fs')

module.exports = {
    writePServiceEnv: (participantID, serviceNameArray, imageVersion, netWorkName, envFileName, dockerComposeFile) => {
        serviceNameArray.forEach(async(serviceObj) => {

            let str = await '  ' + participantID + serviceObj.serviceUrl + ':\n' +
                '    image: ' + serviceObj.imageName + ':' + imageVersion + '\n' +
                '    env_file: \n' +
                '    - "' + envFileName + '"\n' +
                '    hostname: ' + participantID + serviceObj.serviceUrl + '\n' +
                '    container_name: ' + participantID + serviceObj.serviceUrl + '\n' +
                '    ports: \n'
            if (serviceObj.serviceEPort) {
                str += '    - "' + serviceObj.serviceEPort + '"\n'
            }
            if (serviceObj.serviceIPort) {
                str += '    - "' + serviceObj.serviceIPort + '"\n'
            }
            if (serviceObj.serviceName == 'payment-service' || serviceObj.serviceName == 'send-service') {
                str += await '    restart: unless-stopped \n' +
                    '    depends_on: \n' +
                    '    - ww-pr\n' +
                    '    networks: \n' +
                    '    - ' + netWorkName + '\n' +
                    '    environment: \n' +
                    '    - SERVICE_NAME=' + serviceObj.serviceName + '\n' +
                    '    - HOME_DOMAIN_NAME=' + participantID + '\n' +
                    '    - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}\n' +
                    '    - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}\n' +
                    '    - WW_JWT_PEPPER_OBJ=${WW_JWT_PEPPER_OBJ}\n' +
                    '    - FIREBASE_CREDENTIALS=${FIREBASE_CREDENTIALS}\n\n'

            } else {
                str += await '    restart: unless-stopped \n' +
                    '    networks: \n' +
                    '    - ' + netWorkName + '\n' +
                    '    environment: \n' +
                    '    - SERVICE_NAME=' + serviceObj.serviceName + '\n' +
                    '    - HOME_DOMAIN_NAME=' + participantID + '\n' +
                    '    - AWS_ACCESS_KEY_ID=${AWS_ACCESS_KEY_ID}\n' +
                    '    - AWS_SECRET_ACCESS_KEY=${AWS_SECRET_ACCESS_KEY}\n' +
                    '    - WW_JWT_PEPPER_OBJ=${WW_JWT_PEPPER_OBJ}\n' +
                    '    - FIREBASE_CREDENTIALS=${FIREBASE_CREDENTIALS}\n\n'

            }

            await fs.appendFile(dockerComposeFile, str, function(err) {
                if (err) throw err;
                // console.log(str);
                console.log('finished write service ' + participantID + ' - ' + serviceObj.serviceName + ' setting to docker-compose');

            });
        })
    },
    writeBasic: (filename) => {
        let str = "version: '3.5'\r\nnetworks:\r\n  wwcicdnet:\r\n    external: \r\n      name: wwcicdnet\r\nservices:\n"
        fs.writeFile(filename, str, function(err) {
            if (err) throw err;
        })
    }


}