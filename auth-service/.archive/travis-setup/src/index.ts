// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { exec } from 'child_process';
import * as fs from 'fs';
// import { WWBox } from './box/box';

class main {

    ignoreErrorStringsIncluding = [
        'detected repository as gftn/gftn-services',
        'warning: Insecure'
    ]

    // to be inserted into .travis.yml see .travis.yml 
    travisDecryptScript: string;

    // array of secret file paths to include in the .tar 
    secretFilePathsArr = [
        '.credentials/dev/admin-service-account.json',
        '.credentials/dev/client-service-account.json'
    ];

    // encrypt env vars in travis 
    envArr: { name: string, val: string, valBase64: string }[] = [
        {
            name: 'some-name',
            val: '',
            valBase64: '' // to be dynamically set in base64
        }
    ];

    constructor() {
        this.travisDecryptScript = '';
        this.start();
    }

    async start() {

        // // init box sdk
        // const box = new WWBox();
        // await box.init();        
        // // wait 3s
        // await new Promise(resolve => setTimeout(resolve, 3000))
        // // try box sdk
        // await box.upload().then(file => {
        //     console.info(file);
        // }, err => {
        //     console.info(err);
        // });

        try {
            await this.checkDependencies();
            await this.compressCredentials();
            await this.encryptCredentialsForTravis();
            await this.updateTravisYaml();
            await this.cleanUp();
            console.info(this.fmtN('FgBlue') + '\nSuccess! Travis encrypted secrets file has been updated along with World Wires ci/cd travis.yaml. Simply commit your changes including the updated .travis-credentials.tar.enc to github and travis will take care of the rest. \n');
            process.exit();
        } catch (error) {
            console.info(this.fmtN('FgError') + '\nError - Setting up travis encryption failed: ' + error + '. \n');
            process.exit(1);
        }

    }

    async checkDependencies() {

        const methodName = 'checkDependencies';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = 'travis -v ; travis whoami'

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
                let ignore = false;
                for (let i = 0; i < this.ignoreErrorStringsIncluding.length; i++) {
                    if (data.includes(this.ignoreErrorStringsIncluding[i])) {
                        ignore = true;
                    }
                }
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
                    console.error(this.fmtN('FgRed') + '\nERROR - MISSING DEPENDENCY: Please install travis cli (https://github.com/travis-ci/travis.rb) and login. Please run: \n$ travis login -e "https://travis.ibm.com/api" --github-token=personal-access-token-from-githubenterprise\n')
                    reject(new Error(methodName + ' failed'));
                }
            });

        });

    }

    async compressCredentials() {

        const methodName = 'setTravisEncryptedFiles';

        return await new Promise((resolve, reject) => {

            // delete existing .enc file if exists
            if (fs.existsSync('travis-credentials.tar.enc')) {
                fs.unlinkSync('travis-credentials.tar.enc');
            }

            // cmd
            let script = 'tar cvf travis-credentials.tar .credentials'

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
                console.error(methodName + ' stderr: ', { data: data });
                _code = _code + 1;
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

    async encryptCredentialsForTravis() {

        const methodName = 'getTravisFileEncryption';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = 'yes | travis encrypt-file travis-credentials.tar'

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
                if(outputText.includes('openssl aes-256-cbc')){
                    const start = 'openssl';
                    const end = '-d';
                    this.travisDecryptScript = outputText.match(new RegExp(start + "(.*)" + end))[0];
                    console.log(methodName + ' middleText: ', this.travisDecryptScript);
                }                
            });

            // count errors out by stderr (needed for scripts separated by ';')
            let _code = 0;

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data as string });
                let ignore = false;
                for (let i = 0; i < this.ignoreErrorStringsIncluding.length; i++) {
                    if (data.includes(this.ignoreErrorStringsIncluding[i])) {
                        ignore = true;
                    }
                }
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

    async updateTravisYaml() {

        const travisYmlPath = '../.travis.yml';

        let travisText = await fs.readFileSync(travisYmlPath, "utf8") as string;

        // get previous travis encrypt command to replace
        const start = '{{{travis-decrypt-credentials-start}}}';
        const end = '{{{travis-decrypt-credentials-end}}}';

        // find text to replace
        const replacementTxt = travisText.substring(
            travisText.lastIndexOf(start) + start.length,
            travisText.lastIndexOf(end) - 2
        );

        // loop multiple times incase there multiple places to replace text in file
        for (let i = 0; i < 10; i++) {
            // replace the text in the travis file
            travisText = travisText.replace(replacementTxt, '\n        - ' + this.travisDecryptScript + '\n        ');
        }

        // delete existing .travis.yml file and replace with updated one
        if (fs.existsSync(travisYmlPath)) {
            fs.unlinkSync(travisYmlPath);
        }

        // write updated .travis.yml
        fs.writeFileSync(travisYmlPath, travisText);

    }

    async cleanUp() {

        const methodName = 'cleanUp';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = 'rm travis-credentials.tar'

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
                console.error(methodName + ' stderr: ', { data: data });
                console.error(methodName + ' stderr: ', { data: data as string });
                let ignore = false;
                for (let i = 0; i < this.ignoreErrorStringsIncluding.length; i++) {
                    if (data.includes(this.ignoreErrorStringsIncluding[i])) {
                        ignore = true;
                    }
                }
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