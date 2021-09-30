#!/bin/bash

beforeInstall()
{
    npm install protractor && webdriver-manager update
    echo "Starting headless browser for testing..."
    google-chrome --no-sandbox --disable-dev-shm-usage --headless --disable-gpu --remote-debugging-port=9222 http://localhost &
}

installDependencies()
{
    echo "Installing node module dependencies..."
    npm install
    npm install -g @angular/cli
    cd functions && npm install && cd ..
}

lintAndCompile()
{
    echo "Starting typescript compilation..."
    tsc -p tsconfig.json
    echo "Starting linting..."
    ng lint
}

testHeadless()
{
    echo "Starting headless browser testing..."
    npm run test -- --no-watch --no-progress --browsers ChromeHeadlessCI
}

testEndToEnd()
{
    echo "Starting end-to-end testing..."
    ng e2e --protractor-config=e2e/protractor-ci.conf.js && ng build --configuration dev
}

afterFailure()
{

  # dump the last 2000 lines of our build to show error
  tail -n 2000 build.log

}

CMD=$1

if [ $CMD = "lintAndCompile" ]; then
    lintAndCompile
elif [ $CMD = "beforeInstall" ]; then
    beforeInstall
elif [ $CMD = "installDependencies" ]; then
    installDependencies
elif [ $CMD = "testEndToEnd" ]; then
    testEndToEnd
elif [ $CMD = "testHeadless" ]; then
    testHeadless
elif [ $CMD = "afterFailure" ]; then
    afterFailure
fi