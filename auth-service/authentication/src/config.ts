// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
import * as cors from 'cors';
import * as express from 'express';
import * as admin from 'firebase-admin';
import * as fs from 'fs';
import * as http from 'http';
import * as https from 'https';
import { PassportConfig } from './auth/passport';
import * as helmet from 'helmet';
import { IGlobalEnvs } from './environment';
import { constants } from 'crypto';
import { TLSSocket } from 'tls';

export class Config {

    app: express.Application;
    private allowedOrigins: string[];

    rt_db: admin.database.Database;
    private env: IGlobalEnvs;

    async init() {

        this.env = global['envs']

        // create express app
        this.app = express();

        // set origins and firebase db
        this.setConfig();

        // this.setTLS();

        // set cors for certain origins
        this.app.use(this.setCorsConfig());

        // set test endpoints
        this.test();

        this.app = new PassportConfig(this.app).app;

        // create firebase real-time database ref
        const rt_db = admin.database();
        this.rt_db = rt_db;

    }

    setTLS(){

        // TODO: Following middleware is a work in progress
        // Disabling weak ciphers has should be enabled at the 
        // network level for now

        // disable tls 1.0 and 1.1 weak ciphers
        this.app.use((req, res, next) => {

            try {

                // // The following is not supported by app-engine deployment 
                // // (ie: getCipher getProtocol work locally but not live)
                // const cipher = ((req.socket) as TLSSocket).getCipher()
                const protocol = ((req.socket) as TLSSocket).getProtocol()
                // console.info('protocol: ', protocol);
                // console.info('cipher: ', cipher);

                // // The following does work on gcloud but the API  
                // // middleware response handling does not function
                // // as expected 
                // const soc = req.socket;
                // const tlsSocket = new TLSSocket(soc);
                // const protocol = tlsSocket.getProtocol()

                // const protocol = ((req.socket) as TLSSocket).getProtocol()
                // console.log('cipher: ', cipher);
                // output eg: { name: 'ECDHE-RSA-AES128-GCM-SHA256', version: 'TLSv1/SSLv3' }
                // console.log('protocol: ', protocol);
                // output eg: TLSv1.2

                if (protocol === 'TLSv1.2' || protocol === 'TLSv1.3') {
                    next();
                } else {
                    res.status(426);
                    res.send('request requires TLSv1.2 or greater, please upgrade. Inbound: ' + protocol);
                }

            } catch (error) {
                res.status(500);
                res.send('unexpected tls error')
            }

        });

    }


    /**
     * Starts the project in either a 'firebase functions' or 'nodeJs' configuration.
     * DEFAULT is 'firebase'
     *
     * @param {('firebase' | 'nodejs')} [buildFor='firebase']
     * @returns {(false)} If false, then not building in firebase functions configuration
     * @memberof App
     */
    // start(): functions.HttpsFunction | false {
    start(): false {

        // // returns an api that can be deployed as a serverless function in firebase
        // if (this.env.build_for === 'firebase') {
        //     // used for both dev and prod deployment to firebase
        //     // debug using:
        //     // terminal 1 > tsc -w
        //     // terminal 2 > firebase serve --only functions
        //     return functions.https.onRequest(this.app);
        // }

        // builds running a nodejs server running on a specified port
        if (this.env.build_for === 'nodejs') {
            this.listen();
        }

        return false;

    }

    private test() {

        // IMPORTANT: Only expose test endpoints in dev so as not
        // increase the public *attack* surface
        if (this.env.build === 'dev') {

            this.app.get('/test', (req: any, res: any) => {
                console.info('Test Log from calling /test');
                res.send('Fired /test successfully. Please check log for test log output.');
            });

            this.app.get('/', (req: any, res: any) => {
                res
                    .status(200)
                    .send('World Wire Authentication Service - running...')
                    .end();
            });

        }

    }

    private setConfig() {

        // setup configs based on build target:
        if (this.env.build === 'dev') {
            // Development configurations
            this.devConfig();
        } else {
            // production configurations
            this.prodConfig();
        }

    }

    private prodConfig() {

        // set basic security related middleware
        // https://expressjs.com/en/advanced/best-practice-security.html
        this.app.use(helmet());

        // view logs to ensure proper deployment
        console.info('===> Running PRODUCTION Configuration <===');

        // cors domains to allow
        this.allowedOrigins = [
            this.env.site_root,
            this.env.site_root + '/'
        ];
        // this.allowedOrigins = [
        //     'https://worldwire.io',
        //     'https://next.worldwire.io',
        //     'https://pen.worldwire.io',

        //     'https://worldwire.io/',
        //     'https://next.worldwire.io/',
        //     'https://pen.worldwire.io/'
        // ];

        this.app.use((req, res, next) => {
            const allowedOrigins = this.allowedOrigins;
            const origin = req.headers.origin as string;
            if (allowedOrigins.indexOf(origin) > -1) {
                res.setHeader('Access-Control-Allow-Origin', origin);
            }
            res.header('Access-Control-Allow-Methods', 'GET, PUT, POST, OPTIONS, DELETE');
            res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization');
            res.header('Access-Control-Allow-Credentials', 'true');
            return next();
        });

        // initialize firebase app:
        admin.initializeApp(this.env.firebaseConfig);

    }

    private devConfig() {

        this.app.use(helmet({
            hsts: false
        }));

        // view logs to ensure proper deployment
        console.info('===> Running DEVELOPMENT Configuration <===');

        // // set and log server port
        // const listener = this.express_app.listen(8888, function () {
        //     console.log('Listening on port ' + listener.address().port); //Listening on port 8888
        // });

        // not a production build, so add testing domains to cors
        this.allowedOrigins = [

            // angular app request
            'http://localhost:4200',
            'http://localhost:4200/',

            // firebase serve
            'http://localhost:5000',
            'http://localhost:5000/',

            // express app
            'http://localhost:3000',
            'http://localhost:3000/',

            // google cloud function emulator
            'http://localhost:8010',
            'http://localhost:8010/',

            'https://worldwire.io',
            'https://worldwire.io/',

            'https://next.worldwire.io',
            'https://next.worldwire.io/'
        ];

        this.app.use((req, res, next) => {
            const allowedOrigins = this.allowedOrigins;
            const origin = req.headers.origin as string;
            if (allowedOrigins.indexOf(origin) > -1) {
                res.setHeader('Access-Control-Allow-Origin', origin);
            }
            res.header('Access-Control-Allow-Methods', 'GET, PUT, POST, OPTIONS, DELETE');
            res.header('Access-Control-Allow-Headers', 'Content-Type, Authorization');
            res.header('Access-Control-Allow-Credentials', 'true');
            return next();
        });

        // initialize firebase app:
        admin.initializeApp(this.env.firebaseConfig);

    }

    private setCorsConfig(): express.RequestHandler {

        // https://github.com/expressjs/cors

        // set cors configuration for all routes
        return cors({
            origin: (origin, callback) => {
                // allow requests with no origin
                // (like mobile apps or curl requests)
                if (!origin) {
                    return callback(null, true);
                }

                if (this.allowedOrigins.indexOf(origin) === -1) {
                    const msg = 'The CORS policy for this site does not ' +
                        'allow access from the specified Origin.';
                    return callback(new Error(msg), false);
                }

                return callback(null, true);
            },
            // allows IBMId credentials to be sent along in headers from client
            credentials: true
        });

    }

    private listen() {

        // local machine development build 
        if (this.env.build === 'dev') {

            // Set certs for https (needed for testing with IBMid
            // since IBMid requires self-signed SSL)

            // self-signed-certs are ONLY used for local development with IBMId
            let cred_dir = '.credentials-debug-v29';
            if (!fs.existsSync(cred_dir)) {
                // must be debugging with ./authentication/build
                cred_dir = `../../${cred_dir}`;
            }
            const httpsOptions: https.ServerOptions = {
                key: fs.readFileSync(`${cred_dir}/.self-signed-certs/` + 'cert.key', "utf8"),
                cert: fs.readFileSync(`${cred_dir}/.self-signed-certs/` + 'cert.pem', "utf8")
            };

            const port = this.env.api_port;

            https.createServer(httpsOptions, this.app)
                .listen(port, () => {
                    // log when server was started
                    console.log('Started at: ', new Date);
                    console.log('Server running at ' + port);
                });

        }

        // production build
        if (this.env.build === 'prod') {

            // // using express:

            // google cloud listen
            // const PORT = process.env.PORT || 8080;

            // this.app.listen(PORT, () => {
            //     console.log(`App listening on port ${PORT}`);
            // });

            // // using http:
            const httpsOptions: https.ServerOptions = {
                secureOptions: constants.SSL_OP_NO_TLSv1 | constants.SSL_OP_NO_TLSv1_1,
                secureProtocol: 'TLSv1_2_server_method'
            };

            const server = http.createServer(httpsOptions, this.app);
            server.listen(process.env.PORT || 8080);

        }

    }

}