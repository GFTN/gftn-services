// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
let AWS = require('../method/AWS')
const readFromAws = require('../feature/readFromAws')

module.exports = {
    createAccountsToDynamoDB:
    /**
     * 
     * @param {Array} accountInfos 
     * @param {String} accountstablename 
     * @param {String} url 
     * @param {String} AppID 
     * @param {String} Safe 
     * @param {String} Folder 
     * @param {String} certPATH 
     * @param {String} keyPATH 
     */
        async function(accountInfos, accountstablename, url, AppID, Safe, Folder, certPATH, keyPATH) {

        accounts = await readFromAws.getDataFromAws(accountInfos)

        let promise = []
        for (let index = 0; index < accounts.length; index++) {
            promise.push(await AWS.createItem(accountstablename, accounts[index]))

        }
        Promise.all(promise)
        await AWS.getAllDatas(accountstablename)
        return promise
    },
    /**
     * 
     * @param {*} contactsInfos 
     * @param {*} contactsTablename 
     * @param {*} topicstablename 
     */
    createContactsToDynamoDB: async function(contactsInfos, contactsTablename, topicstablename) {

        let promise = []

        for (let index = 0; index < contactsInfos.length; index++) {
            let Topicarn = await AWS.getTopicArn(topicstablename, contactsInfos[index].topicName)

            if (Topicarn == null) {
                throw ({
                    statusCode: 400,
                    Message: {
                        ErrorMsg: "Topic not exist"
                    }
                })
            } else {
                contactsInfos[index].Topicarn = Topicarn
            }
        }
        for (let index = 0; index < contactsInfos.length; index++) {
            let res = await AWS.subscribeTopic(contactsInfos[index].Topicarn, contactsInfos[index].phoneNumber)
                // console.log(res);
            contactsInfos[index].SubscriptionArn = res.SubscriptionArn
            promise.push(await AWS.createItem(contactsTablename, contactsInfos[index]))
        }
        Promise.all(promise)
        return promise


    },
    /**
     * 
     * @param {*} topicInfos 
     * @param {*} topicsTablename 
     */
    createTopicsToDynamoDB: async function(topicInfos, topicsTablename) {
        let promise = []

        for (let index = 0; index < topicInfos.length; index++) {
            /**
             * create topic to SNS
             * get the topic arn
             */
            let response = await AWS.createTopic(topicInfos[index].TopicName, topicInfos[index].DisplayName)

            let info = {
                TopicName: topicInfos[index].TopicName,
                displayName: topicInfos[index].DisplayName,
                TopicArn: response.TopicArn
            }
            promise.push(await AWS.createItem(topicsTablename, info))

        }
        Promise.all(promise)
        return promise
    }


}