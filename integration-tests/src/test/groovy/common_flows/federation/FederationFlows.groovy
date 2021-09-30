package common_flows.federation

import groovyx.net.http.RESTClient

class FederationFlows {

    private RESTClient client

    FederationFlows(RESTClient rc) {
        this.client = rc
    }

    String FindParticipant(String participantDomain, String accountIdentifier) {

        def response = this.client.get(
                path: "/v1/find/" + participantDomain + "/" + accountIdentifier,
                requestContentType: "application/json"
        )

        if (response.status != 200) {
            return ""
        }

        println(response.data)
        def stellarAddress = response.data["stellar_network_address"]
        return stellarAddress

    }


}
