// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as fs from 'fs';
// import { exec } from 'child_process';
import {
    cloneDeep,
    last,
    forEach,
    size
} from 'lodash';
import { APP_YAML, GCLOUD_GITIGNORE, PACKAGE_JSON, TSCONFIG_PROD, TSCONFIG_DEV } from './templates';
import { sharedFiles, microServices, IMicroServiceDefItem } from './micro-services-constants';

export class MicroServicesDeployment {

    // private version: string;
    private microServices = microServices;
    private sharedFiles = sharedFiles;

    constructor(
        // api version to use for gcloud endpoints
        // version: string
    ) {
        // this.version = version;
    }

    async build() {
        await this.createServiceDirectories()
        await this.createAppDotYamls();
        await this.createPackageDotJsons();
        await this.createGcloudGitignores();
        await this.createGcloudTsconfig();
        await this.controllerImportsAppTs();
    }

    /**
     * creates file directories for micro-services in ./micro-services
     *
     * @memberof MicroServicesDeployment
     */
    private async createServiceDirectories() {

        // merge in shared files 
        forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {
            // add shared files
            this.microServices[microServiceDir].paths = this.microServices[microServiceDir].paths.concat(this.sharedFiles);
        });

        // delete out old files if they exist
        forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

            let deleteFolderRecursive = function (path: string) {
                if (fs.existsSync(path)) {
                    fs.readdirSync(path).forEach(function (file, index) {
                        let curPath = path + "/" + file;
                        if (fs.lstatSync(curPath).isDirectory()) {
                            // recurse
                            deleteFolderRecursive(curPath);
                        } else {
                            // delete file
                            fs.unlinkSync(curPath);
                        }
                    });
                    fs.rmdirSync(path);
                }
            };

            // init
            deleteFolderRecursive(microServiceDir);

        });

        // add copies from ./../functions (used for development)
        forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

            // add missing paths
            // NOTE: ok to include root functions and file name since it will be replaced
            // and only used to create the directory structure
            const _paths = val.paths.concat([
                // filenames below are not relevant, only paths: 
                // 'lib/app.js',
                'def/swagger.json',
                'build/package.json'
            ]);

            // copy over updated files ./../*
            forEach(_paths, (originalPath: string) => {

                // file to create
                const toPath = microServiceDir + originalPath.replace('functions', '');

                // dir path to create
                const dir = toPath.replace('/' + last(toPath.split('/')), '');

                fs.mkdirSync(dir, { recursive: true });
                console.log(dir + ' directory was created');

            });

        });

        // add copies from ./../functions (used for development)
        forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

            // copy over updated files ./../*
            forEach(val.paths, (originalPath: string) => {

                // file to create
                const toPath = microServiceDir + originalPath.replace('functions', '');

                fs.copyFileSync(originalPath, toPath);
                console.log(toPath + ' file was created');

            });

        });

    }

    /**
     * Creates the various package.json files for gcloud app engine micro-service deployments
     *
     * @returns
     * @memberof MicroServicesDeployment
     */
    private async createPackageDotJsons() {

        // only resovle once loop is over
        let count = 0;

        return new Promise((resolve, reject) => {

            fs.readFile('package.json', (err, data) => {
                if (err) {
                    console.error('unable to read package.json: ', err)
                    resolve();
                } else {

                    // get the current package.json used for development
                    const devPackageJson: { dependencies: {}, devDependencies: {} } = JSON.parse(data as any);

                    forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

                        // copy production dependencies 
                        const productionsDependenciesObj = cloneDeep(devPackageJson.dependencies);
                        const devDependenciesObj = cloneDeep(devPackageJson.devDependencies);

                        // remove out unused dependencies
                        forEach(val.omitDependencies, (depVersion: string, depKey: string) => {
                            delete productionsDependenciesObj[depKey];
                        });

                        // stringify the production dependencies to add to package.json
                        const prodDependencies = JSON.stringify(productionsDependenciesObj)
                        const devDependencies = JSON.stringify(devDependenciesObj)

                        // create package.json for micro service
                        fs.writeFile(microServiceDir + '/package.json', PACKAGE_JSON(prodDependencies, devDependencies, val.name), (err) => {

                            if (err) {
                                console.error('error creating package.json: ' + microServiceDir, err);
                            }

                            fs.writeFile(microServiceDir + '/build/package.json', PACKAGE_JSON(prodDependencies, '{}', val.name), (err2) => {

                                if (err2) {
                                    console.error('error creating package.json: ' + microServiceDir, err2);
                                }

                                // increment and check if loop is completed
                                count = count + 1;
                                if (count = size(this.microServices)) {
                                    resolve();
                                }

                            });

                        });

                    });

                }
            })

        });

    }

    /**
     * Dynamically sets the environment variables for the various micro-service app.yaml deployment files to Google Cloud App Engine
     *
     * @returns
     * @memberof MicroServicesDeployment
     */
    private async createAppDotYamls() {

        // get env vars from launch.json
        let launchJsonTxt = fs.readFileSync('.vscode/launch.json', "utf8") as string;
        launchJsonTxt = launchJsonTxt.replace(/(\/\/ (.*)$)/gm, '');

        // Remove comments from launch.json to parse it
        // IMPORTANT: not that to parse the launch.json properly  
        //      1) all comments must be have a space after '//'
        //      2) redundant commas ',' after a line will fail to parse json properly 
        // NOTE: a quick way to check the JSON output for proper formatting is to copy the  
        // launchJsonObj output via the "Debug Console" and past it into a postman json body
        // and look for formatting errors 
        const launchJsonObj = JSON.parse(launchJsonTxt);

        if (!launchJsonObj) {
            console.error('cannot find the related launch.json config')
        }

        // const envs = launchJsonObj.configurations[Number(process.env.launch_json_order)].env as { [name: string]: string };
        // const envs = global['envs']
        const e = { cred: global['envs'].cred };

        // only resovle once loop is over
        let count = 0;

        return new Promise((resolve, reject) => {

            forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

                // create app.yaml for micro service
                fs.writeFile(microServiceDir + '/build/app.yaml', APP_YAML(e, val.name), (err) => {

                    if (err) {
                        console.error('error creating app.yaml' + microServiceDir, err);
                    }

                    // increment and check if loop is completed
                    count = count + 1;
                    if (count = size(this.microServices)) {
                        resolve();
                    }

                })

            });

        });

    }

    /**
     * Creates the various gitignores for gcloud app engine micro-services deployment
     *
     * @returns
     * @memberof MicroServicesDeployment
     */
    private async createGcloudGitignores() {

        // only resovle once loop is over
        let count = 0;

        return new Promise((resolve, reject) => {

            forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

                fs.writeFile(microServiceDir + '/build/.gcloudignore', GCLOUD_GITIGNORE(), (err) => {

                    if (err) {
                        console.error('error creating .gcloudignore: ' + microServiceDir, err);
                    }

                    // increment and check if loop is completed
                    count = count + 1;
                    if (count = size(this.microServices)) {
                        resolve();
                    }

                })

            });

        });

    }

    /**
     * Creates the various tsconfigs micro-services dev build
     *
     * @returns
     * @memberof MicroServicesDeployment
     */
    private async createGcloudTsconfig() {

        // only resovle once loop is over
        let count = 0;

        return new Promise((resolve, reject) => {

            forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

                fs.writeFile(microServiceDir + '/tsconfig.dev.json', TSCONFIG_DEV(), (err) => {

                    if (err) {
                        console.error('error creating tsconfig: ' + microServiceDir, err);
                    }

                    fs.writeFile(microServiceDir + '/tsconfig.prod.json', TSCONFIG_PROD(), (err2) => {

                        if (err2) {
                            console.error('error creating tsconfig: ' + microServiceDir, err2);
                        }

                        // increment and check if loop is completed
                        count = count + 1;
                        if (count = size(this.microServices)) {
                            resolve();
                        }

                    });

                });

            });

        });

    }

    /**
     * updates the imports in the app.ts to allign with the imports needed for specific the micro-service
     *
     * @private
     * @returns
     * @memberof MicroServicesDeployment
     */
    private async controllerImportsAppTs() {

        // only resovle once loop is over
        let count = 0;

        return new Promise((resolve, reject) => {

            fs.readFile('src/app.ts', "utf8", (err, _appTsStr) => {

                if (err) {
                    console.error('unable to read app.ts: ', err)
                    resolve();
                } else {

                    // create an array where [1] index is to be the part replaced by specific tsoa contollers
                    const appTsStrArr = cloneDeep(_appTsStr.split('// ============ Controllers ============='));
                    // console.log(appTsStrArr);

                    // update controller for each micro-service
                    forEach(this.microServices, (val: IMicroServiceDefItem, microServiceDir: string) => {

                        // create controller imports to insert
                        let controllerImports = '';
                        forEach(val.tsoa.controllerFiles, (controllerFileName: string) => {
                            controllerImports = controllerImports +
                                "\r\nimport './controllers/" + controllerFileName.replace('.ts', '') + "';\r\n";
                        });

                        // replace current tsoa controller imports
                        appTsStrArr[1] = controllerImports;

                        // concat strings that make up app.ts
                        const newAppTs = appTsStrArr.join('');

                        // create package.json for micro service
                        fs.writeFile(microServiceDir + '/src/app.ts', newAppTs, (err) => {

                            // increment and check if loop is completed
                            count = count + 1;
                            if (count = size(this.microServices)) {
                                resolve();
                            }

                        });

                    });

                }

            });

        });

    }

}