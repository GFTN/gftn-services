// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// NOTE: How to list files recursively in command prompt
// see https://superuser.com/questions/1010287/how-to-recursively-list-files-and-only-files-in-windows-command-propmpt
interface IMicroServiceDef { [microServiceDir: string]: IMicroServiceDefItem }

export interface IMicroServiceDefItem {
    // specific file paths from parent to includ in mirco-services
    paths: string[];
    // omit these depenedencies from shared dependencies
    omitDependencies: {};
    // name = micro-service name to go in package.json name
    name: string
    swagger: {
        title: string;
        host: string;
    };
    tsoa: {
        // list of the contollers to import into app.ts:
        controllerFiles: string[];
    }
}

// files shared with each microservice
export const sharedFiles: string[] = [
    'package.json',
    'tsoa.json',
    'tslint.json',
    'src/app.ts',
    'src/config.ts',
    'src/environment.ts',
    // use 'tsoa routes' on each micro-service to generate individually
    // 'src/routes.ts', 
    'src/auth/auth-helpers.ts',
    'src/auth/passport.ts',
    'src/auth/certs/cert.key',
    'src/auth/certs/cert.pem',
    'src/auth/certs/README.md',
    'src/auth/certs/req.cnf',
    // two-factor controller needed for middleware
    'src/controllers/two-factor.controller.ts',
    'src/constants/scripts.ts',
    'src/email/email.ts',
    'src/email/html-templates.ts',
    'src/email/plain-templates.ts',
    'src/middleware/auth.middleware.ts',
    'src/middleware/authentication.ts',
    'src/middleware/route-user-permissions.constants.ts',
    'src/models/auth.model.ts',
    'src/models/node.interface.d.ts',
    'src/models/participant.interface.d.ts',
    'src/models/token.interface.d.ts',
    'src/models/user.interface.d.ts',
    'src/shared/encryption.ts'
];

// specific files and dependencies need to separate out into separate micro services
export const microServices: IMicroServiceDef = {
    'micro-services/participant': {
        name: 'participant-portal-service',
        // specific files needed only for this micro service
        paths: [
            'src/controllers/permissions.controller.ts',
            'src/controllers/jwt.controller.ts'
        ],
        // depenedencies not needed for this micro-service
        omitDependencies: {
            "axios": "^0.18.0",
            "aws-sdk": "^2.416.0",
            "googleapis": "^39.2.0"
        },
        swagger: {
            title: 'Participant Portal Service',
            host: 'participant-portal-service-dot-next-gftn.appspot.com'
        },
        tsoa:{
            controllerFiles:[
                'permissions.controller.ts',
                'jwt.controller.ts'
            ]
        }
    },
    'micro-services/public': {
        name: 'public-portal-service',
        // specific files needed only for this micro service
        paths: [
            'src/controllers/ibmid.controller.ts',
            'src/controllers/helper.controller.ts'

        ],
        // depenedencies not needed for this micro-service
        omitDependencies: {
            "aws-sdk": "^2.416.0",
            "googleapis": "^39.2.0"
        },
        swagger: {
            title: 'Public Portal Service',
            host: 'public-portal-service-dot-next-gftn.appspot.com'
        },
        tsoa:{
            controllerFiles:[
                'ibmid.controller.ts',
                'helper.controller.ts'
            ]
        }
    },
    'micro-services/super': {
        name: 'super-portal-service',
        // specific files needed only for this micro service
        paths: [
            'src/controllers/automation.controller.ts',
            'src/controllers/institution.controller.ts',
            'src/env/awsParameter.ts',
            'src/env/awsSecret.ts',
            'src/env/env.ts',
            'src/env/utility/common.ts',
            'src/env/utility/var.ts'
        ],
        // depenedencies not needed for this micro-service
        omitDependencies: {
            "axios": "^0.18.0"
        },
        swagger: {
            title: 'Super Portal Service',
            host: 'super-portal-service-dot-next-gftn.appspot.com'
        },
        tsoa:{
            controllerFiles:[
                'automation.controller.ts',
                'institution.controller.ts'
            ]
        }
    }

};