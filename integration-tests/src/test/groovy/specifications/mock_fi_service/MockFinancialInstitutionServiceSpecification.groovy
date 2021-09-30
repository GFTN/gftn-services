package specifications.mock_fi_service

import gftn.util.TestUtil
import groovyx.net.http.ContentType
import groovyx.net.http.RESTClient
import spock.lang.Shared
import spock.lang.Specification

class MockFinancialInstitutionServiceSpecification extends Specification {


    @Shared
    def api = new RESTClient(TestUtil.getMockODFIServiceApiBaseUrl())


    def setupSpec() {
        api.handler.failure = api.handler.success
        reportHeader("Mock FI Service Tests")
    }


    def "Mock service responds positively to Federation Service requests for an account number"() {

        given: "An account number (0000000101) that the mock service will respond positively"
        def accountNumber = "0000000101"
        def requestBody = [
                "account_identifier":accountNumber
        ]

        when: "We call the mock FI with this account number"
        def response = api.post(
                path: "/verifications/account",
                requestContentType: ContentType.JSON,
                body: requestBody
        )

        then: "The service responds with a HTTP OK"
        response != null
        response.status == 200

    }


    def "Mock service responds negatively to Federation Service requests for an account number"() {

        given: "An account number (0000000901) that the mock service will respond negatively"
        def accountNumber = "0000000901"
        def requestBody = [
                "account_identifier":accountNumber
        ]

        when: "We call the mock FI with this account number"
        def response = api.post(
                path: "/verifications/account",
                requestContentType: ContentType.JSON,
                body: requestBody
        )

        then: "The service responds with a HTTP Not Found"
        response != null
        response.status == 404

    }





}
