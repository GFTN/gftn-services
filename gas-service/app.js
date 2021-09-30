// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
var express = require('express');
var app = express();
var bodyParser = require('body-parser');
var bigInt = require("big-integer");
const LOGGER = require('./method/logger')
const log = new LOGGER('API Call')
const environment = require('./environment/env')
let maxRetry = 5

let timer = setInterval(() => {

    if (maxRetry <= 0) {
        process.exit(1)
    }

    if (!process.env[environment.ENV_KEY_GAS_SERVICE_PORT]) {
        log.error("Error reading environment variables", "retrying...")
        maxRetry--
        return
    }

    const AWS = require('./method/AWS')
    const updater = require('./method/updater')
    const logFile = require('./method/logFile')
    const lockaccount = require("./feature/lockAccount.js")
    const unlockAccount = require("./feature/unlockAccount")
    const createTXE = require("./feature/createTXE.js")
    const signXDRAndExecuteXDR = require("./feature/signXDRAndExecuteXDR.js")
    const monitor = require("./feature/monitorLockAccounts.js")
    const monitorBalancesH = require('./feature/monitorBalancesH')
    const deleteTopic = require('./feature/deleteTopic')
    const deleteContact = require('./feature/deleteContact')
    const register = require('./feature/register')


    const AppID = process.env[environment.ENV_KEY_VAULT_APPID]
    const Safe = process.env[environment.ENV_KEY_VAULT_SAFE]
    const Folder = process.env[environment.ENV_KEY_VAULT_FOLDER]
    const vaultUrl = process.env[environment.ENV_KEY_VAULT_URL]
    const certPATH = process.env[environment.ENV_KEY_VAULT_CERT_PATH]
    const keyPATH = process.env[environment.ENV_KEY_VAULT_KEY_PATH]
    const accountstablename = process.env[environment.ENV_KEY_DYNAMODB_ACCOUNTS_TABLE_NAME]
    const contactstablename = process.env[environment.ENV_KEY_DYNAMODB_CONTACTS_TABLE_NAME]
    const topicstablename = process.env[environment.ENV_KEY_DYNAMODB_GROUPS_TABLE_NAME]

    const highThresholdBalance = process.env[environment.ENV_KEY_HIGH_THRESHOLD_BALANCE]
    const highThresholdTimeout = process.env[environment.ENV_KEY_HIGH_THRESHOLD_TIMEOUT]
    const lowThresholdBalance = process.env[environment.ENV_KEY_LOW_THRESHOLD_BALANCE]
    const lowThresholdTimeout = process.env[environment.ENV_KEY_LOW_THRESHOLD_TIMEOUT]
    let lockStatus = false
    let unlockStatus = true
    let lockAccounts = []
    let unlockAccounts
        // let highThresholdAccounts=[]
    let highThresholdAccounts = []
    let lowThresholdAccounts


    /**
     * check lock per lockAccountsMonitorTime
     * if lock timestamp + expireTime > now
     * unlock
     */
    let lockAccountsMonitorTime = process.env[environment.ENV_KEY_GAS_SERVICE_MONITOR_LOCKACCOUNT_FEQ]
    let expireTime = process.env[environment.ENV_KEY_GAS_SERVICE_EXPIRE_TIME]


    app.use(bodyParser.json()); // for parsing application/json
    app.use(bodyParser.urlencoded({
        extended: true
    })); // for parsing application/x-www-form-urlencoded
    //set header
    app.all('*', function(req, res, next) {

        res.header("Access-Control-Allow-Origin", "*");
        res.header("Access-Control-Allow-Headers", "X-Requested-With");
        res.header("Access-Control-Allow-Methods", "PUT,POST,GET,DELETE,OPTIONS");
        res.header("X-Powered-By", ' 3.2.1');
        res.header("Content-Type", "application/json;charset=utf-8");
        next();
    });

    app.post('/createAccounts', async function(req, res) {
        try {
            log.logger("/createAccounts", "REQUEST")
            console.log(req.body);


            let newAccounts = await register.createAccountsToDynamoDB(req.body, accountstablename, vaultUrl, AppID, Safe, Folder, certPATH, keyPATH)
            await updater.updateUnlockQueue(unlockAccounts, newAccounts)

            await updater.addHighThresholdAccountsQueue(highThresholdAccounts, newAccounts)
            res.status(200).json(newAccounts)

        } catch (err) {
            console.log(err);
            if (typeof err.statusCode !== 'undefined') {
                res.status(err.statusCode).json(err.Message)
            } else {
                res.status(500).json(err)
            }

        }
    });
    app.post('/createContacts', async function(req, res) {
        try {
            log.logger("/createContacts", "REQUEST")
            console.log(req.body);
            let result = await register.createContactsToDynamoDB(req.body, contactstablename, topicstablename)
            res.status(200).json(result)
        } catch (err) {
            log.error("/createContacts", JSON.stringify(err, 2, null))
            if (typeof err.statusCode !== 'undefined') {
                res.status(err.statusCode).json(err.Message)
            } else {
                res.status(500).json(err)
            }
        }

    });

    app.post('/createTopics', async function(req, res) {

        try {
            log.logger("/createTopics", "REQUEST")
            console.log(req.body);
            let result = await register.createTopicsToDynamoDB(req.body, topicstablename)

            res.status(200).json(result)

        } catch (err) {

            res.status(500).json(err)
        }

    });


    app.delete('/deleteTopic', async function(req, res) {

        try {
            let result = await deleteTopic(topicstablename, req.body.TopicName, req.body.TopicArn)
            res.status(200).json(result)

        } catch (err) {

            res.status(500).json(err)
        }

    });


    app.delete('/deleteContact', async function(req, res) {

        try {
            let result = await deleteContact(contactstablename, req.body.topicName, req.body.email)
            res.status(200).json(result)

        } catch (err) {
            res.status(err.statusCode).json(err.Message)
        }

    });

    app.delete('/deleteAccount', async function(req, res) {

        try {
            let result = await AWS.deleteAccount(accountstablename, req.body)
            res.status(200).json(result)

        } catch (err) {
            res.status(500).json(err)
        }

    });


    app.get('/getTopics', async function(req, res) {

        try {
            let result = await AWS.getAllDatas(topicstablename)
            res.status(200).json(result)

        } catch (err) {
            res.status(500).json(err)
        }

    });

    app.get('/getContacts', async function(req, res) {

        try {
            let result = await AWS.getAllDatas(contactstablename)
            res.status(200).json(result)

        } catch (err) {
            res.status(500).json(err)
        }

    });


    app.post('/unlockAccount', async function(req, res) {

        try {
            log.logger("/unlockAccount", "REQUEST")
            console.log(req.body);
            let account = {
                accountStatus: true,
                pkey: req.body.pkey
            }
            let result = await unlockAccount(account, accountstablename, lockAccounts, unlockAccounts)
            res.status(200).json(result)
        } catch (err) {
            console.log(err);

            if (typeof err.statusCode !== 'undefined') {
                res.status(err.statusCode).json(err.Message)
            } else {
                res.status(500).json(err)
            }
        }

    });


    app.get('/getMockTx', async function(req, res) {
        log.logger("/getMockTx", "REQUEST")
        console.log(req.body);
        let signedXDRin
        try {
            if (typeof req.body.sequenceNumber == 'undefined') {
                throw ({
                    Message: "body lost sequenceNumber"
                })
            }
            signedXDRin =
                await createTXE(req.body.sequenceNumber,
                    req.body.from.pkey,
                    req.body.from.secret,
                    req.body.to.pkey,
                    req.body.from.asset,
                    req.body.to.asset,
                    req.body.from.asset.amount,
                    req.body.to.asset.amount)

            res.status(200).json({
                oneSignedXDR: signedXDRin
            })

        } catch (err) {
            log.error('createTXE', err)
            res.status(400).json(err)
        }
    });

    app.get('/lockaccount', async function(req, res) {
        log.logger("/lockaccount", "REQUEST")
        console.log(req.body);

        let result = await lockaccount(lockAccounts, unlockAccounts, accountstablename)
        log.logger("lockaccount() ", 'return')
        console.log(result)


        let response = {}
        if (result == null) {
            response.failure_reason = "no avaible account now"
            res.status(500).json(response)

        } else {
            response.pkey = result.pkey
            response.sequenceNumber = bigInt(result.sequenceNumber).add(1).toString()
            log.logger("sequencenum:", response.sequenceNumber)
            res.status(200).json(response)
        }

    });



    app.post('/signXDRAndExecuteXDR', async function(req, res) {
        log.logger("/signXDRAndExecuteXDR", "REQUEST")
        console.log(req.body);
        try {

            let result = await signXDRAndExecuteXDR(req.body.oneSignedXDR, lockAccounts, unlockAccounts, accountstablename)


            if (result.title == "Source Account Expire") {
                res.status(400).json(result)
            }
            if (result.title == "Transaction Failed") {
                res.status(403).json(result)
            }
            if (result.title == "Transaction successful") {
                res.status(200).json(result)
            }

        } catch (err) {
            if (typeof err.statusCode !== 'undefined') {
                res.status(err.statusCode).json(err.Message)
            } else {
                log.error(err)
                res.status(500).json(err)
            }
        }

    });




    //service port
    var server = app.listen(process.env[environment.ENV_KEY_GAS_SERVICE_PORT], async function() {
        try {

            lockAccounts = await AWS.getAccountsInfo(accountstablename, lockStatus, true)
            unlockAccounts = await AWS.getAccountsInfo(accountstablename, unlockStatus, false)


            monitor(lockAccountsMonitorTime, expireTime, lockAccounts, unlockAccounts, accountstablename)
            let accounts = await AWS.getAllDatas(accountstablename)

            accounts.forEach(function(item) {
                highThresholdAccounts.push(item.pkey)

            })

            lowThresholdAccounts = []

            monitorBalancesH(highThresholdAccounts, highThresholdBalance, highThresholdTimeout, lowThresholdAccounts, lowThresholdBalance, lowThresholdTimeout)

            var host = server.address().address;
            var port = server.address().port;
            console.log('Gas service app listening at http://%s:%s', host, port);
            logFile()

        } catch (err) {

            log.error('ERROR', JSON.stringify(err, undefined, 2))

        }

    })
    clearInterval(timer)
}, 1000)