# integration-tests


## Your very own GFTN testnet

* Clone this repository and cd to its root:
  * `git clone https://github.com/GFTN/integration-tests`
  * `cd integration-tests`

* Run docker-compose against the integration test suite:
  * `docker-compose -f src/test/resources/deployment/docker/docker-compose-ci-integration-tests.yaml up`

* The tests will run, hopefully pass.
* You will have an ODFI and RDFI GFTN instances running.  

* docker-compose.yaml is added to set up odfi, rdfi and participant registry running on one machine locally using docker images


## Spock Integration Tests for GFTN


