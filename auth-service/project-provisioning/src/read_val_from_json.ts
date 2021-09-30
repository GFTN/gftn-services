// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// USAGE EXAMPLE: file=".credentials/$env/env.json" key="gae_service"  node ./project-provisioning/lib/read_val_from_json.js
import * as fs from 'fs';
import { exec } from 'child_process';

export class ReadJson {

    constructor() {
        this.init();
    }

    async pwdLs() {

        const methodName = 'pwdLs';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = 'echo \'From read_val_from_json.ts\' ; pwd ; ls'

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
                // let ignore = false;
                // for (let i = 0; i < this.ignoreErrorStringsIncluding.length; i++) {
                //     if (data.includes(this.ignoreErrorStringsIncluding[i])) {
                //         ignore = true;
                //     }
                // }
                // if (!ignore) {
                    _code = _code + 1;
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

    async init() {

        // await this.pwdLs()

        if(!process.env.key || !process.env.file) {
            console.error('missing json key or file name')
            process.exit(1);
        }

        try {
            const jsonTxt: string = await fs.readFileSync(process.env.file, "utf8");  
            const jsonObj = JSON.parse(jsonTxt);
            console.log(jsonObj[process.env.key]);
            process.exit(0)
        } catch (error) {
            console.error('unable to read value from json', error);
            process.exit(1);
        }
        
    }

}

new ReadJson();