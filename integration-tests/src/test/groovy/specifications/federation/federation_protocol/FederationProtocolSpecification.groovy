package specifications.federation.federation_protocol

import gftn.util.TestUtil
import groovyx.net.http.RESTClient
import spock.lang.Shared
import spock.lang.Specification

class FederationProtocolSpecification extends Specification {

    @Shared
    def api = new RESTClient(TestUtil.getFederationServiceInternalApiUrl())


    def setupSpec() {
        api.handler.failure = api.handler.success
        reportHeader("Stellar Federation Protocol Tests")
    }


    def "Service should comply with the Stellar Federation Protocol"() {




    }


}
