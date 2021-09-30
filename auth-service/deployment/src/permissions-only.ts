// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// PURPOSE: The node application builds the permissions.go file only (faster-build time) 
// for the micro-services instead of having to run the full build process via index.ts
// used by gererateMiddlewarePermissions() in travis-microservices.sh 

import { SingleServiceDeployment } from './single-service';

class Main {

    constructor() {
        this.init();
    }

    private async init() {

        console.log('Generating permissions.go for middleware...');

        try {

            // build files for individual single services deployment:
            const permissionsYamlPath = './authorization/middleware/permissions/archive/permissions.yml';
            const permissionsGoPath = './authorization/middleware/permissions/permissions.go';
            
            const s = new SingleServiceDeployment()
            await s.replacePermissions(permissionsYamlPath, permissionsGoPath);

            console.log('Success, permissions.go generated!');

            process.exit(0);

        } catch (error) {
            console.error(error);
            process.exit(1);
        }

    }

}
new Main();