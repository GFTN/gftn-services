// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// application
import { Config } from './config';
import { Environment } from './environment';
import { AuthMiddleware } from './middleware/auth.middleware';
// import { exec } from 'child_process';

// tsoa controllers 
// ============ Controllers =============
import './controllers/helper.controller';
import './controllers/ibmid.controller';
import './controllers/permissions.controller';
import './controllers/totp-two-factor.controller';
// institution controller - not being used at the moment wip
import './controllers/institution.controller';
import './controllers/jwt.controller';
// ============ Controllers =============

// register tsoa routes
import { RegisterRoutes } from './routes';

// 3rd party
// import { isUndefined } from 'lodash';

class Main {

    public config: Config;

    async init() {
        
        // // To see path on deployment:
        // const cmd = exec('pwd ; ls');
        // cmd.stdout.once('data', (data) => {
        //     console.log('current directory Location: ', data);
        // });

        // set global application layer env vars
        await new Environment().init();

        // set configurations
        this.config = new Config();
        await this.config.init();

        // register tsoa routes with express
        new AuthMiddleware(this.config.app);
        RegisterRoutes(this.config.app as any);

        // process.on('unhandledRejection', (reason, p) => {
        //     console.log('Unhandled Rejection at: Promise', p, 'reason:', reason);
        //     // application specific logging, throwing an error, or other logic here
        // });
    }

}

// set https endpoints for firebase functions:
//export const api = new Main().app.start();

// Start application as nodejs application:
export const api = async function () {
    const main = new Main()
    await main.init();
    return main.config.start();
}();
