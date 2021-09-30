package gftn.util;

public class TestUtil {

    public static String getFederationServiceInternalApiUrl() {
        String baseUrl = System.getenv("FEDERATION_SERVICE_INTERNAL_API_BASE_URL");
        return baseUrl;
    }

    public static String getMockODFIServiceApiBaseUrl() {
        String baseUrl = System.getenv("MOCK_ODFI_SERVICE_API_BASE_URL");
        return baseUrl;
    }

    public static String getConfigurationResponderBaseUrl() {
        String baseUrl = System.getenv("CONFIGURATION_RESPONDER_BASE_URL");
        return baseUrl;
    }


}
