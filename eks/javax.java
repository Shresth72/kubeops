import javax.jws.WebMethod;
import javax.jws.WebService;
import javax.xml.ws.Endpoint;

@WebService
public class SoapWebService {

    @WebMethod
    public String getXmlResponse() {
        return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
                    "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">\n" +
                    "    <soap:Body>\n" +
                    "        <Response>\n" +
                    "            <Message>Hello, World!</Message>\n" +
                    "        </Response>\n" +
                    "    </soap:Body>\n" +
                    "</soap:Envelope>";
    }

    public static void main(String[] args) {
        // Specify the address at which the SOAP web service will be available
        String address = "http://localhost:8080/soapWebService";
        
        // Create an instance of the SOAP web service
        SoapWebService soapWebService = new SoapWebService();
        
        // Publish the SOAP web service at the specified address
        Endpoint.publish(address, soapWebService);
        
        // Output to console to indicate that the server has started
        System.out.println("SOAP web service started at: " + address);
    }
}
