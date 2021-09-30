# api-service
Client-facing API Service, orchestrator of World Wire operations

[![Build Status](https://travis.ibm.com/gftn/api-service.svg?token=QUdYcSEpKprUxsFBpz8s&branch=master)](https://travis.ibm.com/gftn/api-service)

Required Software Dependencies:
 - Go/Golang: https://golang.org/doc/install
 - Glide, a package manager for Go: https://glide.sh/
 - A package manager. We use:
  - NPM: https://docs.npmjs.com/downloading-and-installing-node-js-and-npm
  - Homebrew (Mac OSX): https://brew.sh/


# Running API Service locally
In order to run and test against the API-Service locally, you will need several steps to get set up.

## Step 1 -- install all dependencies
The first step is to run dep at the root of the project:
```
make dep
```

## Step 2 -- Configure local environment
Change to newly cloned utility repo and edit the `nodeconfig.toml` file to contain the following:

(make sure to replace filler values with actual values)
```
[DISTRIBUTION_ACCOUNTS]

 [DISTRIBUTION_ACCOUNTS.default]
   NODE_ADDRESS = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
   NODE_SEED = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

[IBM_TOKEN_ACCOUNT]
 NODE_ADDRESS = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
 NODE_SEED = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"

[ISSUING_ACCOUNT]
 NODE_ADDRESS = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
 NODE_SEED = "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX"
```

Set environment variables by copying the variables found in `.env.example` to your OS's config file path: `~/.profile`.

Replace <GOPATH> with your set GOPATH for your project.

Alternatively, for VSCode user, copy over:
- `launch.json.example` to `launch.json`
- `settings.json.example` to `settings.json`

DON'T FORGET: Replace <GOPATH> with your set GOPATH for your project.




> To open the `~/.profile` with vscode use `code ~/.profile` and then after updating the file use `source ~/.profile` to enable and set the updates



### Environment Variables/Config

#### Required


* SERVICE_PORT
  * The porn number this service should bind to at startup

* VERIFY_ACCOUNT_IDENTIFIER_URL
  * The URL the API Service should call in order to verify an account number

* NOTIFY_NEW_ADDRESS_URL
  * The URL the API Service should call in order to notify of a new crypto address generated for the client

* FEDERATION_SERVICE_INTERNAL_API_URL
  * The base URL for the internal Federation Service Endpoint for this Participant

* QUOTES_SERVICE_INTERNAL_API_URL
  * The base URL for the internal QUOTES Service Endpoint for this Participant

* ISSUING_TOKEN
  * token issued by issuing account by this participant

* HORIZON_CLIENT_URL
  * Horizon client url

* FRIENDBOT_URL
  * Friendbot url to fund accounts

* SERVICE_DIST_ACCOUNT_KEYS_FILE
  * Distribution account keys File path

* SERVICE_ISSUE_ACCOUNT_KEYS_FILE
  * issuing account keys file path

* PARTICIPANT_REGISTRY_URL
  * participant registry url

* OPTIONAL: Generate JWT RS256 key (Must set ENABLE_JWT = true). To do so, run:
    * ssh-keygen -t rsa -b 4096 -f jwtRS256.key

      Don't add passphrase
    * openssl rsa -in jwtRS256.key -pubout -outform PEM -out jwtRS256.key.pub
    * cat jwtRS256.key
    * cat jwtRS256.key.pub

## Step 3 -- Install dependencies

This step assumes you have npm installed from nodejs.org and already have Golang etc. installed on your system

You will need to change directories to the api-service directory:
```
cd ~/<GOPATH>/src/github.com/GFTN/api-service
```

Mac only: Install homebrew if not installed:
```
npm install -g homebrew
```

After finishing the homebrew installation or if already installed previously, use the following command to install Glide:
```
brew install glide # Mac OSX

curl https://glide.sh/get | sh # all other OSes
```

With Glide installed, you can now run the following from inside the `api-service` directory:
```
glide install
```

This installs all the golang dependencies found in the project files in api-service.

After these packages are installed, you can build the project with the Makefile:

```
make
```

Alternative you can run each of the individual commands found in the Makefile, if the Makefile is broken or you need to do an intermediary step in between.

To build the project manually, run the following command to compile the project into an executable:
```
go build # general command for all OSes
go build -o api-service.exe # Windows
go build -o api-service.osx # Mac OSX

```

To run the project:
```
./api-service # or ./api-service.exe or ./api-service.osx, etc.
```

Alternatively for VSCode users: go into the Debug panel, click the play button on the launch service for api-service to run the project.

The API Service should now be running locally and accepting requests on `localhost:8080`, or whichever SERVICE_PORT you have set it to `localhost:<SERVICE_PORT>`.
