// © Copyright IBM Corporation 2020. All rights reserved.
// SPDX-License-Identifier: Apache2.0
// 
const request = require('request')
const should = require('should');
const StellarSdk = require('stellar-sdk');
// StellarSdk.Network.use(new StellarSdk.Network('Standalone World Wire Network ; Mar 2019'));

// console.log(StellarSdk.Network.networkPassphrase());



/*

Impotant !!!!!!!
NEED TO CHANGE THESE TO USE, OR IT WILL FAIL
1. UPDATE URL 
2. UPDATE ANCHOR SEED 
3. UPDATE FUNDING OBJ PARTICIPANT ID 
4. UPDATE FUNDING ACCOUNT NAME
5. UPDATE FUNDING JWT-TOKEN

calling fundingParticipantAcc(fundingObject,anchorSVCurl,anchorSeed,anchorToken) to funding


{
    account_name: "OPERATING ACCOUNT NAME",
    amount_funding: FUNDING AMOUNT,
    anchor_id: "ANCHOR ID",
    asset_code_issued: "ASSET CODE",
    end_to_end_id: "CAN BE RANDOM NUMBER",
    participant_id: "TARGET PARTICIPANT ID",
    memo_transaction: "DESCRIPTION"
}
*/



let fundingObj = {
    account_name: "default",
    amount_funding: 499,
    anchor_id: "ibmtaiwan1",
    asset_code_issued: "USD",
    end_to_end_id: "32939",
    participant_id: "test1",
    memo_transaction: "worldwire funding operating account"
}
let anchorUrl = 'https://ibmtaiwan1.worldwire-qa.io/global/anchor/v1'
let anchorSeed = 'SAJHQPYYZIS6GYFVFSTPOMAVGJDYKXDABXXOXG67TDUNTYIIGJRNLZKL'
let anchorToken = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6Ii1McFJBNTlSTzB3OU9NZHhtLUUwLjEtMjgifQ.eyJhY2MiOlsiaXNzdWluZyIsImRlZmF1bHQiXSwidmVyIjoiMi45LjMuMTNfUkMxIiwiaXBzIjpbIjIwMi4xMzUuMjQ1LjM5IiwiMjAyLjEzNS4yNDUuNCIsIjIxOS43NC4xNS4yMDciLCIyMDIuMTM1LjI0NS4yIiwiMTI1LjE4LjkuMjAiLCIxNjkuNjMuMTQwLjEwNiIsIjIwMi4xMzUuMjQ1LjM1IiwiMzYuMjI2LjI1MS4yMjAiLCIzOS4xMi4yMDIuMTUiXSwiZW52Ijoic3QiLCJlbnAiOlsiL3YxL2FkbWluL3ByIiwiL3YxL2FkbWluL3ByL2RvbWFpbiIsIi92MS9hbmNob3IvYWRkcmVzcyIsIi92MS9hbmNob3IvYXNzZXRzL2lzc3VlZCIsIi92MS9hbmNob3IvZnVuZGluZ3MvaW5zdHJ1Y3Rpb24iLCIvdjEvYW5jaG9yL2Z1bmRpbmdzL3NlbmQiLCIvdjEvYW5jaG9yL3RydXN0IiwiL3YxL2NsaWVudC9hY2NvdW50cyIsIi92MS9jbGllbnQvYXNzZXRzIiwiL3YxL2NsaWVudC9hc3NldHMvYWNjb3VudHMiLCIvdjEvY2xpZW50L2Fzc2V0cy9pc3N1ZWQiLCIvdjEvY2xpZW50L2Fzc2V0cy9wYXJ0aWNpcGFudHMiLCIvdjEvY2xpZW50L2JhbGFuY2VzL2FjY291bnRzIiwiL3YxL2NsaWVudC9iYWxhbmNlcy9kbyIsIi92MS9jbGllbnQvZXhjaGFuZ2UiLCIvdjEvY2xpZW50L2ZlZXMvcmVxdWVzdCIsIi92MS9jbGllbnQvZmVlcy9yZXNwb25zZSIsIi92MS9jbGllbnQvbWVzc2FnZSIsIi92MS9jbGllbnQvcGFydGljaXBhbnRzIiwiL3YxL2NsaWVudC9wYXJ0aWNpcGFudHMvd2hpdGVsaXN0IiwiL3YxL2NsaWVudC9wYXlvdXQiLCIvdjEvY2xpZW50L3F1b3RlcyIsIi92MS9jbGllbnQvcXVvdGVzL3JlcXVlc3QiLCIvdjEvY2xpZW50L3NpZ24iLCIvdjEvY2xpZW50L3Rva2VuL3JlZnJlc2giLCIvdjEvY2xpZW50L3RyYW5zYWN0aW9ucyIsIi92MS9jbGllbnQvdHJhbnNhY3Rpb25zL3JlcGx5IiwiL3YxL2NsaWVudC90cmFuc2FjdGlvbnMvc2VuZCIsIi92MS9jbGllbnQvdHJhbnNhY3Rpb25zL3NldHRsZS9kYSIsIi92MS9jbGllbnQvdHJ1c3QiLCIvdjEvb25ib2FyZGluZy9hY2NvdW50cyJdLCJuIjowLCJpYXQiOjE1NjkyMDk0NzUsIm5iZiI6MTU2OTIwOTQ4MCwiZXhwIjoxNTY5NDY4Njc1LCJhdWQiOiJpYm1hbmNob3IiLCJzdWIiOiItTG9ESDQwLXhVRjJUV29zazRTNiIsImp0aSI6Ii1McFI5dXNGd3AwbnFIbC1IeThQIn0.H5HrNUmb4KX9mvF-m58O3SCYW2pO01y4GKvFme0RILM'

fundingParticipantAcc(fundingObj, anchorUrl, anchorSeed, anchorToken)




async function fundingParticipantAcc(fundingObj, anchorSVCurl, anchorSeed, anchorToken) {
    // Get Instruction From WorldWire
    let detailNinstr = await postFundingsInst(fundingObj, anchorSVCurl, anchorToken)
        // Using anchor seed to sign 
    let signedFundingNInstr = await signFundingInstr(anchorSeed, detailNinstr)
        // Send request to WW, then WW will send Tx to stellar via gas service
    await postFundingSend(fundingObj, signedFundingNInstr, anchorSVCurl, anchorToken)

}


function postFundingsInst(fundingObj, url, anchorToken) {
    return new Promise(function(res, rej) {
        var options = {
            method: 'POST',
            url: url + '/anchor/fundings/instruction',
            headers: {
                Authorization: 'Bearer ' + anchorToken
            },
            body: fundingObj,
            json: true
        };
        console.log('======================= Sending request : =============================');
        console.log('Sending request to : ' + options.url);
        console.log('Request Body: ' + JSON.stringify(options.body));


        // Get Instruction From WorldWire
        request(options, function(err, response, body) {
            if (err) {
                rej(err)
            } else {
                should(response.statusCode).be.exactly(200, "response data: " + JSON.stringify(response.body) + "\n request data : " + JSON.stringify(options))
                let details_funding = JSON.stringify(body.details_funding)
                let instruction_unsigned = body.instruction_unsigned
                console.log('=============== Get Instruction From WorldWire: ===============');
                console.log('Response status code: ' + response.statusCode);
                console.log('Response body: ' + JSON.stringify(body));
                res([details_funding, instruction_unsigned])
            }
        });

    })
}

// Using anchor seed to sign 
function signFundingInstr(anchorSeed, detailNinstr) {

    return new Promise(function(res, rej) {

        const source = StellarSdk.Keypair.fromSecret(anchorSeed)
        console.log('=============== Using anchor public key and seed: ===============');
        console.log('pkey: ' + source.publicKey());
        console.log('seed: ' + source.secret());

        let details_funding = detailNinstr[0]
        let instruction_unsigned = detailNinstr[1]

        let afrString = details_funding
        let buf = Buffer.from(afrString, 'ascii');
        let base64afrString = buf.toString('base64')
        let signereq = source.sign(base64afrString)
        let funding_signed = signereq.toString('base64')

        let transaction = new StellarSdk.Transaction(instruction_unsigned, 'Standalone World Wire Network ; Mar 2019')
        transaction.sign(source)
        let instruction_signed = transaction.toEnvelope().toXDR('base64')

        console.log('=============== Signed Funding And Instruction: ===============');
        console.log('funding_signed: ' + funding_signed);
        console.log('instruction_signed: ' + instruction_signed);

        res([funding_signed, instruction_signed])
    })
}

// Send request to WW, then WW will send Tx to stellar 
function postFundingSend(fundingObj, signedFundingNInstr, anchorSVCurl, anchorToken) {
    return new Promise(function(res, rej) {
        var options = {
            method: 'POST',
            url: anchorSVCurl + '/anchor/fundings/send',
            qs: {
                funding_signed: encodeURI(signedFundingNInstr[0]),
                instruction_signed: encodeURI(signedFundingNInstr[1])
            },
            headers: {
                Authorization: 'Bearer ' + anchorToken
            },
            body: fundingObj,
            json: true
        };

        console.log('Sending request to : ' + options.url);
        console.log('Request body: ' + JSON.stringify(options.body));

        request(options, function(err, response, body) {
            if (err) {
                rej([option, err])
            } else {
                should(response.statusCode).be.exactly(200, "response data: " + JSON.stringify(response.body) + "\n request data : " + JSON.stringify(options))

                console.log('=============== Funding via WorldWire Result: ===============');
                console.log('Response status code: ' + response.statusCode);
                console.log('Response body: ' + JSON.stringify(body));
                console.log('===============================================================');
                res(body)
            }
        });
    })
}





// module.exports = async function fundingParticipantAcc(fundingObj, anchorSVCurl, anchorSeed) {
//     // Get Instruction From WorldWire
//     let detailNinstr = await postFundingsInst(fundingObj, anchorSVCurl,anchorToken)
//         // Using anchor seed to sign 
//     let signedFundingNInstr = await signFundingInstr(anchorSeed, detailNinstr)
//         // Send request to WW, then WW will send Tx to stellar via gas service
//     await postFundingSend(fundingObj, signedFundingNInstr, anchorSVCurl,anchorToken)

// }