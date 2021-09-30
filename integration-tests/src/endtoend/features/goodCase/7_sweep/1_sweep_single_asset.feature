
# Feature: World Wire
#   In order to helping a customer send money
#   As a developer
#   I want to issue anchor asset in WW

#   Scenario Outline: sweep accounts
#     Given id: <participant_id> check id: <participant_id> operating account:<targer_account> exist, participantAPIURL: <participant_api_url>
#     And id: <participant_id> asset_code: <asset_code> asset_type: <asset_type> issuer: <asset_issuer> in trust_account <targer_account> trusted list, apiURL:<participant_api_url>
#     And id: <participant_id> check id: <participant_id> operating account:<source_account> exist, participantAPIURL: <participant_api_url>
#     And id: <participant_id> asset_code: <asset_code> asset_type: <asset_type> issuer: <asset_issuer> in trust_account <targer_account> trusted list, apiURL:<participant_api_url>
#     And id: <participant_id> check asset_code: <asset_code> issuer: <asset_issuer> account_name: <source_account> balance is greater than sweep amount: <sweep_amount>, participant_api_url:<participant_api_url>
#     And id: <participant_id> check asset_code: <asset_code> issuer: <asset_issuer> account_name: <source_account> balance before transaction, participant_api_url:<participant_api_url>
#     And id: <participant_id> check asset_code: <asset_code> issuer: <asset_issuer> account_name: <targer_account> balance before transaction, participant_api_url:<participant_api_url>
#     When id: <participant_id> sweep asset_code: <asset_code> issuer: <asset_issuer> asset_type: <asset_type> amount: <sweep_amount> account_name: <source_account> to account_name: <targer_account>, participant_api_url:<participant_api_url>
#     Then id: <participant_id> check asset_code: <asset_code> issuer: <asset_issuer> account_name: <targer_account> balance increase: <sweep_amount>, participant_api_url:<participant_api_url>
#     And id: <participant_id> check asset_code: <asset_code> issuer: <asset_issuer> account_name: <source_account> balance decrease: <sweep_amount>, participant_api_url:<participant_api_url>


#     Examples:
#       | participant_api_url             | asset_code                  | asset_type                  | participant_id             | targer_account                                 | source_account                                   | asset_issuer        | sweep_amount             |
#       | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_ASSET_CODE" | "ENV_KEY_ANCHOR_ASSET_TYPE" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_PARTICIPANT_1_OPERATING_ACCOUNT_NAME" | "ENV_KEY_PARTICIPANT_1_OPERATING_ACCOUNT_2_NAME" | "ENV_KEY_ANCHOR_ID" | "ENV_KEY_SWEEP_AMOUNT_1" |
