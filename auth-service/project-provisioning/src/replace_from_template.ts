// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

// USAGE EXAMPLE: template_file="firebase-template.json" replace_file="firebase.json" template_txt="{{{PORTAL_DOMAIN}}}" replace_txt="worldwire-qa.io"  node ../auth-service/project-provisioning/lib/replace_from_template.js
import * as fs from 'fs';

export class Main {

    replace_file: string;
    template_file: string;
    template_txt: string;
    replace_txt: string


    constructor() {

        this.template_file = process.env.template_file
        this.replace_file = process.env.replace_file
        this.template_txt = process.env.template_txt
        this.replace_txt = process.env.replace_txt

        if(!this.template_file){
            console.log('Missing env var - template_file')
            process.exit(1)
        }

        if(!this.replace_file){
            console.log('Missing env var - replace_file')
            process.exit(1)
        }

        if(!this.template_txt){
            console.log('Missing env var - template_txt')
            process.exit(1)
        }

        if(!this.replace_txt){
            console.log('Missing env var - replace_txt')
            process.exit(1)
        }

        this.init();
    }

    async init() {

        try {

            // remove file if exists
            if (fs.existsSync(this.replace_file)) {
                fs.unlinkSync(this.replace_file);
            }

            // get template text
            const stringText = await fs.readFileSync(this.template_file, "utf8");  

            // replace text
            const updatedText = stringText.replace(this.template_txt, this.replace_txt);

            // write update text to replace file
            await fs.writeFileSync(this.replace_file, updatedText, "utf8");  

            console.log('success, ' + this.replace_file + ' updated')

            process.exit(0);

        } catch (error) {
            
            console.error('replace file failed');
            
            process.exit(1);

        }
        
    }

}

new Main();