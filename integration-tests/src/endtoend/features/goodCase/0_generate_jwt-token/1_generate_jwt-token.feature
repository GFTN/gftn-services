
Feature: World Wire goodCase/0_generate_JWT-token/1_generate_jwt-token
    In order to calling world wire API
    As a developer
    I want to generate jwt token for request

    Scenario Outline: Generate JWT-TOKEN
        Given user: <user1_email> generate a FID
        And user: <user2_email> generate a FID
        And user: <user1_email> send request to request a jwt-token for participant: <participant_id> IID: <participant_iid> using totpkey: <user1_totp_key>, auth_url:<auth_url>
        When user: <user2_email> send request to approve a jwt-token for participant: <participant_id> IID: <participant_iid> using totpkey: <user2_totp_key>, auth_url:<auth_url>
        Then user: <user1_email> send request to get a jwt-token for participant: <participant_id> IID: <participant_iid> using totpkey: <user1_totp_key> naming as: <token_variables> , auth_url:<auth_url>
        Examples:
            | user1_email            | user2_email            | user1_totp_key            | user2_totp_key            | participant_id             | participant_iid             | token_variables        | auth_url           |
            | "ENV_KEY_USER_1_EMAIL" | "ENV_KEY_USER_2_EMAIL" | "ENV_KEY_USER_1_TOTP_KEY" | "ENV_KEY_USER_2_TOTP_KEY" | "ENV_KEY_PARTICIPANT_1_ID" | "ENV_KEY_PARTICIPANT_1_IID" | "PARTICIPANT_1_JWT_TOKEN" | "ENV_KEY_AUTH_URL" |
            | "ENV_KEY_USER_1_EMAIL" | "ENV_KEY_USER_2_EMAIL" | "ENV_KEY_USER_1_TOTP_KEY" | "ENV_KEY_USER_2_TOTP_KEY" | "ENV_KEY_PARTICIPANT_2_ID" | "ENV_KEY_PARTICIPANT_2_IID" | "PARTICIPANT_2_JWT_TOKEN" | "ENV_KEY_AUTH_URL" |
            | "ENV_KEY_USER_1_EMAIL" | "ENV_KEY_USER_2_EMAIL" | "ENV_KEY_USER_1_TOTP_KEY" | "ENV_KEY_USER_2_TOTP_KEY" | "ENV_KEY_ANCHOR_ID"        | "ENV_KEY_ANCHOR_IID"        | "ANCHOR_JWT_TOKEN"       | "ENV_KEY_AUTH_URL" |