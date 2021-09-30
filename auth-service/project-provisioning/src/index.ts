// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { exec } from 'child_process';
import * as fs from 'fs';
import * as _ from 'lodash';
import * as readline from 'readline';
import { Helpers } from './helpers';

// Test ReadJson()
// import { ReadJson } from './read_val';
// new ReadJson();

class main {

    // private ignoreErrorStringsIncluding = [];

    private adminUid = 'rootsuperuser';
    private orgId = '137052568356';

    // user input project name:
    private envShortName: string;
    private envLongName: string;
    // private envShortName = 'pen1';
    // private envLongName = 'Pentesting1'

    // newly created gcloud project id:
    private gcloudProjectName: string;
    // private gcloudProjectName = 'pen1-260919';

    // seed initial super admin email address:
    private superAdminEmail: string;
    // private superAdminEmail = 'chaseo@ibm.com';

    // iam for deployment:
    private deployServiceAccountName: string;
    // private deployServiceAccountName = 'gae-deploy-968';
    
    // iam for firebase: 
    // private firebaseAdminSDKServiceAccountName: string;
    // private firebaseAdminSDKServiceAccountName = 'firebase-adminsdk-1828';
    // private gcloudRegion='us-central';

    private credKey: string;
    private credIv: string;

    private credDir: string;
    private credVersion: string;
    // private credDir = '.credentials-v22/raw/pen1'
    // private credVersion = '22';

    // domain that hosts the portal eg: worldwire-dev.io:
    private portalDomain: string;
    // private portalDomain = 'worldwire-pen.io';

    constructor() {
        this.init();
    }

    async justOne() {

    }

    async init() {

        try {

            // create a random  gae deploy name by adding 4 digits to the end
            this.deployServiceAccountName = 'gae-deploy-' + Math.round(Math.random() * 10000);
            // this.firebaseAdminSDKServiceAccountName = 'firebase-adminsdk-' + Math.round(Math.random() * 10000);

            // check if you are logged into GCLoud locally and if firebase and GCloud exist
            await this.checkDependencies();
            await this.requireCliLogins();

            // prompt to name the gcloud + firebase project
            await this.setNewEnvName();
            await this.setNewEnvLongName();

            // prompt super user for email address
            await this.setInitialAdminEmail();

            // prompt for portal domain
            await this.setPortalDomain();

            // prompt for setting current credentials version
            await this.setCredDir();

            // create gcloud project
            await this.createGCloudProject();

            // add firebase project to gcloud project
            await this.addFirebase();

            // create first super admin user the firebase authentication user directory
            await this.addSuperUserToFirebaseAuth();

            // add new firebase project name to gftn-web/.firebaserc
            await this.addFirebaseUseConfig();

            // deploy firebase functions for users (IMPORTANT: must be before adding the user to the database)
            await this.deployFirebaseFunctions();

            // create first super admin user in the firebase database    
            await this.addSuperUserToFirebaseDatabase();

            // apply bolt security rules
            await this.applyFirebaseBoltRules();

            // create credentials for new project 
            await this.createSecretsEnvs();
            await this.createServiceAccounts();
            await this.addServiceAccountRoles();

            // enable billing for gcloud account
            await this.enableBilling();

            // setup callback url for IBMid
            await this.setupIBMIdCallback();

            // deploy portal
            await this.createBuildConfig();

            // rotate credentials 
            // (necessary because deployment of the authentication 
            // service will require updated encrypted env vars and secrets)
            await this.rotateCreds();
            await this.exportCredDecryptionKeys();

            // deploy auth-service, portal, functions, database rules
            await this.deployAuthService();

            // config DNS and Firewall
            // must come before configureAuthServiceDNS() because verification of the domain occurs first with configurePortalDNS()
            await this.configurePortalDNS();
            await this.requireCliLogins(); // require admin login again
            await this.configureAuthServiceDNS();
            await this.configureFirewall();

            console.info('gcloud provisioning completed successfully...')

        } catch (error) {
            console.error(error);
            process.exit(1);
        }

    }

    // async execScriptExample() {

    //     const methodName = 'execScriptExample';

    //     return await new Promise((resolve, reject) => {

    //         // cmd
    //         let script = 'SOME SHELL SCRIPT HERE'

    //         // execute shell
    //         const cmd = exec(script);

    //         // output
    //         cmd.stdout.on('data', (outputText: string) => {
    //             console.log(methodName + ' stdout: ', outputText);
    //             // const start = 'START TEXT HERE';
    //             // const end = 'END TEXT HERE';
    //             // // const middleText = outputText.match(new RegExp(start + "(.*)" + end))[0];
    //             // const const middleText = outputText.split(start)[1].split(end)[0].trim()
    //             // console.log(methodName + ' middleText: ', middleText);
    //         });

    //         // count errors out by stderr (needed for scripts separated by ';')
    //         let _code = 0;

    //         // log out error info
    //         cmd.stderr.on('data', (data) => {
    //             console.error(methodName + ' stderr: ', { data: data as string });
    //             let ignore = false;
    //             for (let i = 0; i < this.ignoreErrorStringsIncluding.length; i++) {
    //                 if (data.includes(this.ignoreErrorStringsIncluding[i])) {
    //                     ignore = true;
    //                 }
    //             }
    //             if (!ignore) {
    //                 _code = _code + 1;
    //             }
    //         });

    //         // script competed
    //         cmd.once('exit', (code: number, signal: string) => {
    //             _code = code + _code;
    //             console.log(methodName + ' exit ', { code: _code, signal: signal });
    //             if (_code === 0) {
    //                 resolve();
    //             } else {
    //                 reject(new Error(methodName + ' failed'));
    //             }
    //         });

    //     });

    // }

    async rotateCreds() {

        const methodName = 'rotateCreds';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = 'tsc ./secret-mgmt/src/index.ts ; node ./secret-mgmt/src/index.js'

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {

                console.log(methodName + ' stdout: ', outputText);
                
                // get the key 
                try {                
                    const start = `CICD_CRED_KEY_V${Number(this.credVersion) + 1}=`;
                    const end = `CICD_CRED_IV_V${Number(this.credVersion) + 1}=`;
                    const middleText = outputText.split(start)[1].split(end)[0].trim()
                    console.log(methodName + ' CICD_CRED_KEY: ', middleText);
                    this.credKey = middleText;
                } catch (error) {
                    
                }
                
                // get the iv 
                try {                
                    const start = `CICD_CRED_IV_V${Number(this.credVersion) + 1}=`;
                    const end = `# Run decryption using: `;
                    const middleText = outputText.split(start)[1].split(end)[0].trim()
                    console.log(methodName + ' CICD_CRED_IV: ', middleText);
                    this.credIv = middleText;
                } catch (error) {
                    
                }
                
            });

            // count errors out by stderr (needed for scripts separated by ';')
            let _code = 0;

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data as string });
                // let ignore = false;
                // for (let i = 0; i < this.ignoreErrorStringsIncluding.length; i++) {
                //     if (data.includes(this.ignoreErrorStringsIncluding[i])) {
                //         ignore = true;
                //     }
                // }
                // if (!ignore) {
                //     _code = _code + 1;
                // }
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                _code = code + _code;
                console.log(methodName + ' exit ', { code: _code, signal: signal });
                if (_code === 0) {
                    resolve();
                } else {
                    reject(new Error(methodName + ' failed'));
                }
            });

        });

    }    

    /**
     * sets environment variables locally to decrypt credentials needed to deploy authentication service (emulating what travis would usually do)
     *
     * @returns
     * @memberof main
     */
    async exportCredDecryptionKeys() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nAdd the the following envs to your ~/.bash_profile (open in vscode using ${this.fmtN('FgWhite')} $ code ~/.bash_profile ${this.fmtN('FgBlue')}). ${this.fmtN('FgYellow')}\nexport CICD_CRED_KEY_V${Number(this.credVersion)+1}="${this.credKey}"\nexport CICD_CRED_IV_V${Number(this.credVersion)+1}="${this.credIv}"${this.fmtN('FgBlue')} \n\nThen run ${this.fmtN('FgWhite')}$ source ~/.bash_profile${this.fmtN('FgBlue')} to persist to the machines env vars\n\nThen add the rotated key and iv to travis env vars at https://travis.ibm.com/gftn/gftn-services/settings \n\nPress enter to continue.`, ans => {
            rl.close();
            resolve(ans);
        }))

        // console.log(`User indicated that that billing is enabled for  ${this.createGCloudProject}`);

    }

    async requireCliLogins() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nLogin by running the following in another terminal: \n\n$ gcloud auth application-default login ; firebase login \n\n Press enter to continue.`, ans => {
            rl.close();
            resolve(ans);
        }))

        // console.log(`User indicated that that billing is enabled for  ${this.createGCloudProject}`);

    }

    async configurePortalDNS() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nCONFIGURE PORTAL HOSTING:\nGo to https://console.firebase.google.com/u/1/project/${this.gcloudProjectName}/hosting/main and follow the instructions to setup DNS in AWS Route 53. Then press enter.`, ans => {
            rl.close();
            resolve(ans);
        }))

        // console.log(`User indicated that that billing is enabled for  ${this.createGCloudProject}`);

    }

    async configureAuthServiceDNS() {

        // // ask user to name project from cmd line
        // const rl = readline.createInterface({
        //     input: process.stdin,
        //     output: process.stdout,
        //     terminal: false
        // });

        // await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nCONFIGURE AUTH-SERVICE HOSTING:\nGo to https://www.google.com/webmasters/verification/verification?domain=auth.${this.portalDomain} then select the "Other" provider, and follow the instructions to create the dns records in aws route 53. Then press enter.`, ans => {
        //     rl.close();
        //     resolve(ans);
        // }))

        const methodName = 'configureAuthServiceDNS';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = `gcloud app domain-mappings create auth.${this.portalDomain} --project="${this.gcloudProjectName}"` 
            console.info('running: '+ script)

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', async (code: number, signal: string) => {
                console.log(methodName + ' exit ', { code: code, signal: signal });

                // ask user to name project from cmd line                
                const rl = readline.createInterface({
                    input: process.stdin,
                    output: process.stdout,
                    terminal: false
                });

                await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nCONFIGURE AUTH-SERVICE HOSTING:\nGo to route53 in aws. Add a CNAME record for 'auth.${this.portalDomain}' and value 'ghs.googlehosted.com.'. It may take 24hrs to provisiion a certificate. Then press enter.`, ans => {
                    rl.close();
                    resolve(ans);
                }))

                if (code === 0) {
                    resolve();
                } else {
                    reject();
                    process.exit(1);
                }
            });

        });

    }

    async configureFirewall() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\n Enable the compute api by clicking enable via https://console.developers.google.com/apis/api/compute.googleapis.com/overview?project=${this.gcloudProjectName}&authuser=1 \n\n Press enter to continue.`, ans => {
            rl.close();
            resolve(ans);
        }))

        const methodName = 'configureFirewall';

        return await new Promise((resolve, reject) => {

            const deleteRules = [
                'default-allow-rdp', // removes tcp:3389
                'default-allow-ssh' // removes tcp:22
            ]

            let script = '';
            for (let i = 0; i < deleteRules.length; i++) {
                script = script + `gcloud compute firewall-rules delete ${deleteRules[i]} --project="${this.gcloudProjectName}" --quiet`
                if (i < deleteRules.length - 1) {
                    // add separator
                    script = script + ' ; '
                }
            }

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    reject();
                    process.exit(1);
                }
            });

        });

    }

    async setupIBMIdCallback() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        // await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nAdd callback url "https://${this.gcloudProjectName}.appspot.com/sso/callback" to the list of IBMId Callbacks here: https://w3.innovate.ibm.com/tools/sso/application/list.html. Press enter to continue.`, ans => {
        await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nAdd callback url "auth.${this.portalDomain}" to the list of IBMId Callbacks here: https://w3.innovate.ibm.com/tools/sso/application/list.html (or try https://w3.ibm.com/tools/sso/application/list.html). Press enter to continue.`, ans => {
            rl.close();
            resolve(ans);
        }))

        // console.log(`User indicated that that billing is enabled for  ${this.createGCloudProject}`);

    }

    async enableBilling() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nEnable billing for this new gcloud account by navigating to  https://console.developers.google.com/apis/api/cloudbuild.googleapis.com/overview?project=${this.gcloudProjectName} . This is required before the auth-service can be deployed to google app-engine (gae). Press enter after you have enabled billing.`, ans => {
            rl.close();
            resolve(ans);
        }))

        // console.log(`User indicated that that billing is enabled for  ${this.createGCloudProject}`);

    }

    async createServiceAccounts() {

        const methodName = 'createServiceAccounts';

        return await new Promise((resolve, reject) => {

            // cmd
            // let script = `gcloud iam service-accounts create ${this.deployServiceAccountName} --project=${this.gcloudProjectName} ; gcloud iam service-accounts create ${this.firebaseAdminSDKServiceAccountName} --project=${this.gcloudProjectName}`
            let script = `gcloud iam service-accounts create ${this.deployServiceAccountName} --project=${this.gcloudProjectName} --description "Deploys the authentication service" --display-name "deploy-${this.gcloudProjectName}"`

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {

                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    reject();
                    process.exit(1);
                }
            });

        });

    }

    async createSecretsEnvs() {

        // ask user to provide firebase admin sdk credential from cmd line
        const rl2 = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        // get the firebase admin sdk credential output
        await new Promise(resolve => rl2.question(this.fmtN('FgBlue') + `\n\nRun the following in another terminal window \n\n$ cp -rf .credentials-v19/raw/dev/. ${this.credDir} ; rm ${this.credDir}/portal.txt ; cd ../gftn-web/ ; firebase use ${this.envShortName} ; firebase apps:create WEB portal ;  firebase apps:sdkconfig -o ../auth-service/${this.credDir}/portal.txt \n\nPress enter to continue.`, ans => {
            rl2.close();
            resolve(ans);
        }));

        // get default dev envs
        const envJsonTxt: string = await fs.readFileSync(`${this.credDir}/img/env.json`, "utf8");
        let envJsonObj = JSON.parse(envJsonTxt);

        // get default dev secrets
        const secretJsonTxt: string = await fs.readFileSync(`${this.credDir}/img/secret.json`, "utf8");
        let secretJsonObj = JSON.parse(secretJsonTxt);

        // update secret and env variables for new environment
        envJsonObj.totp_label = `World Wire - ${this.envShortName}`;
        envJsonObj.site_root = `https://${this.portalDomain}`;
        // envJsonObj.app_root = `https://${this.gcloudProjectName}.appspot.com`;
        envJsonObj.app_root = `https://auth.${this.portalDomain}`;
        envJsonObj.gae_service = `${this.gcloudProjectName}`;
        envJsonObj.fb_database_url = `https://${this.gcloudProjectName}.firebaseio.com`;

        // ask user to provide firebase admin sdk credential from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        // get the firebase admin sdk credential output
        const cred: string = await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nSet the service account. Navigate to https://console.firebase.google.com/u/1/project/${this.gcloudProjectName}/settings/serviceaccounts/adminsdk . Click \" Generate new private key\ at the bottom of the page to download the json credential. Remove all new lines and whitespace from the service account json credential and paste here: "`, ans => {
            rl.close();
            resolve(ans);
        }));

        const helpers = new Helpers();
        secretJsonObj.fb_admin = await helpers.encodeFbCred(JSON.parse(cred));
        secretJsonObj.ww_jwt_pepper_encoded = Buffer.from(JSON.stringify(helpers.getPepperObj())).toString('base64');
        secretJsonObj.passport_secret = helpers.randStr(32);

        // overwrite env.json
        await fs.writeFileSync(`${this.credDir}/img/env.json`, JSON.stringify(envJsonObj));

        // overwrite secret.json
        await fs.writeFileSync(`${this.credDir}/img/secret.json`, JSON.stringify(secretJsonObj));

        return;

    }

    async addServiceAccountRoles() {

        const methodName = 'addServiceAccountRoles';

        return await new Promise((resolve, reject) => {

            // create iam policy scripts (must be added one at a time):
            // gcloud app engine roles
            const gaeDeploymentRoles = [
                'appengine.appViewer',
                'appengine.codeViewer',
                'appengine.deployer',
                'appengine.serviceAdmin',
                'cloudbuild.builds.editor',
                'compute.storageAdmin',
                'storage.admin'
            ];

            // // firebase admin sdk roles
            // const firebaseAdminSdkRoles = [
            //     'editor'
            // ];

            // add iam policies one by one
            let deploymentIamPolicies = '';
            for (let i = 0; i < gaeDeploymentRoles.length; i++) {
                deploymentIamPolicies = deploymentIamPolicies + `gcloud projects add-iam-policy-binding ${this.gcloudProjectName} --member=serviceAccount:${this.deployServiceAccountName}@${this.gcloudProjectName}.iam.gserviceaccount.com --role=roles/${gaeDeploymentRoles[i]} ; `
            }

            // let firebaseAdminSdkIamPolicies = '';
            // for (let i = 0; i < firebaseAdminSdkRoles.length; i++) {
            //     firebaseAdminSdkIamPolicies = firebaseAdminSdkIamPolicies + `gcloud projects add-iam-policy-binding ${this.gcloudProjectName} --member=serviceAccount:${this.firebaseAdminSDKServiceAccountName}@${this.gcloudProjectName}.iam.gserviceaccount.com --role=roles/${firebaseAdminSdkRoles[i]} ; `
            // }

            // output json service account credentials to ./credentials/{taget env}:
            // deployment credential
            const generateDeploymentJsonKey = `gcloud iam service-accounts keys create ${this.credDir}/deploy/deploy.json --iam-account ${this.deployServiceAccountName}@${this.gcloudProjectName}.iam.gserviceaccount.com --project=${this.gcloudProjectName}`;
            
            // // firebase Admin sdk credential
            // const generateFirebaseAdminSDKJsonKey = `gcloud iam service-accounts keys create ${this.credDir}/adminsdk.json --iam-account ${this.firebaseAdminSDKServiceAccountName}@${this.gcloudProjectName}.iam.gserviceaccount.com --project=${this.gcloudProjectName} ; gcloud app create --project=${this.gcloudProjectName} --region ${this.gcloudRegion}`;

            // cmd
            let script = deploymentIamPolicies + 
            // firebaseAdminSdkIamPolicies + 
            generateDeploymentJsonKey 
            // + ' ; ' + generateFirebaseAdminSDKJsonKey
            ;

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {

                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    reject();
                    process.exit(1);
                }
            });

        });

    }

    async deployAuthService() {
         
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        await new Promise(resolve => rl.question(this.fmtN('FgBlue') + `\n\nDeploy the authentication service to gcloud. If this is a completely new environment you may have to: \n1. Update the ./deployment/deploy.sh , see deploy.sh line 254 to add the travis specific branch and environment name. \n2. Update ./secret-mgmt/src/index.ts line 28 to add the specific env. \n3. Enable cloud build deployment, click enable via https://console.developers.google.com/apis/api/cloudbuild.googleapis.com/overview?project=${this.gcloudProjectName} \n4. Run the following ${this.fmtN('FgWhite')}$ export TRAVIS_BRANCH=${this.envShortName}-gftn ; export TRAVIS_PULL_REQUEST=false ; bash -x ./deployment/deploy.sh --env ${this.envShortName} \n\n${this.fmtN('FgBlue')}Press enter to continue.`, ans => {
            rl.close();
            resolve(ans);
        }))

        // const methodName = 'deployAuthService';

        // return await new Promise((resolve, reject) => {

        //     // cmd
        //     let script = `export TRAVIS_BRANCH=${this.envShortName}-gftn ; export TRAVIS_PULL_REQUEST=false ; bash -x ./deployment/deploy.sh --env ${this.envShortName}`

        //     // execute shell
        //     const cmd = exec(script);

        //     // output
        //     cmd.stdout.on('data', (outputText: string) => {
        //         console.log(methodName + ' stdout: ', outputText);
        //     });

        //     // log out error info
        //     cmd.stderr.on('data', (data) => {
        //         console.error(methodName + ' stderr: ', { data: data });
        //     });

        //     // script competed
        //     cmd.once('exit', (code: number, signal: string) => {
        //         // console.log(methodName + ' exit ', { code: code, signal: signal });
        //         if (code === 0) {
        //             resolve();
        //         } else {
        //             reject();
        //             process.exit(1);
        //         }
        //     });

        // });

    }

    async createBuildConfig() {


        //  add the environment.{target}.ts:
        const environmentJsonPath = `../gftn-web/src/environments/environment.${this.envShortName}.ts`;
        const portalConfigTxt: string = await fs.readFileSync(`${this.credDir}/portal.txt`, "utf8");

        //         const environmentJsonObj = `
        // export const environment = {
        //     production: true,
        //     firebase: ${portalConfigTxt},
        //     apiRootUrl: 'https://${this.gcloudProjectName}.appspot.com',
        //     supported_env: {
        //         text: '${this.envLongName}',
        //         val: '${this.envShortName}',
        //         envApiRoot: '${this.portalDomain}/local/api',
        //         envGlobalRoot: '${this.portalDomain}/global',
        //     },
        //     inactivityTimeout: 15
        //     };

        // `;

        const environmentJsonObj =
            `export const environment = {
    production: true,
    firebase: ${portalConfigTxt},
    apiRootUrl: 'https://auth.${this.portalDomain}'
    supported_env: {
        text: '${this.envLongName}',
        val: '${this.envShortName}',
        envApiRoot: '${this.portalDomain}/local/api',
        envGlobalRoot: '${this.portalDomain}/global',
    },
    inactivityTimeout: 15
};`;

        await fs.writeFileSync(environmentJsonPath, environmentJsonObj);

        // update the angularJson:
        const angularJsonPath = '../gftn-web/angular.json';
        const angularJsonTxt: string = await fs.readFileSync(angularJsonPath, "utf8");
        const angularJsonObj: any = JSON.parse(angularJsonTxt);

        const config = {
            optimization: true,
            outputHashing: "all",
            sourceMap: false,
            extractCss: true,
            namedChunks: false,
            aot: true,
            extractLicenses: true,
            vendorChunk: false,
            buildOptimizer: true,
            fileReplacements: [
                {
                    replace: "src/environments/environment.ts",
                    with: `src/environments/environment.${this.envShortName}.ts`
                }
            ]
        }

        // add new env config to the angular.json
        const newAngularObj = _.set(angularJsonObj, 'projects.gftn-web.architect.build.configurations.' + this.envShortName, config);

        await fs.writeFileSync(angularJsonPath, JSON.stringify(newAngularObj));

    }

    async buildAndDeployPortal() {

        const methodName = 'deployPortal';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = `cd ../gftn-web; ng build -c=${this.envShortName} ; firebase deploy --only hosting --project=${this.gcloudProjectName}`
            console.log('running: ' + script)

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    reject();
                    process.exit(1);
                }
            });

        });

    }

    async addFirebaseUseConfig() {

        const firebaseRcPath = '../gftn-web/.firebaserc';
        const firebaseRcTxt: string = await fs.readFileSync(firebaseRcPath, "utf8");

        const firebaseRcObj: any = JSON.parse(firebaseRcTxt);

        const newObj = _.set(firebaseRcObj, 'projects.' + this.envShortName, this.gcloudProjectName);

        await fs.writeFileSync(firebaseRcPath, JSON.stringify(newObj));

    }

    async applyFirebaseBoltRules() {

        const methodName = 'applyFirebaseBoltRules';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = `cd ../gftn-web; firebase-bolt database.rules.bolt ; firebase deploy --only database --project=${this.gcloudProjectName}`

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    reject();
                    process.exit(1);
                }
            });

        });

    }

    async deployFirebaseFunctions() {

        const methodName = 'deployFirebaseFunctions';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = `cd ../gftn-web ; npm i ; cd functions ; npm i ; firebase use ${this.envShortName} ; firebase deploy --only functions`;

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    reject();
                    process.exit(1);
                }
            });

        });

    }

    async addSuperUserToFirebaseDatabase() {

        const methodName = 'addUserToFirebaseDatabase';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = `firebase database:update /super_permissions --data '{"${this.adminUid}":{"email":"${this.superAdminEmail}","roles":{"admin":true}}}' --project=${this.gcloudProjectName} -y`
            script = script + ` ; firebase database:update /users --data '{"${this.adminUid}":{"profile":{"email":"${this.superAdminEmail}"}}}' --project=${this.gcloudProjectName} -y`

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    process.exit(1);
                    reject();
                }
            });

        });

    }

    async addSuperUserToFirebaseAuth() {

        // JSON - eg: https://firebase.google.com/docs/cli/auth#JSON
        const adminUser = {
            "users": [
                {
                    localId: "rootsuperuser", // UID
                    email: this.superAdminEmail
                }
            ]
        };

        // create a temp file (not pushed to git)
        const userJsonPath = './provisionAdminUser.json';

        // convert js object to json
        fs.writeFileSync(userJsonPath, JSON.stringify(adminUser));

        const methodName = 'addUserToFirebaseAuth';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = `firebase auth:import ${userJsonPath} --project=${this.gcloudProjectName}`;

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    process.exit(1);
                    reject();
                }
            });

        });

    }

    async addFirebase() {

        const methodName = 'addFirebase';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = `firebase projects:addfirebase ${this.gcloudProjectName}`;

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    process.exit(1);
                    reject();
                }
            });

        });

    }

    async setPortalDomain() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        const portalDomain = await new Promise(resolve => rl.question(this.fmtN('FgBlue') + '\n\nWhat domain would you like to host the portal at? (eg: worldwire-dev.io) ', ans => {
            rl.close();
            resolve(ans);
        }))

        this.portalDomain = portalDomain as string;

    }


    async setNewEnvName() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        // display question
        let result: string = await new Promise(resolve => rl.question(this.fmtN('FgBlue') + '\n\nGCloud & Firebase project name? (should be short alias, minimum 4 characters, NO Spaces, NO Special Characters, numbers are ok - eg: pen1) ', (ans: string) => {
            rl.close();
            resolve(ans);            
        }))

        // validate
        if( result.length < 4 ){
            
            // 4 char requirement is from gcloud
            console.info(this.fmtN('FgRed')+'must be more than 4 characters')
            
            // retry
            await this.setNewEnvName()
        }
        
        // set the short env name
        this.envShortName = result as string;

    }

    async setCredDir() {

        // ask user to provide firebase admin sdk credential from cmd line
        const rl1 = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        // get the firebase admin sdk credential output
        const currentCredVersion = await new Promise(resolve => rl1.question(this.fmtN('FgBlue') + `\n\nCurrent credentials version? (ie: .credentials-v[VERSION]/...) `, ( ans: string ) => {
            rl1.close();
            resolve(ans);
        }));

        this.credDir =`.credentials-v${currentCredVersion}/raw/${this.envShortName}`;
        
    }


    async setNewEnvLongName() {

        // ask user to name project from cmd line
        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        const result = await new Promise(resolve => rl.question(this.fmtN('FgBlue') + '\n\nHuman readable environment name? (eg: Development) ', ans => {
            rl.close();
            resolve(ans);
        }))

        this.envLongName = result as string;

    }

    /**
     * Name the super admin from cmd line
     *
     * @memberof main
     */
    async setInitialAdminEmail() {

        const rl = readline.createInterface({
            input: process.stdin,
            output: process.stdout,
            terminal: false
        });

        const superAdminEmail = await new Promise(resolve => rl.question(this.fmtN('FgBlue') + '\n\nInitial super admin email? (eg: chaseo@ibm.com) ', ans => {
            rl.close();
            resolve(ans);
        }))

        this.superAdminEmail = superAdminEmail as string;

        // console.log('super admin email is: ', this.superAdminEmail);

    }

    /**
     * create firebase and gcloud project
     *
     * @memberof main
     */
    async createGCloudProject() {

        const methodName = 'createGcloudProject';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = `gcloud projects create --name=${this.envShortName} --organization=${this.orgId} --quiet`;

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {

                console.log(methodName + ' stdout: ', outputText);
                let gcloudProjectName = '';

                try {
                    // parse output to get the project name
                    const start = 'projects/';
                    const end = '].';
                    gcloudProjectName = outputText.match(new RegExp(start + "(.*)" + end))[0];

                    console.log('before install script: ', gcloudProjectName);
                } catch (error) {
                    // do nothing
                    // console.log('')
                }


            });

            // log out error info
            cmd.stderr.on('data', (outputText: string) => {

                console.error(methodName + ' stderr: ', { data: outputText });

                try {
                    // parse output to get the project name
                    const start = 'projects/';
                    const end = '].';
                    this.gcloudProjectName = outputText.match(new RegExp(start + "(.*)" + end))[0];

                    // remove start and end text from parsed text
                    this.gcloudProjectName = this.gcloudProjectName.replace(start, '');
                    this.gcloudProjectName = this.gcloudProjectName.replace(end, '');

                    console.log('New GCloud Project Name: ', this.gcloudProjectName);
                } catch (error) {
                    // do nothing
                    // console.log('')
                }

            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                resolve();
            });

        });

    }

    /**
     * check if gcloud cli and firebase-tools cli installed and user logged in
     *
     * @returns
     * @memberof main
     */
    async checkDependencies() {

        const methodName = 'checkDependencies';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = 'gcloud --version ; gcloud -v ; gcloud organizations list'

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                // console.log(methodName + ' exit ', { code: code, signal: signal });
                if (code === 0) {
                    resolve();
                } else {
                    console.error(this.fmtN('FgRed') + '\n\nPlease install gcloud cli, firebase-tools, and firebase-bolt. Also login using \'gcloud auth login\' and \'firebase login\'' + this.fmtN('Reset'));
                    process.exit(1);
                    reject();
                }
            });

        });

    }

    // fmtSh(formatName: string) {
    //     const formats = {
    //         normal: '$(tput sgr0)',         // back to normal
    //         bold: '$(tput bold)',           // bold mode
    //         dim: '$(tput dim)',             // dim (half-bright) mode
    //         underline_on: '$(tput smul)',   // underline mode
    //         underline_off: '$(tput rmul)',  // underline mode off
    //         standout_on: '$(tput smso)',    // standout (bold) mode on 
    //         standout_off: '$(tput rmso)',   // standout (bold) mode off 
    //         black: '$(tput setaf 0)',       // black     COLOR_BLACK     0,0,0
    //         red: '$(tput setaf 1)',         // red       COLOR_RED       1,0,0
    //         green: '$(tput setaf 2)',       // green     COLOR_GREEN     0,1,0
    //         yellow: '$(tput setaf 3)',      // yellow    COLOR_YELLOW    1,1,0
    //         blue: '$(tput setaf 4)',        // blue      COLOR_BLUE      0,0,1
    //         magenta: '$(tput setaf 5)',     // magenta   COLOR_MAGENTA   1,0,1
    //         cyan: '$(tput setaf 6)',        // cyan      COLOR_CYAN      0,1,1
    //         white: '$(tput setaf 7)'       // white     COLOR_WHITE     1,1,1
    //     }
    //     console.log(formats[formatName]);
    //     return formats[formatName];
    // }

    /**
     * Format node output
     *
     * @param {string} formatName
     * @returns
     * @memberof main
     */
    fmtN(formatName: string) {
        const formats = {
            Reset: "\x1b[0m",
            Bright: "\x1b[1m",
            Dim: "\x1b[2m",
            Underscore: "\x1b[4m",
            Blink: "\x1b[5m",
            Reverse: "\x1b[7m",
            Hidden: "\x1b[8m",

            FgBlack: "\x1b[30m",
            FgRed: "\x1b[31m",
            FgGreen: "\x1b[32m",
            FgYellow: "\x1b[33m",
            FgBlue: "\x1b[34m",
            FgMagenta: "\x1b[35m",
            FgCyan: "\x1b[36m",
            FgWhite: "\x1b[37m",

            BgBlack: "\x1b[40m",
            BgRed: "\x1b[41m",
            BgGreen: "\x1b[42m",
            BgYellow: "\x1b[43m",
            BgBlue: "\x1b[44m",
            BgMagenta: "\x1b[45m",
            BgCyan: "\x1b[46m",
            BgWhite: "\x1b[47m",
        }
        // console.log(formats[formatName]);
        return formats[formatName];
    }

}

new main();