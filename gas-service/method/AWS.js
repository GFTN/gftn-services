// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 

const AWS = require('aws-sdk')
const LOGGER = require('./logger');
const environment = require('../environment/env')
let log = new LOGGER('AWS Call')
AWS.config.update({
    region: process.env[environment.ENV_KEY_DYNAMODB_REGION],
    accessKeyId: process.env[environment.ENV_KEY_DYNAMODB_ACCESSKEYID],
    secretAccessKey: process.env[environment.ENV_KEY_DYNAMODB_SECRECT_ACCESS_KEY]
});
var dynamodb = new AWS.DynamoDB();
var docClient = new AWS.DynamoDB.DocumentClient();

module.exports = {

    createTable:
        /**
         * 
         * @param {Object} create table content 
         */
        function (params) {
            return new Promise(function (res, rej) {
                dynamodb.createTable(params, function (err, data) {
                    if (err) {
                        res(err)
                    } else {
                        res(data)
                    }
                });


            })
        },

    deleteTable:
        /**
         * 
         * @param {String} tablename 
         */
        function (accountstablename) {
            return new Promise(function (res, rej) {

                var params = {
                    TableName: accountstablename
                };
                dynamodb.deleteTable(params, function (err, data) {
                    if (err) {
                        rej(err)
                    } else {
                        res(data)
                    }
                });
            })
        },

    createItem:
        function (tablename, item) {
            return new Promise(function (res, rej) {

                var params = {
                    TableName: tablename,
                    Item: item
                };

                docClient.put(params, function (err, data) {
                    if (err) {
                        throw err
                    } else {
                        res(item)
                    }
                });


            })
        },
    deleteItem:
        function (params) {
            return new Promise(function (res, rej) {
                dynamodb.deleteItem(params, function (err, data) {
                    console.log(err);
                    console.log(data);
                                        
                    if (err) {
                        throw err
                    }
                    else {
                        res(data)
                    }
                });

            })
        },
    updateItem:
        function (params) {
            return new Promise(function (res, rej) {
                docClient.update(params, function (err, data) {
                    if (err) {
                        throw ({
                            statusCode: 500,
                            Message: {
                                ErrorMsg: err
                            }
                        })
                    } else {
                        res(data)
                    }
                });

            })
        },
    getAccountsInfo:
        async function (accountstablename, status, getTS) {
            return new Promise(function (res, rej) {
                try {
                    
                    let accountArray = []
                    let params = {
                        TableName: accountstablename,
                        ProjectionExpression: "#st,#ts,pkey, secret",
                        FilterExpression: "#st = :status ",
                        ExpressionAttributeNames: {
                            "#st": "accountStatus",
                            "#ts": "lockTimestamp"
                        },
                        ExpressionAttributeValues: {
                            ":status": status
                        }

                    };
                    docClient.scan(params, onScan);

                    function onScan(err, data) {
                        if (err) {
                            console.log(err);
                            throw ({
                                statusCode: 500,
                                Message: {
                                    ErrorMsg: err
                                }
                            })
                        } else {
                            if (getTS) {
                                data.Items.forEach(function (item) {
                                    let obj = {
                                        pkey: item.pkey,
                                        lockTimestamp: item.lockTimestamp
                                    }
                                    accountArray.push(obj)
                                })
                                res(accountArray)

                            }
                            else {
                                data.Items.forEach(function (item) {
                                    accountArray.push(item.pkey)
                                })

                                res(accountArray)
                            }

                        }
                    }

                } catch (err) {
                    throw ({
                        statusCode: 500,
                        Message: {
                            ErrorMsg: err
                        }
                    })
                }

            })

        },

    queryData:
    function (params) {
        return new Promise(function (res, rej) {
            docClient.query(params, function (err, data) {
                if (err) {
                    res(err)
                } else {
                    if (data.Items.length == 0) {
                        res(null)
                    }
                    else {
                        res(data.Items[0])
                    }
                }
            });

        })
    },
    // queryDataSecret:
    //     function (accountstablename, pkey) {
    //         return new Promise(function (res, rej) {
    //             var params = {
    //                 TableName: accountstablename,
    //                 KeyConditionExpression: "#pk = :pkey",
    //                 ExpressionAttributeNames: {
    //                     "#pk": "pkey"
    //                 },
    //                 ExpressionAttributeValues: {
    //                     ":pkey": pkey
    //                 }
    //             };

    //             docClient.query(params, function (err, data) {
    //                 if (err) {
    //                     res(err)
    //                 } else {
    //                     if (data.Items.length == 0) {
    //                         res(null)
    //                     }
    //                     else {
    //                         res(data.Items[0].secret)
    //                     }
    //                 }
    //             });

    //         })
    //     },
    // queryDataGroupID:
    //     function (accountstablename, pkey) {
    //         return new Promise(function (res, rej) {

    //             var params = {
    //                 TableName: accountstablename,
    //                 KeyConditionExpression: "#pk = :pkey",
    //                 ExpressionAttributeNames: {
    //                     "#pk": "pkey"
    //                 },
    //                 ExpressionAttributeValues: {
    //                     ":pkey": pkey
    //                 }
    //             };

    //             docClient.query(params, function (err, data) {
    //                 if (err) {
    //                     res(err)
    //                 } else {
    //                     res(data.Items[0].groupName)
    //                 }
    //             });

    //         })
    //     },
    getDataEmail:
        /**
         * 
         * @param {String} tablename 
         * @param {Bool} status 
         * @param {Array} lockAccounts 
         */
        function (contactsTableName, topicName) {
            return new Promise(function (res, rej) {
                let Array = []

                var params = {
                    TableName: contactsTableName,
                    ProjectionExpression: "#topicName,email, phoneNumber",
                    FilterExpression: "topicName = :topicName ",
                    ExpressionAttributeNames: {
                        "#topicName": "topicName"
                    },
                    ExpressionAttributeValues: {
                        ":topicName": topicName
                    }

                };
                
                docClient.scan(params, onScan);

                function onScan(err, data) {
                    if (err) {
                        rej(err)
                    } else {
                        res(data.Items)
                    }
                }

            })

        },
    getSubscriptionArn:
        function (contactstablename, topicName, email) {
            return new Promise(function (res, rej) {

                var params = {
                    TableName: contactstablename,
                    ProjectionExpression: "#topicName, #email, phoneNumber,Topicarn,SubscriptionArn",
                    FilterExpression: "topicName = :topicName and email = :email",
                    ExpressionAttributeNames: {
                        "#topicName": "topicName",
                        "#email": "email"
                    },
                    ExpressionAttributeValues: {
                        ":topicName": topicName,
                        ":email": email
                    }

                };

                docClient.scan(params, onScan);
                function onScan(err, data) {
                    if (err) {
                        throw ({
                            statusCode: 500,
                            Message: {
                                ErrorMsg: err
                            }
                        })
                    } else {
                        if (data.Items.length > 0) res(data.Items[0].SubscriptionArn)
                        else {
                            res(null)
                        }
                    }
                }

            })

        },
    getTopicArn:
        function (topicstablename, topicName) {
            return new Promise(function (res, rej) {

                var params = {
                    TableName: topicstablename,
                    ProjectionExpression: "#TopicName,TopicArn, displayName",
                    FilterExpression: "TopicName = :TopicName ",
                    ExpressionAttributeNames: {
                        "#TopicName": "TopicName"
                    },
                    ExpressionAttributeValues: {
                        ":TopicName": topicName
                    }

                };

                docClient.scan(params, onScan);

                function onScan(err, data) {
                    if (err) {
                        throw ({
                            statusCode: 500,
                            Message: {
                                ErrorMsg: err
                            }
                        })
                    } else {
                        if (data.Items.length > 0) res(data.Items[0].TopicArn)
                        else {
                            res(null)
                        }
                    }
                }

            })

        },

    getAllDatas:
        /**
         * 
         * @param {String} tablename 
         */
        function (accountstablename) {
            return new Promise(function (res, rej) {
                var params = {
                    TableName: accountstablename
                };
                docClient.scan(params, onScan);

                function onScan(err, data) {
                    if (err) {
                        res(err)
                    } else {
                        res(data.Items)

                    }
                }
            })

        },
    unsubscribe:
        function (SubscriptionArn) {
            return new Promise(function (res, rej) {

                var params = {
                    SubscriptionArn: SubscriptionArn /* required */
                };
                var sns = new AWS.SNS();
                var request = sns.unsubscribe(params);

                console.log(params);
                request.
                    on('success', function (response) {
                        
                    }).
                    on('error', function (response) {
                        
                    }).
                    on('complete', function (response) {
                        
                        res(response.data)
                    }).
                    send();
            })
        },
    createTopic:
        function (topicName, displayName) {
            return new Promise(function (res, rej) {


                var params = {
                    Name: topicName, /* required */
                    Attributes: {
                        'DisplayName': displayName
                    }
                };
                var sns = new AWS.SNS();
                var request = sns.createTopic(params);

                request.
                    on('success', function (response) {
                    }).
                    on('error', function (response) {
                    }).
                    on('complete', function (response) {
                        res(response.data)
                    }).
                    send();
            })
        },
    deleteTopic: function (TopicArn) {

        return new Promise(function (res, rej) {
            var params = {
                TopicArn: TopicArn /* required */
            };

            var sns = new AWS.SNS();
            var request = sns.deleteTopic(params);

            request.
                on('success', function (response) {
                }).
                on('error', function (response) {

                }).
                on('complete', function (response) {

                    res(response.data)
                }).
                send();
        })
    },
    subscribeTopic:
        function (TopicArn, phoneNumber) {
            return new Promise(function (res, rej) {
                // Create publish parameters
                var params = {
                    Protocol: 'sms', /* required */
                    TopicArn: TopicArn.toString(),
                    Endpoint: phoneNumber.toString(),
                    ReturnSubscriptionArn: true
                };

                var sns = new AWS.SNS();
                let request = sns.subscribe(params);

                request.
                    on('success', function (response) {
                        // console.log("Success!");
                    }).
                    on('error', function (response) {

                        // console.log("Error!");
                    }).
                    on('complete', function (response) {
                        // console.log("Always!");
                        res(response.data)
                    }).
                    send();
            })
        },
    sendSMS:
        function (msg, topic) {

            return new Promise(function (res, rej) {
                // Create publish parameters
                var params = {
                    Message: msg, /* required */
                    TopicArn: topic.toString()
                };

                // Create promise and SNS service object
                var publishTextPromise = new AWS.SNS({ apiVersion: '2010-03-31' }).publish(params).promise();

                // Handle promise's fulfilled/rejected states
                publishTextPromise.then(
                    function (data) {
                        res(`Message : ${params.Message} send sent to the topic ${params.TopicArn}` +" MessageID is " + data.MessageId )
                    }).catch(
                        function (err) {
                            rej(err)// console.error(err, err.stack);
                        });
            })
        }


}