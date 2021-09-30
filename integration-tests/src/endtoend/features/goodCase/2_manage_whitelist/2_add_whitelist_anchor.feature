
Feature: World Wire
  In order to helping a customer send money
  As a developer
  I want to add another participant to my whitelist in worldwire

  Scenario Outline: Add participant to whitelist
    # Given id: <participant_id> check participant: <whitelist_target_participant_id>, was in Workd Wire, apiURL: <participant_api_url>
    When id: <participant_id> add participant: <whitelist_target_participant_id> - whitelistURL: <global_whitelist_url>
    Then id: <participant_id> check participant: <whitelist_target_participant_id> was in the whitelist - whitelistURL: <global_whitelist_url>

    Examples:
      | global_whitelist_url           | participant_api_url             | participant_id             | whitelist_target_participant_id |
      | "ENV_KEY_PARTICIPANT_1_WL_URL" | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_ANCHOR_ID"             |
      | "ENV_KEY_PARTICIPANT_2_WL_URL" | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_ANCHOR_ID"             |
      | "ENV_KEY_PARTICIPANT_1_WL_URL" | "ENV_KEY_PARTICIPANT_1_API_URL" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_1_ID"      |
      | "ENV_KEY_PARTICIPANT_2_WL_URL" | "ENV_KEY_PARTICIPANT_2_API_URL" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_PARTICIPANT_2_ID"      |
