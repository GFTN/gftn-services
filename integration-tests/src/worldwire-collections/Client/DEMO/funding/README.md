# WorldWire-Anchor-Funding

### HOW TO USE

NEED TO UPDATE THESE TO USE, OR IT WILL FAIL
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

### Execute
1. `npm install`
2. DO UPDATE INFORMATION
3. `node fundingWWapi.js`