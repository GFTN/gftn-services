#!/bin/bash

beforeInstall()
{
    google-chrome --no-sandbox --disable-dev-shm-usage --headless --disable-gpu --remote-debugging-port=9222 http://localhost &
}

installDependencies()
{
    npm install
    cd functions && npm install && cd ..
}

runScript()
{
    tsc -p tsconfig.json
    ng build --configuration staging
}

afterFailure()
{

  # dump the last 2000 lines of our build to show error
  tail -n 2000 build.log

}

afterSuccess()
{
    echo "Build and tests succeeded for Dev environment"
    echo "Beginning deployment..."
    npm install -g firebase-tools
    firebase deploy -m "Travis deploy" --non-interactive --project st21-251107 --token $FIREBASE_TOKEN_STAGING --only functions,hosting
}

CMD=$1

if [ $CMD = "beforeInstall" ]; then
    beforeInstall
elif [ $CMD = "runScript" ]; then
    runScript
elif [ $CMD = "installDependencies" ]; then
    installDependencies
elif [ $CMD = "afterFailure" ]; then
    afterFailure
elif [ $CMD = "afterSuccess" ]; then
    afterSuccess
fi