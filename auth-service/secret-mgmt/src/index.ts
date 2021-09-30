// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as fs from 'fs'
import * as path from 'path'
import { Encrypt } from './encrypt'
import { exec } from 'child_process'
import * as _ from 'lodash'

class main {

    private currentVersion: number;
    private currentDir: string;
    private newDir: string;

    constructor() {
        this.start();
    }

    async start() {

        try {

            // get int representing the current version creds (since it can continually be rotated)
            await this.setCurrentVersion()

            // create copy of old creds in to new dir (incremented version)
            await this.createNewVersionDir();

            // encrypt/rotate all env creds
            const envArr = ['dev', 'qa', 'st', 'prod', 'local', 'pen1'];
            for (let i = 0; i < envArr.length; i++) {
                await this.encryptCredentials(envArr[i]);
            }

            // encrypt secrets for cicd pipeline, ie: cicd-cred-vN.tgz.enc
            // to be decrypted and consumed by the CICD provider secret management console
            await this.createCicdCredentials()

            console.info(this.fmtN('FgGreen') + '\nSuccess! Secrets have been encrypted for all environments \n' + this.fmtN('Reset'));

            process.exit();

        } catch (error) {
            console.info(this.fmtN('FgRed') + '\nError - failure creating secret encryption: ' + error + '. \n' + this.fmtN('Reset'));
            process.exit(1);
        }

    }

    async createCicdCredentials() {

        const methodName = 'createCicdCreds';

        console.log('starting createCicdCreds:');

        return await new Promise((resolve, reject) => {

            // console.log('current', this.currentVersion);

            // cmd
            let script = `bash secret-mgmt/cicd/generate.sh ${this.currentVersion}`

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(outputText);
            });

            // count errors out by stderr (needed for scripts separated by ';')
            let _code = 0;

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data as string });
                _code = _code + 1;
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                _code = code + _code;
                if (_code === 0) {
                    console.log('createCicdCreds succeeded');
                    resolve();
                } else {
                    console.error('createCicdCreds failed in ./secret-mgmt');
                    reject(new Error(methodName + ' failed'));
                }
            });

        });

    }

    /**
     * set dynamically by looking up the most recent version from file system
     *
     * @memberof main
     */
    async setCurrentVersion() {

        // loop increment dir names until version is found and then subsequent version is not found
        let version = 0;

        const versionsFound: number[] = [];

        // stop loop if no version is found by 1000
        const maxVersion = 1000

        for (let i = 0; i < maxVersion; i++) {
            if (fs.existsSync('.credentials-v' + version)) {
                versionsFound.push(version)
            }

            // increment version number 
            version+=1;
        }

        if(versionsFound.length <= 0){
            console.log('No un-encrypted .credentials-v{VERSION_NUMBER} available. Unable to rotate. Please download () and extract into ./gftn-services/auth-service and re-run');            
            process.exit(1);
        }

        this.currentVersion = _.max(versionsFound);
        console.log('max version found: ', this.currentVersion);
        this.currentDir = `.credentials-v${this.currentVersion}`
        this.newDir = `.credentials-v${this.currentVersion + 1}`

    }

    async createNewVersionDir() {

        const methodName = 'createNewRaw';

        return await new Promise((resolve, reject) => {

            // make copy of credentials with increment version
            let script = `cp -Rf ${this.currentDir}/. ${this.newDir} ; \\
             echo -n '${this.currentVersion + 1}' > ${this.newDir}/version`

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            // count errors out by stderr (needed for scripts separated by ';')
            let _code = 0;

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data as string });
                _code = _code + 1;
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                _code = code + _code;
                if (_code === 0) {
                    resolve();
                } else {
                    reject(new Error(methodName + ' failed'));
                }
            });

        });

    }

    async encryptCredentials(env: string | 'prod' | 'dev' | 'qa' | 'st'): Promise<void> {

        const c = new Encrypt();

        // get plain text to encrypt
        const secretPlainTxt = fs.readFileSync(`${this.newDir}/raw/${env}/img/secret.json`, 'utf8')

        // encrypt text
        const e = await c.encrypt(secretPlainTxt);

        // write output files:
        
        // encrypted secrets file
        this.writeFile(`${this.newDir}/secret-mgr/${env}/.secret.enc`, e.enc)
        
        // decryption details - NEVER PUSH TO GIT
        this.writeFile(`${this.newDir}/secret-mgr/${env}/.decrypt.json`, JSON.stringify(e.dec));

        return;
        
    }

    writeFile(filePath: string, content: string ){
        
        const dirPath = path.dirname(filePath)
        
        if(!fs.existsSync(dirPath)){
            fs.mkdirSync(dirPath, { recursive: true });
        }

        fs.writeFileSync(filePath, content, 'utf8')
        return;

    }

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