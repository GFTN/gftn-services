// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import { use, expect } from 'chai'
import {before, describe, it, after} from 'mocha'
import { exec } from 'child_process';
import * as fs from 'fs';
// import { Encrypt } from './encrypt'
import 'chai-fs'
use(require('chai-fs'));

const workingTestDir = 'secret-mgmt/test';
const compressedFileName = 'test-dir.tgz'
let encryptionContainDir: string = '/enc';
let passContainDir: string = '/k8-secrets';
let rawContainDir: string = '/raw';
let original_extract_dev_creds: string = '';

before(() => {
    // make a copy of the const extract-debug-creds.sh to restore at the end of running test
    original_extract_dev_creds = fs.readFileSync('extract-debug-creds.sh', 'utf8')
});

describe('typescript secret management', () => {

    // TODO: new encrypt.ts test

});

describe('ci/cd (openssl) encryption workflow', () => {

    // using a high version number to avoid collision
    const currentTestVersion: number = 999;
    const nextTestVersion: number = currentTestVersion + 1;
    let prodDecryptScript = ''


    before(async () => {
        // start clean by deleting out temp test compressed files
        await new Promise((resolve, reject) => {

            const envArr = ['local', 'dev', 'qa', 'st', 'prod'];
            const cicdWorkingTestDir = 'secret-mgmt/cicdtest'

            // mock out stub out test dir
            let script = `
            rm -Rf .credentials ; \
            rm -Rf ${cicdWorkingTestDir} ; \
            mkdir ${cicdWorkingTestDir} ; \
            echo -n 'mockvalue' > ./cicd-cred-v${currentTestVersion}.tgz.enc ; \
            echo -n 'mockvalue' > ./cicd-cred-debug-v${currentTestVersion}.tgz.enc ; \
            mkdir ${cicdWorkingTestDir + rawContainDir} ; \
            mkdir ${cicdWorkingTestDir + encryptionContainDir} ; \
            mkdir ${cicdWorkingTestDir + passContainDir} ; \
            `

            // files to create
            let mockFiles = [
                cicdWorkingTestDir + '/enc/{{env}}/.cred_salt.txt',
                cicdWorkingTestDir + '/enc/{{env}}/.credentials.tgz.enc',
                cicdWorkingTestDir + '/k8-secrets/{{env}}/iv.txt',
                cicdWorkingTestDir + '/k8-secrets/{{env}}/pass.txt',
                cicdWorkingTestDir + '/raw/{{env}}/deploy/deploy.json',
                cicdWorkingTestDir + '/raw/{{env}}/img/adminsdk.json',
                cicdWorkingTestDir + '/raw/{{env}}/img/env.json',
                cicdWorkingTestDir + '/raw/{{env}}/img/secret.json',
                cicdWorkingTestDir + '/raw/{{env}}/portal.txt',
                cicdWorkingTestDir + '/fb_deploy.json',
                cicdWorkingTestDir + '/README.md',
                cicdWorkingTestDir + '/version'
            ]

            // each env
            for (let e = 0; e < envArr.length; e++) {
                script += `
                mkdir ${cicdWorkingTestDir + rawContainDir + '/' + envArr[e]} ; \
                mkdir ${cicdWorkingTestDir + rawContainDir + '/' + envArr[e] + '/deploy'} ; \
                mkdir ${cicdWorkingTestDir + rawContainDir + '/' + envArr[e] + '/img'} ; \
                mkdir ${cicdWorkingTestDir + encryptionContainDir + '/' + envArr[e]} ; \
                mkdir ${cicdWorkingTestDir + passContainDir + '/' + envArr[e]} ; \
                `
                // each file                
                for (let i = 0; i < mockFiles.length; i++) {
                    // write file contents
                    let path: string = mockFiles[i].replace('{{env}}', envArr[e])
                    script += `
                    echo -n 'testvalue' > ${path} ; \
                    `
                }
            }

            // set mock cred dir version
            script += `echo -n '${currentTestVersion}' > ${cicdWorkingTestDir}/version ; \ `

            script += `cp -Rf ${cicdWorkingTestDir} .credentials-v${currentTestVersion}`

            // execute shell
            const cmd = exec(script);

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error('stderr: ', { data: data as string });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                if (code === 0) {
                    resolve();
                } else {
                    console.error('faild to create the test dir for cicd test');
                    reject();
                }
            });

        });

    })

    before(async () => {

        // start clean by deleting out temp test compressed files
        await new Promise((resolve, reject) => {

            // rotate cicd credentials in addition to typescript output
            let script = `tsc -p secret-mgmt/tsconfig.json ; \
            version=${currentTestVersion} node ./secret-mgmt/lib/index.js `

            // execute shell
            const cmd = exec(script);

            // log out info
            // cmd.stdout.on('data', (data) => {
            //     console.error('stdout: ', { data: data as string });
            // });

            // output
            cmd.stdout.on('data', (outputText: string) => {
                if (outputText.includes('Run decryption using')) {
                    // console.log(' stdout: ', outputText);
                    const start = 'bash ./secret-mgmt';
                    const end = '# end';
                    const middleText = outputText.match(new RegExp(start + "(.*)" + end))[0];
                    // console.log(' middleText: ', middleText);
                    prodDecryptScript = middleText.replace('# end', '')
                }
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error('stderr: ', { data: data as string });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                if (code === 0) {
                    resolve();
                } else {
                    console.error('failed run ./secret-mgmt #1');
                    reject();
                }
            });

        });

    })

    it('should create local readable (with /raw) directory', () => {
        expect(`.credentials-v${nextTestVersion}`).to.be.a.directory();
    });

    it('should create updated encrypted files', () => {
        expect(`cicd-cred-debug-v${nextTestVersion}.tgz.enc`).to.be.a.file();
        expect(`cicd-cred-v${nextTestVersion}.tgz.enc`).to.be.a.file();
    });

    // removing old credentials is important because it will
    // force the ci/cd to error out if the secrets are rotated 
    // and old decryption keys are still present
    it('should remove old encrypted credentials', () => {
        expect(fs.existsSync(`cicd-cred-v${currentTestVersion}.tgz.enc`)).to.be.false;
        expect(fs.existsSync(`cicd-cred-debug-v${currentTestVersion}.tgz.enc`)).to.be.false;
    });

    it('should increment version file', () => {
        expect(`.credentials-v${nextTestVersion}/version`).to.be.a.file().with.content(String(currentTestVersion + 1));
    });

    it('should decrypt using key and iv', async () => {

        let scriptSuccess = false // default

        await new Promise((resolve, reject) => {

            // cmd
            let script = prodDecryptScript + ' ; sh extract-debug-creds.sh'

            // execute shell
            const cmd = exec(script);

            // count errors out by stderr (needed for scripts separated by ';')
            let _code = 0;

            // output
            cmd.stdout.on('data', (outputText: string) => {
                if (outputText.includes('decyption integrity compromised rotate travis creds immediately')) {
                    console.log('stdout: ', outputText);
                    _code = _code + 1;
                }
            });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error('stderr: ', { data: data as string });
                _code = _code + 1;
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                _code = code + _code;
                if (_code === 0) {
                    scriptSuccess = true;
                    resolve();
                } else {
                    reject();
                }
            });

        });

        // should be no error anticipated with running decryption
        expect(scriptSuccess).to.be.equal(true);
        expect(`.credentials/version`).to.be.a.file().with.content(String(nextTestVersion));
    });

    it('should have different salts, passpharases, initialization vectors', async () => {


        // start clean by deleting out temp test compressed files
        await new Promise((resolve, reject) => {

            // run the rotate again (twice total) to compare actual encrypted output
            // encrypted salts, pass, and iv to be different between next two version
            // since mock is just meaningless dummy value 
            let script = `tsc -p secret-mgmt/tsconfig.json ; \
            version=${nextTestVersion} node ./secret-mgmt/lib/index.js`

            // execute shell
            const cmd = exec(script);

            // log out info
            // cmd.stdout.on('data', (data) => {
            //     console.error('stdout: ', { data: data as string });
            // });

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error('stderr: ', { data: data as string });
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                if (code === 0) {
                    resolve();
                } else {
                    console.error('failed run ./secret-mgmt #2');
                    reject();
                }
            });

        });

        // arbitraily choosing one env to check 
        const env = 'st';

        let filesArr = [
            `.credentials-v{{version}}/enc/${env}/.cred_salt.txt`,
            `.credentials-v{{version}}/k8-secrets/${env}/iv.txt`,
            `.credentials-v{{version}}/k8-secrets/${env}/pass.txt`
        ]

        // create key value with filename and contents
        for (let i = 0; i < filesArr.length; i++) {
            const nextV0Contents = fs.readFileSync(filesArr[i].replace('{{version}}', String(nextTestVersion)))
            const nextV1Contents = fs.readFileSync(filesArr[i].replace('{{version}}', String(nextTestVersion + 1)))
            expect(nextV0Contents + '').to.not.equal(nextV1Contents + '')
        }

    });

    after(async () => {

        await new Promise((resolve, reject) => {

            // clean-up
            let script = `
            rm -Rf .credentials ; \
            rm -Rf .credentials-v${currentTestVersion} ; \
            rm -Rf .credentials-v${nextTestVersion} ; \
            rm -Rf .credentials-v${nextTestVersion + 1} ; \
            rm cicd-cred-debug-v${currentTestVersion}.tgz.enc ; \
            rm cicd-cred-v${currentTestVersion}.tgz.enc ; \
            rm cicd-cred-debug-v${nextTestVersion}.tgz.enc ; \
            rm cicd-cred-v${nextTestVersion}.tgz.enc ; \
            rm cicd-cred-debug-v${nextTestVersion + 1}.tgz.enc ; \
            rm cicd-cred-v${nextTestVersion + 1}.tgz.enc ; \
            `

            // execute shell
            const cmd = exec(script);

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                resolve();
            });
        });
    });

});

after(() => {
    // restor extract-debug-creds.sh
    fs.writeFileSync('extract-debug-creds.sh', original_extract_dev_creds, 'utf8')
});


