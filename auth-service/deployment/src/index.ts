// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// PURPOSE: to create a app.yaml (needed for deployment to GCP App Engine) 
// with appropriate secrets and env vars for targeted environment
// ie: prod(aka: worldwire.io) vs. dev (aka: next.worldwire.io)
// ENV vars are retrieved from the .vscode/launch.json for the targeted 
// environment. Since the launch.json is not pushed to github this prevents
// the accidental leaking of env vars 

import { exec } from 'child_process';
// import { MicroServicesDeployment } from './micro-services';
import { SingleServiceDeployment } from './single-service';

class Main {

    constructor() {
        //start to create deployment files
        this.init();
    }

    private async init() {

        console.log('Starting - Generating build files for deployment (GAE & Docker)...');

        const cmd = exec('pwd');
        cmd.stdout.once('data', (data) => {
            console.log('current directory Location: ', data);
        });

        try {

            // build files for individual micro-services deployment:
            // await new MicroServicesDeployment();

            // build files for individual single services deployment:
            await new SingleServiceDeployment().build();

            console.log('Success! Deployment files generated at ./build...');

            process.exit(0);

        } catch (error) {
            console.error(error);
            process.exit(1);
        }

    }

}

// start
new Main();