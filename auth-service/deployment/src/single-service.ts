// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as fs from 'fs';
import * as path from 'path';
import { APP_YAML, GCLOUD_GITIGNORE, PACKAGE_JSON } from './templates';
import { exec } from 'child_process';
import { IAppEnvs } from './shared/environment';
import * as _ from 'lodash';
import * as YAML from 'yamljs';

export class SingleServiceDeployment {

    // private ignoreErrorStringsIncluding = [];

    async build() {

        try {

            const env = process.env.env;

            let debugPath = '-debug'
            if (env !== 'dev' && env !== 'local') {
                console.info('Using production credentials to build project...')
                debugPath = '';
            }

            const credDir = `.credentials${debugPath}-v${process.env.cred_version}`;

            // check if the build dir exists otherwise create it
            const buildPath = './authentication/build';
            if (!fs.existsSync(buildPath)) {
                fs.mkdirSync(buildPath, { recursive: true } as any);
            }

            const serviceName = 'default' // 'gftn-authentication-service';

            const appEnvs = this.getAppEnvs(credDir, env);

            // permissions are defined in yml so that comments can be included 
            // and transpiled into json to work across languages (ie: node & Go)
            // furthermore we hardcode those middleware permissions into the codebase 
            // for the micro-services which the function below accomplishes
            const permissionsYamlPath = './authorization/middleware/permissions/archive/permissions.yml';
            const permissionsGoPath = './authorization/middleware/permissions/permissions.go';
            await this.replacePermissions(permissionsYamlPath, permissionsGoPath);

            await this.createAppDotYaml(appEnvs, serviceName);
            await this.createPackageDotJson(appEnvs, serviceName, buildPath);
            await this.createGcloudGitignore();
            // await this.createAppEngineEndpoints();
            await this.createDockerfile();

            // UPDATE: Now loading certs encrypted and base64 encoded ibmid certs dynamically
            // await this.copyIbmidCertsToBuild(buildPath, credDir);
            // this.copyFolderRecursiveSync(`${credDir}/.ibmid-certs`, `${buildPath}`);

            // .self-signed-certs - only needed for local development testing not a security risk 
            this.copyFolderRecursiveSync(`${credDir}/.self-signed-certs`, `${buildPath}`);

        } catch (error) {
            console.error(error)
            process.exit(1);
        }

    }

    /**
     * gets node application level encrypted secrets, decryption, and env vars
     *
     * @returns {IAppEnvs}
     * @memberof SingleServiceDeployment
     */
    getAppEnvs(credDir: string, env): IAppEnvs {

        // get env vars for deployment
        const enc = fs.readFileSync(`${credDir}/secret-mgr/${env}/.secret.enc`, "utf8")
        const decrypt = JSON.parse(fs.readFileSync(`${credDir}/secret-mgr/${env}/.decrypt.json`, "utf8"))
        const envs = JSON.parse(fs.readFileSync(`${credDir}/raw/${env}/img/env.json`, "utf8"))

        const appEnvs = {}

        // set envs 
        _.forEach(envs, (val: string, key: string) => {
            _.set(appEnvs, key, val);
        })

        // set decrypt 
        _.forEach(decrypt, (val: string, key: string) => {
            _.set(appEnvs, key, val);
        })

        // set enc
        _.set(appEnvs, 'enc', enc);

        return appEnvs as IAppEnvs;

    }

    // async copyIbmidCertsToBuild(buildPath: string, credDir: string) {

    //     console.log('copy certs');


    //     const methodName = 'copyIbmidCertsToBuild';

    //     return await new Promise((resolve, reject) => {

    //         // cmd
    //         // NOTE: need to copy over self signed for local debug ONLY
    //         // .certs are the organization's universal IBMid certs needed by passport
    //         console.log(` ${buildPath}/.certs vs ./authentication/build/.certs - ${credDir}/.certs vs .credentials-v17/.certs`);

    //         let script = `cp -Rf ${credDir}/.certs/. ${buildPath}/.certs ; cp -Rf ${credDir}/.self-signed-certs/. ${buildPath}/.self-signed-certs`

    //         // execute shell
    //         const cmd = exec(script);

    //         // output
    //         cmd.stdout.on('data', (outputText: string) => {
    //             console.log(methodName + ' stdout: ', outputText);
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
    //             if (_code === 0) {
    //                 resolve();
    //             } else {
    //                 reject(new Error(methodName + ' failed'));
    //             }
    //         });

    //     });

    // }

    private copyFileSync(source, target) {

        var targetFile = target;

        //if target is a directory a new file with the same name will be created
        if (fs.existsSync(target)) {
            if (fs.lstatSync(target).isDirectory()) {
                targetFile = path.join(target, path.basename(source));
            }
        }

        fs.writeFileSync(targetFile, fs.readFileSync(source));
    }

    private copyFolderRecursiveSync(source, target) {

        let self = this;

        var files = [];

        //check if folder needs to be created or integrated
        var targetFolder = path.join(target, path.basename(source));
        if (!fs.existsSync(targetFolder)) {
            fs.mkdirSync(targetFolder);
        }

        //copy
        if (fs.lstatSync(source).isDirectory()) {
            files = fs.readdirSync(source);
            files.forEach(function (file) {
                var curSource = path.join(source, file);
                if (fs.lstatSync(curSource).isDirectory()) {
                    self.copyFolderRecursiveSync(curSource, targetFolder);
                } else {
                    self.copyFileSync(curSource, targetFolder);
                }
            });
        }
    }

    /**
     * copy over dockerfile needed for build
     *
     * @returns
     * @memberof SingleServiceDeployment
     */
    async createDockerfile() {

        const methodName = 'createDockerfile';

        return await new Promise((resolve, reject) => {

            // cmd
            let script = 'cp ./deployment/docker/Dockerfile ./authentication/build ; cp ./authentication/src/docker-entrypoint.sh ./authentication/build'

            // execute shell
            const cmd = exec(script);

            // output
            cmd.stdout.on('data', (outputText: string) => {
                console.log(methodName + ' stdout: ', outputText);
            });

            let _code = 0;

            // log out error info
            cmd.stderr.on('data', (data) => {
                console.error(methodName + ' stderr: ', { data: data });
                _code = _code + 1;
            });

            // script competed
            cmd.once('exit', (code: number, signal: string) => {
                _code = code + _code;
                if (_code === 0) {
                    resolve();
                } else {
                    console.log(methodName + ' exit ', { code: _code, signal: signal });
                    reject(new Error(methodName + ' failed'));
                }
            });

        });

    }

    private createPackageDotJson(appEnvs: IAppEnvs, serviceName: string, buildPath: string) {

        const p = fs.readFileSync('authentication/package.json', 'utf8')

        // get the current package.json used for development
        const devPackageJson: { dependencies: {} } = JSON.parse(p as any);

        // parse the development package.json file
        // so that we get the production dependencies
        const productionDependencies = JSON.stringify(devPackageJson.dependencies)

        fs.writeFileSync(
            buildPath + '/package.json',
            PACKAGE_JSON(appEnvs as {}, serviceName, productionDependencies, '{}')
        );

        return;

    }

    // dynamically sets the environment variables for the app.yaml deployment file to Google Cloud App Engine
    private async createAppDotYaml(appEnvs, serviceName) {

        return new Promise((resolve, reject) => {

            fs.writeFile('authentication/build/app.yaml', APP_YAML(appEnvs as {}, serviceName), (err) => {

                if (err) {
                    console.error('error creating app.yaml for google cloud app engine, please ensure that ./build directory structure exists, error from index.ts: ', err);
                    resolve();
                } else {
                    resolve();
                }

            })

        });

    }

    private async createGcloudGitignore() {

        return new Promise((resolve, reject) => {

            fs.writeFile('authentication/build/.gcloudignore', GCLOUD_GITIGNORE(), (err) => {
                // success creating the app.yaml file
                if (err) {
                    console.error('error creating .gcloudignore for google cloud app engine, please ensure that ./build directory structure exists, error from index.ts: ', err);
                    resolve();
                } else {
                    resolve();
                }
            })

        });

    }
    
    /**
     * This function replaces the old permissions.go with a new permissions.go transpiled from permissions.yml
     *
     * @private
     * @param {string} srcPath
     * @param {string} destPath
     * @returns
     * @memberof SingleServiceDeployment
     */
    replacePermissions(srcPath: string, destPath: string): Promise<string> {
        return new Promise((res, rej) => {
            
            // convert yaml to json
            const json = this.transpileYamlToJson(srcPath);
            
            // create contents of the go file
            let goStr = `package permissions\n\n// Permissions : defined permissions for middleware\nfunc Permissions() string {\n\n\treturn ${'`' + json + '`'}\n\n}`;
            
            // write go file with permissions to middleware
            fs.writeFile(destPath, goStr, (err) => {
                if (!err) {
                    res('Permissions replacement succeeded!');
                }
                else {
                    rej(err);
                }
            });

        })
    }
    
    /**
     * This function transpiles YAML files to JSON strings.
     *
     * @private
     * @param {string} path
     * @returns {string}
     * @memberof SingleServiceDeployment
     */
    private transpileYamlToJson(path: string): string {
        const obj = YAML.parseFile(path);
        const json = JSON.stringify(obj, null, '\t');
        return json;
    }

}