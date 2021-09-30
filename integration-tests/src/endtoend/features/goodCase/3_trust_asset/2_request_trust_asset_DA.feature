
Feature: World Wire
  In order to helping a customer send money
  As a developer
  I want to create trust line in Worldwire

  Scenario Outline: request trust asset
    # Given id: <participant_id> check participant: <participant_trust_asset_issuer>, was in Workd Wire, apiURL: <participant_api_url>
    When id: <participant_id> check the asset_code: <participant_trust_asset_code> asset_type: <participant_trust_asset_type> issuer: <participant_trust_asset_issuer> already issue by worldwire, apiURL:<participant_api_url>
    And id: <participant_id> check participant: <participant_trust_asset_issuer> was in the whitelist - whitelistURL: <global_whitelist_url>
    Then id: <participant_id> send trust request to trust asset_code: <participant_trust_asset_code> asset_issuer: <participant_trust_asset_issuer> with trust_limit <participant_trust_asset_limit> in trust_account <participant_trust_asset_account_name>, apiURL:<participant_api_url>
    And allower anchor: <allow_trust_participant_id> send alow trust request to allow trust DA asset_code: <participant_trust_asset_code> trustSender: <participant_id> trust_account <participant_trust_asset_account_name>, allowURL:<alow_trust_url>
    And id: <participant_id> asset_code: <participant_trust_asset_code> asset_type: <participant_trust_asset_type> issuer: <participant_trust_asset_issuer> in trust_account <participant_trust_asset_account_name> trusted list, apiURL:<participant_api_url>

    Examples:
      | global_whitelist_url           | participant_id             | allow_trust_participant_id | participant_api_url             | participant_trust_asset_code  | participant_trust_asset_type | participant_trust_asset_issuer | participant_trust_asset_limit             | participant_trust_asset_account_name           | alow_trust_url               |
      | "ENV_KEY_PARTICIPANT_1_WL_URL" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE"   | "ENV_KEY_ANCHOR_ASSET_TYPE"  | "ENV_KEY_ANCHOR_ID"            | "ENV_KEY_PARTICIPANT_1_TRUST_ASSET_LIMIT" | "ENV_KEY_PARTICIPANT_1_ISSUING_ACCOUNT_NAME"   | "ENV_KEY_ANCHOR_SERVICE_URL" |
      | "ENV_KEY_PARTICIPANT_1_WL_URL" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE"   | "ENV_KEY_ANCHOR_ASSET_TYPE"  | "ENV_KEY_ANCHOR_ID"            | "ENV_KEY_PARTICIPANT_1_TRUST_ASSET_LIMIT" | "ENV_KEY_PARTICIPANT_1_OPERATING_ACCOUNT_NAME" | "ENV_KEY_ANCHOR_SERVICE_URL" |
      | "ENV_KEY_PARTICIPANT_1_WL_URL" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE_2" | "ENV_KEY_ANCHOR_ASSET_TYPE"  | "ENV_KEY_ANCHOR_ID"            | "ENV_KEY_PARTICIPANT_1_TRUST_ASSET_LIMIT" | "ENV_KEY_PARTICIPANT_1_ISSUING_ACCOUNT_NAME"   | "ENV_KEY_ANCHOR_SERVICE_URL" |
      | "ENV_KEY_PARTICIPANT_2_WL_URL" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE_2" | "ENV_KEY_ANCHOR_ASSET_TYPE"  | "ENV_KEY_ANCHOR_ID"            | "ENV_KEY_PARTICIPANT_1_TRUST_ASSET_LIMIT" | "ENV_KEY_PARTICIPANT_1_OPERATING_ACCOUNT_NAME" | "ENV_KEY_ANCHOR_SERVICE_URL" |
      | "ENV_KEY_PARTICIPANT_2_WL_URL" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE"   | "ENV_KEY_ANCHOR_ASSET_TYPE"  | "ENV_KEY_ANCHOR_ID"            | "ENV_KEY_PARTICIPANT_2_TRUST_ASSET_LIMIT" | "ENV_KEY_PARTICIPANT_2_ISSUING_ACCOUNT_NAME"   | "ENV_KEY_ANCHOR_SERVICE_URL" |
      | "ENV_KEY_PARTICIPANT_2_WL_URL" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE"   | "ENV_KEY_ANCHOR_ASSET_TYPE"  | "ENV_KEY_ANCHOR_ID"            | "ENV_KEY_PARTICIPANT_2_TRUST_ASSET_LIMIT" | "ENV_KEY_PARTICIPANT_2_OPERATING_ACCOUNT_NAME" | "ENV_KEY_ANCHOR_SERVICE_URL" |
      | "ENV_KEY_PARTICIPANT_2_WL_URL" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE_2" | "ENV_KEY_ANCHOR_ASSET_TYPE"  | "ENV_KEY_ANCHOR_ID"            | "ENV_KEY_PARTICIPANT_2_TRUST_ASSET_LIMIT" | "ENV_KEY_PARTICIPANT_2_ISSUING_ACCOUNT_NAME"   | "ENV_KEY_ANCHOR_SERVICE_URL" |
      | "ENV_KEY_PARTICIPANT_2_WL_URL" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE_2" | "ENV_KEY_ANCHOR_ASSET_TYPE"  | "ENV_KEY_ANCHOR_ID"            | "ENV_KEY_PARTICIPANT_2_TRUST_ASSET_LIMIT" | "ENV_KEY_PARTICIPANT_2_OPERATING_ACCOUNT_NAME" | "ENV_KEY_ANCHOR_SERVICE_URL" |
