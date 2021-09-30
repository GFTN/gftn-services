// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// PURPOSE: to create the google cloud endpoints deployment file 

import * as fs from 'fs';
import { cloneDeep, set, forEach, size } from 'lodash';
import * as YAML from 'yamljs';
import { microServices, IMicroServiceDefItem } from './micro-services-constants'

class Main {

    constructor() {
        //start to create deployment files
        this.init();
    }

    private async init() {

        console.log('Starting - Generating Endpoints Gcloud Deployments files...');

        // build files for individual micro-services deployment:
        await new MicroServicesDeployment('1.0.0').build();

        console.log('Success! GCloud Endpoints deployment files created.\r\nApplication shutting down...');

        // terminate application
        process.exit();

    }

}

class MicroServicesDeployment {

    private version: string;
    private microServices = microServices;

    constructor(
        // api version to use for gcloud endpoints
        version: string
    ) {
        this.version = version;
        this.build();
    }

    async build() {
        await this.createAppEngineEndpoints();
    }

    /**
     * Creates swagger definition to be consumed by gcloud endpoints
     *
     * @returns
     * @memberof MicroServicesDeployment
     */
    private async createAppEngineEndpoints() {

        // only resovle once loop is over
        let count = 0;

        return new Promise((resolve, reject) => {

            forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

                fs.readFile(microServiceDir + '/def/swagger.json', 'utf8', (err, _devSwagger) => {

                    const devSwagger = cloneDeep(JSON.parse(_devSwagger));

                    if (err) {
                        console.error('unable to read swagger.json: ' + microServiceDir, err)
                        resolve();
                    } else {

                        // update properties of swagger.json
                        set(devSwagger, 'host', val.swagger.host);
                        set(devSwagger, 'info.version', this.version);
                        set(devSwagger, 'info.title', val.swagger.title);

                        // At the top level of the file (not indented or nested), add an empty security directive to apply it to the entire API
                        // https://cloud.google.com/endpoints/docs/openapi/restricting-api-access-with-api-keys#restricting_access_to_specific_api_methods
                        set(devSwagger, 'security', []);

                        // turn off API key validation for a particular method even when you've restricted API access for the API
                        // https://cloud.google.com/endpoints/docs/openapi/restricting-api-access-with-api-keys#restricting_access_to_specific_api_methods
                        // The /callback endpoint only interacts with IBMId so it is not possible to pass along an api key
                        set(devSwagger, 'paths./auth/sso/callback.get.security', [])

                        // convert js obj in yaml
                        const productionYamlSwagger = YAML.stringify(devSwagger, 20);

                        fs.writeFile(microServiceDir + '/build/openapi-appengine.yaml', productionYamlSwagger, (err) => {

                            // success creating the app.yaml file
                            if (err) {
                                console.error('error creating openapi-appengine.yaml: ' + microServiceDir, err);
                            }

                            // increment and check if loop is completed
                            count = count + 1;
                            if (count = size(this.microServices)) {
                                resolve();
                            }

                        });

                    }

                });

            });

        });

    }

}

// start
new Main();