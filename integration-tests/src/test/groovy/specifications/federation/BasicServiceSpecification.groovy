package specifications.federation

import common_flows.service_check.ServiceCheckFlows
import gftn.util.TestUtil
import groovyx.net.http.RESTClient
import spock.lang.Shared
import spock.lang.Specification

class BasicServiceSpecification extends Specification {

    @Shared
    def api = new RESTClient(TestUtil.getFederationServiceInternalApiUrl())


    def setupSpec() {
        api.handler.failure = api.handler.success
        reportHeader("Basic Services Tests")
    }


    def "The Federation Service will have an internal API Endpoint"() {

        given: "A federation service with an internal service check API endpoint"
        def serviceCheckFlow = new ServiceCheckFlows(api)


        when: "The internal API endpoint is checked for system status"
        def didPass = serviceCheckFlow.ServiceCheck()


        then: "The service will return a positive result"
        didPass == true

    }

}
