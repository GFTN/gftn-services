// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

const AWS = require('../method/AWS')
const LOGGER = require('../method/logger')
const log = new LOGGER('Delete Call')

module.exports =
  /**
   * 
   * @param {String} signedXDRin 
   * @param {String} accountstablename 
   */
  async function (groupstablename, TopicName,TopicArn) {
    
    let result = await AWS.deleteTopic(TopicArn)
    let param = params = {
      Key: {
          "TopicName": {
              S: TopicName
          },
          "TopicArn": {
              S: TopicArn
          }
      },
      TableName: groupstablename
  };
    await AWS.deleteItem(param)
    return result
  }
