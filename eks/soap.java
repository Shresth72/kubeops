import static spark.Spark.*;

public class Main {
    public static void main(String[] args) {
        get("/", (req, res) -> {
            res.type("text/xml");
            return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
                    "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">\n" +
                    "    <soap:Body>\n" +
                    "        <Response>\n" +
                    "            <Message>Hello, World!</Message>\n" +
                    "        </Response>\n" +
                    "    </soap:Body>\n" +
                    "</soap:Envelope>";
        });

        get("/error", (req, res) -> {
            res.type("text/xml");
            return "<?xml version=\"1.0\" encoding=\"UTF-8\"?>\n" +
                    "<soap:Envelope xmlns:soap=\"http://schemas.xmlsoap.org/soap/envelope/\">\n" +
                    "    <soap:Body>\n" +
                    "        <Response>\n" +
                    "            <Message>This is error 404</Message>\n" +
                    "        </Response>\n" +
                    "    </soap:Body>\n" +
                    "</soap:Envelope>";
        })
    }   
}
