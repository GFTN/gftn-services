#!/bin/bash

afterSuccess()
{
  cd ~/build/gftn/gftn-services
  version=`cat VERSION`
  export version=$version
  cd integration-tests/
  make push-dockers

}

afterFailure()
{

  # dump the last 2000 lines of our build to show error
  tail -n 2000 build.log

}

CMD=$1

if [ $CMD = "afterSuccess" ]; then
    afterSuccess
elif [ $CMD = "afterFailure" ]; then
    afterFailure
fi