// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

// USAGE EXAMPLE: str_original="" str_substr="" str_replace="" node ../auth-service/project-provisioning/lib/replace.js

export class Main {

    str_substr: string;
    str_original: string;
    str_replace: string;

    constructor() {

        this.str_original = process.env.str_original
        this.str_substr = process.env.str_substr
        this.str_replace = process.env.str_replace
        
        if(!this.str_original){
            console.log('Missing env var - str_original')
            process.exit(1)
        }

        if(!this.str_substr){
            console.log('Missing env var - str_substr')
            process.exit(1)
        }

        // empty string ok
        if(this.str_replace === undefined){
            console.log('Missing env var - str_replace')
            process.exit(1)
        }


        this.init();
    }

    async init() {

        try {

            const result = this.str_original.replace(this.str_substr,this.str_replace);
            console.log(result);            

            process.exit(0);

        } catch (error) {
            
            console.error('replacement sub str failed');
            
            process.exit(1);

        }
        
    }

}

new Main();