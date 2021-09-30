package common_flows.service_check

import groovyx.net.http.RESTClient

class ServiceCheckFlows {


    private RESTClient client

    ServiceCheckFlows(RESTClient rc) {
        this.client = rc
    }

    boolean ServiceCheck() {

        def response = this.client.get(
                path: "/v1/service_check",
                requestContentType: "application/json"
        )

        if (response.status == 200) {
            return true
        }

        println("Response status:  " + response.status)
        return false

    }

}
