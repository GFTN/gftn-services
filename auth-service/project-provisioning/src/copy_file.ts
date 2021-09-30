// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// USAGE EXAMPLE: file=".credentials/$env/env.json" key="gae_service"  node ./project-provisioning/lib/read_val_from_json.js
import * as fs from 'fs';
import * as path from 'path';
import * as _ from 'lodash';

export class main {

    fromPath: string
    toPath: string
    comment: '// ' | '# ' | string
    
    // Optional: 
    searchTxt: string
    replaceTxt: string

    constructor() {
        this.init();
    }

    async init() {

        this.fromPath = process.env.fromPath;
        this.toPath = process.env.toPath;
        this.comment= process.env.comment;

        this.searchTxt = process.env.searchTxt;
        this.replaceTxt = process.env.replaceTxt;
        

        if(!this.fromPath){
            console.error('missing \'fromPath\' from copy_file.ts')
            process.exit(1)
        }

        if(!this.toPath){
            console.error('missing \'toPath\' from copy_file.ts')
            process.exit(1)
        }

        if(!this.comment){
            console.error('missing \'comment\' from copy_file.ts')
            process.exit(1)
        }

        try {
            await this.copy();
            process.exit(0)
        } catch (error) {
            console.error('error: see copy_file.ts: ', error);
            process.exit(1);
        }
        
    }

    async copy(){
        
        const txt = fs.readFileSync(this.fromPath, 'utf8')
        let body = `
${this.comment} IMPORTANT: DO NOT MODIFY - THIS FILE IS AUTO-GENERATED FROM ${this.fromPath} 
${txt}        
`
        if(this.searchTxt){
           body = this.overrideContents(body, this.searchTxt, this.replaceTxt)
        }

        const dirPath = path.dirname(this.toPath)
        
        if(!fs.existsSync(dirPath)){
            fs.mkdirSync(dirPath, { recursive: true });
        }

        fs.writeFileSync(this.toPath, body, 'utf8')

    }

    private overrideContents(fullTxt: string, searchTxt: string, replaceTxt: string){
        return _.replace(fullTxt, new RegExp(searchTxt, 'g'), replaceTxt);
    }

}

new main();