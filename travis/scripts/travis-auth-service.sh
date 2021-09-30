#!/bin/bash

authInstall()
{

   cd auth-service/functions
   npm uninstall typescript --no-save # ensure typscript is not installed
   npm uninstall tslint --no-save # ensure tslint is not installed
   npm install -g typescript # install typescript globally
   npm install -g tsoa # install TSOA globally
   npm install # install dependencies at  ./auth-service/functions/package.json

}

authRun()
{

   cd auth-service/functions
   tsoa swagger
   tsoa routes
   cd .. # in ./auth-service
   tsc -p tsconfig.json
   # cd ${TRAVIS_BUILD_DIR} # back to root

}

afterFailure()
{

  # dump the last 2000 lines of our build to show error
  tail -n 2000 build.log

}

afterSuccess()
{

  # Log that the build was a success
  echo "build succeeded.."

}

CMD=$1

if [ $CMD = "authInstall" ]; then
 authInstall
elif [ $CMD = "authRun" ]; then
 authRun
elif [ $CMD = "afterFailure" ]; then
 afterFailure
elif [ $CMD = "afterSuccess" ]; then
 afterSuccess
fi