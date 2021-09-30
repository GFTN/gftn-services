
Feature: World Wire goodCase/6_funding/1_anchor_funding_participant
  In order to helping a customer send money
  As a developer
  I want to issue anchor asset in WW

  Scenario Outline: anchor funding participant
    Given id: <anchor_id> check participant: <participant_id>, was in Workd Wire, anchorURL: <anchor_service_url>
    And id: <participant_id> check asset_code: <asset_code> issuer: <anchor_id> account_name: <account_name> balance before transaction, participant_api_url:<receive_participant_api_url>
    When id: <anchor_id> get instruction to fund account_name: <account_name> amount_funding: <amount_funding> anchor_id: <anchor_id> asset_code_issued: <asset_code_issued> end_to_end_id: <end_to_end_id> participant_id: <participant_id> memo_transaction: <memo_transaction>, anchor_service_url: <anchor_service_url>
    Then id: <anchor_id> signed instruction and details_funding to fund account_name: <account_name> amount_funding: <amount_funding> anchor_id: <anchor_id> asset_code_issued: <asset_code_issued> end_to_end_id: <end_to_end_id> participant_id: <participant_id> memo_transaction: <memo_transaction> with anchor_seed: <anchor_seed>, anchor_service_url: <anchor_service_url>
    And id: <participant_id> check asset_code: <asset_code> issuer: <anchor_id> account_name: <account_name> balance increase: <amount_funding>, participant_api_url:<receive_participant_api_url>


    Examples:
      | receive_participant_api_url     | anchor_service_url           | asset_code                    | asset_type                  | participant_role             | participant_bic             | participant_country_code             | participant_id             | account_name                                   | anchor_id           | amount_funding           | asset_code_issued             | end_to_end_id                   | memo_transaction       | anchor_seed           |
      | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_SERVICE_URL" | "ENV_KEY_ANCHOR_ASSET_CODE"   | "ENV_KEY_ANCHOR_ASSET_TYPE" | "ENV_KEY_PARTICIPANT_1_ROLE" | "ENV_KEY_PARTICIPANT_1_BIC" | "ENV_KEY_PARTICIPANT_1_COUNTRY_CODE" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_PARTICIPANT_1_OPERATING_ACCOUNT_NAME" | "ENV_KEY_ANCHOR_ID" | "ENV_KEY_FUNDING_AMOUNT" | "ENV_KEY_ANCHOR_ASSET_CODE"   | "ENV_KEY_FUNDING_END_TO_END_ID" | "ENV_KEY_FUNDINT_MEMO" | "ENV_KEY_ANCHOR_SEED" |
      | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_ANCHOR_SERVICE_URL" | "ENV_KEY_ANCHOR_ASSET_CODE"   | "ENV_KEY_ANCHOR_ASSET_TYPE" | "ENV_KEY_PARTICIPANT_2_ROLE" | "ENV_KEY_PARTICIPANT_2_BIC" | "ENV_KEY_PARTICIPANT_2_COUNTRY_CODE" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_PARTICIPANT_2_OPERATING_ACCOUNT_NAME" | "ENV_KEY_ANCHOR_ID" | "ENV_KEY_FUNDING_AMOUNT" | "ENV_KEY_ANCHOR_ASSET_CODE"   | "ENV_KEY_FUNDING_END_TO_END_ID" | "ENV_KEY_FUNDINT_MEMO" | "ENV_KEY_ANCHOR_SEED" |
      # | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_SERVICE_URL" | "ENV_KEY_ANCHOR_ASSET_CODE_2" | "ENV_KEY_ANCHOR_ASSET_TYPE" | "ENV_KEY_PARTICIPANT_1_ROLE" | "ENV_KEY_PARTICIPANT_1_BIC" | "ENV_KEY_PARTICIPANT_1_COUNTRY_CODE" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_PARTICIPANT_1_OPERATING_ACCOUNT_NAME" | "ENV_KEY_ANCHOR_ID" | "ENV_KEY_FUNDING_AMOUNT" | "ENV_KEY_ANCHOR_ASSET_CODE_2" | "ENV_KEY_FUNDING_END_TO_END_ID" | "ENV_KEY_FUNDINT_MEMO" | "ENV_KEY_ANCHOR_SEED" |
      # | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_ANCHOR_SERVICE_URL" | "ENV_KEY_ANCHOR_ASSET_CODE_2" | "ENV_KEY_ANCHOR_ASSET_TYPE" | "ENV_KEY_PARTICIPANT_2_ROLE" | "ENV_KEY_PARTICIPANT_2_BIC" | "ENV_KEY_PARTICIPANT_2_COUNTRY_CODE" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_PARTICIPANT_2_OPERATING_ACCOUNT_NAME" | "ENV_KEY_ANCHOR_ID" | "ENV_KEY_FUNDING_AMOUNT" | "ENV_KEY_ANCHOR_ASSET_CODE_2" | "ENV_KEY_FUNDING_END_TO_END_ID" | "ENV_KEY_FUNDINT_MEMO" | "ENV_KEY_ANCHOR_SEED" |
