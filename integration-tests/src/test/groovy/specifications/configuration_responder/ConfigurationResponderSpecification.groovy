package specifications.configuration_responder

import com.moandjiezana.toml.Toml
import gftn.util.TestUtil
import groovyx.net.http.ContentType
import groovyx.net.http.RESTClient
import spock.lang.Shared
import spock.lang.Specification

class ConfigurationResponderSpecification extends Specification {


    @Shared def api = new RESTClient(TestUtil.getConfigurationResponderBaseUrl())

    def setupSpec() {
        api.handler.failure = api.handler.success
        reportHeader("Configuration Responder Tests")
    }




    def "Able to respond with a standard stellar TOML configuration"() {

        when: "A request for a stellar.toml is completed"
        def response = api.get(
                path: "/.well_known/stellar.toml",
                requestContentType: "application/toml"
        )

        def tomlResponse = null;
        if (response.responseData instanceof InputStream) {
            tomlResponse = new Toml().read((InputStream)response.responseData)
        }



        then: "The participant's configuration is returned"
        response != null
        response.status == 200
        tomlResponse != null
        tomlResponse.getString("AUTH_SERVER") == "http://compliance.specification.integrationtest.gftn.io"
        tomlResponse.getString("FEDERATION_SERVER") == "http://federation.specification.integrationtest.gftn.io"


    }





}
