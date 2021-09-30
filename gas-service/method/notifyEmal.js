// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
// using SendGrid's v3 Node.js Library
// https://github.com/sendgrid/sendgrid-nodejs
const sgMail = require('@sendgrid/mail');
const LOGGER = require('../method/logger')
const log = new LOGGER('SendGrid')
const environment = require('../environment/env')
const contactsTableName = process.env[environment.ENV_KEY_DYNAMODB_CONTACTS_TABLE_NAME]
const AWS = require('../method/AWS')
const sender = process.env[environment.ENV_KEY_EMAIL_SENDER]
const sendGridAPIKey = process.env[environment.ENV_KEY_SENDGRID_API_KEY]
const sendGrudURL = process.env[environment.ENV_KEY_SENDGRIID_ENDPINT]
const request = require('request')
sgMail.setApiKey(sendGridAPIKey);

module.exports = async function (topicName, subject, text, html) {

  let groupArray = await AWS.getDataEmail(contactsTableName, topicName)
  if (groupArray.length>0){
    let sendList = await getList(groupArray)
    // await gridSemdEmail(sendGridAPIKey, sendList, subject, text, html)
    sendmail(sendList, subject, text, html)
  }
}


function getList(groupArray) {
  return new Promise(function (res, rej) {
    try {
      let p = []
      let sendList
      p.push(groupArray.forEach(async (info, index) => {
        if (index == 0) {
          sendList = info.email
        }
        else {
          sendList += ';' + info.email
        }
      }))

      Promise.all(p)
      res(sendList)
    }
    catch (err) {
      console.log(err);

      rej(err)
    }


  })
}

/**
 * ues api to send email
 */
function gridSemdEmail(sendGridAPIKey, sendList, subject, text, html) {
  return new Promise(function (resolve, reject) {
    let options = {
      headers: {
        "authorization": "Bearer " + sendGridAPIKey,
        "content-type": "application/json"
      },
      method: 'POST',
      body: {
        "personalizations": [
          { "to": [{ "email": sendList }], "subject": subject }
        ],
        "content": [
          { "type": "text/plain", "value": text }
        ],
        "from": { "email": sender }
      },
      json: true,
      url: sendGrudURL
    }

    log.info('Send email', JSON.stringify(options))

    request(options, function (err, res, body) {
      if (err) { console.log(err); reject(err) }
      else {
        log.info('Send email result: ', res.statusCode)
        resolve(res.statusCode)
      }
    });

  })
}

/**
 * sometime it is not working
 * ues sdk to send email
 */
function sendmail(sendList, subject, text, html) {
  const msg = {
    to: sendList,
    from: sender,
    subject: subject,
    text: text,
    html: html,
  };
  log.info('Send email', JSON.stringify(msg))
  sgMail.send(msg);

}
