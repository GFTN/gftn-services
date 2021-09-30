// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const request = require('request')
const fs = require('fs')
const LOGGER = require('../method/logger')
const log = new LOGGER('Vault')
const checkResponse = require('../method/checkResponse')
const environment = require('../environment/env')
const accountstablename = process.env[environment.ENV_KEY_DYNAMODB_ACCOUNTS_TABLE_NAME]
const AWS = require('../method/AWS')


function getFromVault(url, AppID, Safe, Folder, Object, certPATH, keyPATH) {
  return new Promise(function (res, rej) {
    try {
      let options = {

        cert: fs.readFileSync(certPATH, 'utf8'),
        key: fs.readFileSync(keyPATH, 'utf8'),
        contentType: 'application/json',
        method: 'GET',
        rejectUnauthorized: false,
        json: true,
        url: url + '/AIMWebService/api/Accounts?AppID=' + AppID + '&Safe=' + Safe + '&Folder=' + Folder + '&Object=' + Object
      }
      request(options, function (err, response, body) {
        if (err) {
          throw ({
            statusCode: 500,
            Message: {
              ErrorMsg: err
            }
          })
        }
        else {
          log.logger('REQUEST', '')
          console.log(options);
          res(response)
          log.logger('RESPONSE :', '')
          console.log(body);
        }
      });

    } catch (err) {
      throw ({
        statusCode: 500,
        Message: {
          ErrorMsg: err
        }
      })
    }

  })


}

module.exports = {
  /**
   * 
   * @param {String} url 
   * @param {String} AppID 
   * @param {String} Safe 
   * @param {String} Folder 
   * @param {String} certPATH 
   * @param {String} keyPATH 
   * @param {String} accountInfos 
   */
  getDataFromVault: async function (url, AppID, Safe, Folder, certPATH, keyPATH, accountInfos) {

    let accountPromise = []
    for (let index = 0; index < accountInfos.length; index++) {

      let key = accountInfos[index].key.Object
      let seed = accountInfos[index].seed.Object
      let pkeyFromVault
      let secretFromVault

      let response = await getFromVault(url, AppID, Safe, Folder, key, certPATH, keyPATH)
      checkResponse(response)
      pkeyFromVault = response.body.Content

      response = await getFromVault(url, AppID, Safe, Folder, seed, certPATH, keyPATH)
      checkResponse(response)
      secretFromVault = response.body.Content

      let allAccounts = await AWS.getAllDatas(accountstablename)

      if (allAccounts.indexOf(pkeyFromVault) == -1) {

        let account = {
          pkey: pkeyFromVault,
          secret: secretFromVault,
          accountStatus: accountInfos[index].accountStatus,
          topicName: accountInfos[index].topicName
        }

        accountPromise.push(account)

      }

    }
    Promise.all(accountPromise)
    if (accountPromise.length == 0) {
      throw ({
        statusCode: 400,
        Message: {
          ErrorMsg: 'All accounts exist'
        }
      })
    }
    return accountPromise
  },
  requestVault: function getFromVault(url, AppID, Safe, Folder, Object, certPATH, keyPATH) {
    return new Promise(function (res, rej) {

      try {
        let options = {

          cert: fs.readFileSync(certPATH, 'utf8'),
          key: fs.readFileSync(keyPATH, 'utf8'),
          contentType: 'application/json',
          method: 'GET',
          rejectUnauthorized: false,
          json: true,
          url: url + '/AIMWebService/api/Accounts?AppID=' + AppID + '&Safe=' + Safe + '&Folder=' + Folder + '&Object=' + Object
        }
        request(options, function (err, response, body) {
          if (err) {
            throw ({
              statusCode: 500,
              Message: {
                ErrorMsg: err
              }
            })
          }
          else {
            res(response)
            console.log(body);
          }
        });


      } catch (err) {
        throw ({
          statusCode: 500,
          Message: {
            ErrorMsg: err
          }
        })
      }

    })


  }
}
