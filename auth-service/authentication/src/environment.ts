// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { AppOptions, credential } from 'firebase-admin';
import { IEncodeFirebaseCred, IRandomPepperObj } from './models/auth.model';
import * as fs from 'fs';
import * as _ from 'lodash';
import { Encrypt, IDecrypt } from './shared/encrypt';
import { exec } from 'child_process';

/**
 * Environment variables consumed at the application layer
 *
 * @interface IEnvs
 */
export interface IEnvs {
    launch_json_order: string;
    totp_label: string;
    build: 'prod' | 'dev' | string;
    build_for:
    // VMs
    'nodejs' |
    // firebase functions
    'firebase' |
    // google cloud functions
    'gcloud' |
    // needed for linter
    string;
    api_port: string;
    site_root: string;
    app_root: string;
    gae_service: string;
    ibmid_authorization_url: string;
    ibmid_token_url: string;
    ibmid_issuer_id: string;
    enable_2fa: string;
    fb_database_url: string;
    refresh_mins: string;
    initial_mins: string;
    // env = targeted env - dev, qa, tn, st
    env: 'dev' | 'qa' | 'tn' | 'st' | string
}

export interface IAppEnvs extends IEnvs, IDecrypt {
    enc: string;
}

/**
 * Secrets consumed at the application layer
 *
 * @interface ISecrets
 */
export interface ISecrets {
    send_in_blue_api_key: string;
    passport_secret: string;
    ibmid_client_id: string;
    ibmid_client_secret: string;
    fb_admin: IEncodeFirebaseCred;
    ww_jwt_pepper_obj: IRandomPepperObj;
    tmp: {
        name: string; value: string
    }[]
}

interface IDynamicEnv {
    firebaseConfig: AppOptions;
    ibmId_callback_url: string;
    ibmId_encoded_certs: string[]
}

export interface IGlobalEnvs extends ISecrets, IEnvs, IDynamicEnv { }

export class Environment {

    async init() {

        // // storing certs in tmp dir 
        // // see https://cloud.google.com/appengine/docs/standard/nodejs/runtime#filesystem
        // let tmpPath = '/tmp';

        // // only for development
        // if (process.env.build === 'dev') {
        //     tmpPath = './authentication/tmp';
        // }

        // // create /tmp dir if it does not exist 
        // if (!fs.existsSync(tmpPath)) {
        //     fs.mkdirSync(tmpPath, { recursive: true } as any)
        // }

        const secretsTxt = await this.decryptSecrets();
        
        const secrets = JSON.parse(secretsTxt) as ISecrets;
        
        this.setAppEnvs(secrets/*, tmpPath*/);
        
        // await this.script(`pwd ;  ls -a ; cd / ; ls -a`);        
    }

    /**
     * extracts credentials in plain text format
     *
     * @private
     * @returns {Promise<{envsTxt: string, secretsTxt: string}>}
     * @memberof Environment
     */
    private async decryptSecrets(): Promise<string> {

        try {

            // set required default envs
            const env = process.env.env;
            let pass = process.env.pass;
            let iv = process.env.iv;
            let salt = process.env.salt;
            let iterate = Number(process.env.iterate);
            let enc = process.env.enc;
            
            // override default envs from file (for local development ONLY)
            if (process.env.cred_version) {

                let debugPath = '-debug'
                if (env !== 'dev' && env !== 'local') {
                    debugPath = '';
                }

                const credDir = `.credentials${debugPath}-v${process.env.cred_version}`;
                
                enc = await fs.readFileSync(`${credDir}/secret-mgr/${env}/.secret.enc`, "utf8")
                const decrypt = await fs.readFileSync(`${credDir}/secret-mgr/${env}/.decrypt.json`, "utf8")
                const envs = await fs.readFileSync(`${credDir}/raw/${env}/img/env.json`, "utf8")

                // set envs from file on env
                _.forEach(JSON.parse(envs), (val, key) => {
                    process.env[key] = val;
                })

                const d: IDecrypt = JSON.parse(decrypt)

                // override envs
                pass = d.pass;
                iv = d.iv;
                salt = d.salt;
                iterate = Number(d.iterate);
            }

            // check require envs
            if (
                !pass ||
                !iv ||
                !salt ||
                !iterate ||
                !enc
            ) {
                console.error('Missing secret decryption info, see environment.ts');
                process.exit(1)
            }

            const c = new Encrypt();

            // decrypt to file
            return await c.decrypt(pass, salt, iterate, iv, enc);


        } catch (error) {

            console.error('Fatal Error: cannot decrypt auth-service credentials. Exiting...', error)
            // terminate application because no need to keep running if 
            // it cannot decrypt credentials

            return process.exit(1);

        }

    }

    async logPath() {

        const methodName = 'logPath';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = 'pwd ; ls -a ; echo "now /app\n" ; ls -a /app ; echo "now app\n" ; ls -a app ;  echo "now .credentials\n" ; ls -a .credentials ;  echo "now .credentials/live\n" ; ls -a .credentials/live'

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data as string });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                resolve();
            });

        });

    }

    /**
     * set shared env vars regardless of deployment environment
     *
     * @memberof Environment
     */
    private setAppEnvs(secrets: ISecrets /*, tmpPath: string*/) {

        // create global object
        let g = {} as IGlobalEnvs;

        const envs: IEnvs = process.env as any;

        _.set(g, 'env', envs.env);

        _.set(g, 'totp_label', envs.totp_label);

        _.set(g, 'build', envs.build);
        _.set(g, 'build_for', envs.build_for);

        // apiRoot and siteRoot should NOT end in '/'
        _.set(g, 'site_root', envs.site_root);
        _.set(g, 'app_root', envs.app_root);

        // config IBMId variables
        _.set(g, 'ibmid_client_id', secrets.ibmid_client_id);
        _.set(g, 'ibmid_client_secret', secrets.ibmid_client_secret);

        _.set(g, 'ibmid_authorization_url', envs.ibmid_authorization_url);
        _.set(g, 'ibmid_token_url', envs.ibmid_token_url);
        _.set(g, 'ibmid_issuer_id', envs.ibmid_issuer_id);

        _.set(g, 'enable_2fa', envs.enable_2fa);

        _.set(g, 'passport_secret', secrets.passport_secret);

        _.set(g, 'api_port', Number(envs.api_port));

        _.set(g, 'fb_database_url', envs.fb_database_url);

        _.set(g, 'firebaseConfig', {
            credential: credential.cert({
                projectId: secrets.fb_admin.project_id,
                clientEmail: secrets.fb_admin.client_email,
                privateKey: secrets.fb_admin.private_key
            }),
            databaseURL: envs.fb_database_url
        });

        // https://stackoverflow.com/questions/16891729/best-practices-salting-peppering-passwords
        // set pepper env var for jwt tokens (can be rotated using jwt/rotate-pepper)
        _.set(g, 'ww_jwt_pepper_obj', secrets.ww_jwt_pepper_obj);

        try {
            _.set(g, 'initial_mins', Number(envs.initial_mins));
        } catch (error) {
            console.error('failed to set env for initial_mins for jwt')
        }

        try {
            _.set(g, 'refresh_mins', Number(envs.refresh_mins));
        } catch (error) {
            console.error('failed to set env for refresh_mins for jwt')
        }

        // same for both prod and dev:
        _.set(g, 'ibmId_callback_url', g.app_root + '/sso/callback');

        _.set(g, 'send_in_blue_api_key', secrets.send_in_blue_api_key);

        _.set(g, 'ibmId_encoded_certs', []);

        _.forEach(secrets.tmp, (obj) => {

            // // decode base64 and write to /tmp file system
            // let certText = Buffer.from(obj.value, 'base64').toString();
            // // write to cert to /tmp file
            // console.log(`writing ${tmpPath}/${obj.name}`);
            // fs.writeFileSync(`${tmpPath}/${obj.name}`, certText, 'utf8')

            // create an array of base64 encoded certs to be 
            // consumed by custom passport strategy 

            // for dev only
            if (g.build === 'dev' && obj.name === 'prepiam_toronto_ca_ibm_com.crt') {
                g.ibmId_encoded_certs.push(obj.value);
            }

            // for prod only
            if (g.build === 'prod' && obj.name === 'idaas_iam_ibm_com.crt') {
                g.ibmId_encoded_certs.push(obj.value);
            }

            // for both prod and dev
            if (
                obj.name === 'digicert-root.pem' ||
                obj.name === 'IBMid-server.crt'
            ) {
                g.ibmId_encoded_certs.push(obj.value);
            }
            
        });

        // set as a global var
        global['envs'] = g;

    }

    async script(script: string) {

        const methodName = 'viewPathDetails';

        return await new Promise((resolve, reject) => {

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
                // const start = 'START TEXT HERE';
                // const end = 'END TEXT HERE';
                // const middleText = outputText.match(new RegExp(start + "(.*)" + end))[0];
                // console.log(methodName + ' middleText: ', middleText);
            });

            // count errors out by stderr (needed for scripts separated by ';')
            let _code = 0;

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data as string });
                let ignore = false;
                // for (let i = 0; i < this.ignoreErrorStringsIncluding.length; i++) {
                //     if (data.includes(this.ignoreErrorStringsIncluding[i])) {
                //         ignore = true;
                //     }
                // }
                if (!ignore) {
                    _code = _code + 1;
                }
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

}
