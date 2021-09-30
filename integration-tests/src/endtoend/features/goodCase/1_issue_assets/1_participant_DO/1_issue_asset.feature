
Feature: World Wire
  In order to helping a customer send money
  As a developer
  I want to create asset in Worldwire

  Scenario Outline: issue asset
    Given id: <participant_id> check participant: <participant_id>, was in Workd Wire, apiURL: <participant_api_url>
    And id: <participant_id> get account_name: <account_name> address from WW, participant_api_url:<participant_api_url>
    When id: <participant_id> issue asset asset_code: <participant_issue_asset>, asset_type: <participant_issue_asset_type> - apiURL: <participant_api_url>
    Then id: <participant_id> check id: <participant_id> has asset_code: <participant_issue_asset>, asset_type: <participant_issue_asset_type> in issued asset list - apiURL: <participant_api_url>
    And id: <participant_id> query id: <participant_id> has asset_code: <participant_issue_asset>, asset_type: <participant_issue_asset_type> in issued asset list - apiURL: <participant_api_url>


    Examples:
      | account_name                                 | participant_api_url             | participant_send_url             | participant_role             | participant_bic             | participant_country_code             | participant_id             | participant_issue_asset                    | participant_issue_asset_type             |
      | "ENV_KEY_PARTICIPANT_1_ISSUING_ACCOUNT_NAME" | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_PARTICIPANT_1_SEND_URL" | "ENV_KEY_PARTICIPANT_1_ROLE" | "ENV_KEY_PARTICIPANT_1_BIC" | "ENV_KEY_PARTICIPANT_1_COUNTRY_CODE" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_PARTICIPANT_1_ISSUE_ASSET"        | "ENV_KEY_PARTICIPANT_1_ISSUE_ASSET_TYPE" |
      | "ENV_KEY_PARTICIPANT_1_ISSUING_ACCOUNT_NAME" | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_PARTICIPANT_1_SEND_URL" | "ENV_KEY_PARTICIPANT_1_ROLE" | "ENV_KEY_PARTICIPANT_1_BIC" | "ENV_KEY_PARTICIPANT_1_COUNTRY_CODE" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_PARTICIPANT_1_ISSUE_REVOKE_ASSET" | "ENV_KEY_PARTICIPANT_1_ISSUE_ASSET_TYPE" |
      | "ENV_KEY_PARTICIPANT_2_ISSUING_ACCOUNT_NAME" | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_PARTICIPANT_2_SEND_URL" | "ENV_KEY_PARTICIPANT_2_ROLE" | "ENV_KEY_PARTICIPANT_2_BIC" | "ENV_KEY_PARTICIPANT_2_COUNTRY_CODE" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_PARTICIPANT_2_ISSUE_ASSET"        | "ENV_KEY_PARTICIPANT_2_ISSUE_ASSET_TYPE" |
