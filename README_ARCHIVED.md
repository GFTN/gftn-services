# world-wire-services

[![Build Status](https://travis.ibm.com/gftn/world-wire-services.svg?token=7Nw63z3Dh19HGjHzQq87&branch=master)](https://travis.ibm.com/gftn/world-wire-services) **Master Banch**
[![Build Status](https://travis.ibm.com/gftn/world-wire-services.svg?token=7Nw63z3Dh19HGjHzQq87&branch=development)](https://travis.ibm.com/gftn/world-wire-services) **Development Branch**

Parent source repository for all micro services under world wire


To install dependencies and check all services build correctly after code change run:
at root path

````
make dep
````
it creates /vendor folder at root


To compile all services, at root:

````
cd integration-tests
make build-go

````

To build all docker images, at root:

````
cd integration-tests
make service-docker build=all

````

To build a particular docker image, at root:

````
cd integration-tests
make service-docker build=anchor-service

````

# Authentication Service
Please see ./auth-service/README.md for instructions including screenshots and code snippets.

# World Wire Web (Portal)
Please see ./world-wire-web/README.md for instructions to run the portal application locally.
