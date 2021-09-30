# anchor-service
This service will be used by stable coin issuers like Stronghold. This API will evolve as new stable coin issuers are boarded


1. Anchor account is created outside WW, (could be created using stellar lab for testing)

1. Add IBM account as a signer to this anchor account (can be done using stellar lab)

1. Anchor issues digital asset (using stellar Lab) outside of WW

1. Onboard an anchor on WW, WW admin creates an entry in PR with role as an anchor with country and domain

1. Register Anchor:
Endpoint: admin/anchor/usd.stronghold.co/register
Validates if this account is valid and has setoptions set correctly with IBM admin account, it then registers with PR

1. WW admin will then have to update PR 
with its role as IS, country, and API service URL 
for stronghold it should be:  "api_service_url": "https://sandbox.stronghold.co/v1/venues",

1. Register asset:
 Endpoint: /admin/anchor/{anchor_domain}/onboard/assets
 Creates a trust with IBM admin account
 
 
 1. Participant will send create trust operation
 
 
 1. Anchor will do allow trust operation
 Endpoint: anchor/trust
 
 
 1. anchor will Discover WW participant account address:
 Endpoint: anchor/address/stellar
 
 
 1. anchor can now send DA to participant using funding endpoint
 
 1. return or withdraw anchor endpoint, 
 
 API service Endpoint: transactions/settle/da
 
 
 
 
 
 
 
 
