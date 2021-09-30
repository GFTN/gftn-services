// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
var admin = require("firebase-admin");
var firebase = require('firebase');
const encoder = require('nodejs-base64-encode');
var firebaseCred = process.env.firebaseCred


let serviceAccount = JSON.parse(encoder.decode(firebaseCred, 'base64'))
var app = admin.initializeApp({
    credential: admin.credential.cert(serviceAccount),
    databaseURL: process.env.databaseURL
});

var app = firebase.initializeApp({
    apiKey: process.env.firebaseAPIKey,
    authDomain: process.env.authDomain,
    databaseURL: process.env.databaseURL,
    projectId: process.env.projectId,
    storageBucket: process.env.storageBucket,
    messagingSenderId: process.env.messagingSenderId,
    appId: process.env.appId,
});


function setFirebaseUser(email) {
    return new Promise((resolve, reject) => {

        admin.auth().getUserByEmail(email)
            .then(function(userRecord) {
                // See the UserRecord reference doc for the contents of userRecord.
                // console.log('Successfully fetched user data:', userRecord.toJSON());
                resolve(userRecord);
            })
            .catch(function(error) {
                console.log('Error fetching user data:', error);
                reject(error);
            });
    })
}



function createCustomToken(userRecord) {
    return new Promise((resolve, reject) => {
        admin.auth().createCustomToken(userRecord.uid, {})
            .then(function(customToken) {
                // Send token back to client
                // console.log(customToken);
                resolve(customToken);

            })
            .catch(function(error) {
                console.log('Error creating custom token:', error);
                reject(error);
            });
    })
}


function getFID(token) {
    return new Promise((resolve, reject) => {
        firebase.auth().signInWithCustomToken(token)
            .then(function() {
                firebase.auth().onAuthStateChanged(function(user) {
                    // console.log(user);
                    if (user) {
                        user.getIdToken().then(function(data) {
                            // console.log(data)
                            resolve(data);
                        });
                    }
                });
            })
            .catch(function(error) {

                // Handle Errors here.
                var errorCode = error.code;
                var errorMessage = error.message;
                if (errorCode === 'auth/invalid-custom-token') {
                    alert('The token you provided is not valid.');
                } else {
                    console.log('The token you provided is not valid:', error);
                    reject(error)
                }
            });
    })
}


function getTotp(key) {
    var token = speakeasy.totp({
        secret: key,
        encoding: 'ascii'
    });

    console.log(token);
}

module.exports = async(email) => {
    try {
        let userRecord = await setFirebaseUser(email)
        let customToken = await createCustomToken(userRecord)
        const token = await getFID(customToken)
        return token
    } catch (error) {
        console.log('Error when generate firebaseID:' + error);
    }

}