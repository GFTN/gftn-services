# Global Financial Transaction Network Services

This code was developed at IBM during 2017-2020, and contributed to open source in September 2021.

## Overview
This is a monorepo of code 
intended to support a scalable cloud service for issuing and trading digital assets.
The website is written in typescript and javascript, and can be found in the directory `world-wire-web`.  

There are two types of users: administrator and participant.  
The administrator is primarily implemented by the `administration-service`, although some methods are also
implemented in the world-wire-web javascript.  The main role of the administrator
is to manage participants.

When participants are [deployed](automation-service/automate/participant/deploy.go) by the administrator, several services
are spawned for the participant and used only by that participant.  These are:
* `api-service`
* `crypto-service`
* `send-service`
* `payment-listener`
* `ww-gateway`

Finally, in addition to the `administration-service`, there are the global services:
* `quotes-service`
* `gas-service`
* `participant-registry`
* `fee-service`
* `anchor-service`
* `auth-service`
* `global-whitelist-service`

## To Developers
Probably the best way to understand what is going on is to look through the postman collections.
For example, [this directory](integration-tests/src/worldwire-collections) 
covers a lot of ground.

The deployment scripts and yaml files are in the `automation-service` directory tree, especially `k8s`.



[Archived Readme](README_ARCHIVED.md).

