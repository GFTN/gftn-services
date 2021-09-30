
Feature: World Wire goodCase/4_payment/DA/1_payment_comple
  In order to helping a customer send money
  As a developer
  I want to send payment in Worldwire

  Scenario Outline: 1_OFI_request_redemption_flow_RFI
    Given id: <participant_id> check participant: <RFI_participant_id>, was in Workd Wire, apiURL: <participant_api_url>
    When  id: <participant_id> check the asset_code: <sending_asset_code> asset_type: <sending_asset_type> issuer: <sending_asset_issuer> already issue by worldwire, apiURL:<participant_api_url>
    And id: <participant_id> add participant: <RFI_participant_id> - whitelistURL: <global_whitelist_url>
    And id: <participant_id> check participant: <RFI_participant_id> was in the whitelist - whitelistURL: <global_whitelist_url>
    And id: <RFI_participant_id> add participant: <participant_id> - whitelistURL: <global_whitelist_url>
    And id: <RFI_participant_id> check participant: <participant_id> was in the whitelist - whitelistURL: <global_whitelist_url>
    Then id: <participant_id> sign payload for redeem asset from sender_bic <sender_bic> sending_account_name <sending_account_name> amount <sending_amount> with settlement_method <settlement_method> asset_code <sending_asset_code> asset_issuer <asset_issuer_participant_id> to receiver <receiver_id> recever_bic <receiver_bic>, cryptoURL: <crypto_service_url>
    And id: <participant_id> sending type: <pac009file> signed payload, sendURL: <send_service_url>

    Examples:
      | participant_id             | participant_api_url             | RFI_participant_id  | asset_issuer_participant_id | sending_asset_code          | sending_asset_type          | sending_asset_issuer | global_whitelist_url           | send_service_url                 | sender_bic                  | sending_account_name                           | settlement_method              | sending_amount                         | receiver_id         | receiver_bic         | rfi_wwGatewayURL                      | crypto_service_url                 | pac009file                 |
      | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_ID" | "ENV_KEY_ANCHOR_ID"         | "ENV_KEY_ANCHOR_ASSET_CODE" | "ENV_KEY_ANCHOR_ASSET_TYPE" | "ENV_KEY_ANCHOR_ID"  | "ENV_KEY_PARTICIPANT_1_WL_URL" | "ENV_KEY_PARTICIPANT_1_SEND_URL" | "ENV_KEY_PARTICIPANT_1_BIC" | "ENV_KEY_PARTICIPANT_1_OPERATING_ACCOUNT_NAME" | "ENV_KEY_SETTLEMENT_METHOD_DA" | "ENV_KEY_PARTICIPANT_1_SENDING_AMOUNT" | "ENV_KEY_ANCHOR_ID" | "ENV_KEY_ANCHOR_BIC" | "ENV_KEY_PARTICIPANT_1_WWGATEWAY_URL" | "ENV_KEY_PARTICIPANT_1_CRYPTO_URL" | "iso20022:pacs.009.001.08" |
