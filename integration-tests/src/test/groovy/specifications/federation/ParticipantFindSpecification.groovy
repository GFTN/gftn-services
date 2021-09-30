package specifications.federation

import common_flows.federation.FederationFlows
import gftn.util.TestUtil
import groovyx.net.http.RESTClient
import spock.lang.Shared
import spock.lang.Specification

class ParticipantFindSpecification extends Specification {


    @Shared
    def api = new RESTClient(TestUtil.getFederationServiceInternalApiUrl())


    def setupSpec() {
        api.handler.failure = api.handler.success
        reportHeader("Participant Find Tests")
    }


    def "Federation service should be able to find a participant on the network when requested by the API Service"() {

        given: "A GFTN participant domain and an account identifier at the participant"
        def participantDomain = "rdfi.payments.gftn.io"
        def accountIdentifier = "0000000101"
        def federationFlows = new FederationFlows(this.api)

        when: "We request the federation server to find a participant on the network"
        def cryptoAddress = federationFlows.FindParticipant(participantDomain, accountIdentifier)


        then: "The participant is found, and a Stellar address is returned for transacting"
        cryptoAddress != null
        cryptoAddress == "GAXMOCKSRTELLARADDRESS"


    }


}
