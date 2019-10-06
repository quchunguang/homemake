/*
 * ESP8266 (WeMosD1) WiFi Relay Control
 * 
 * learnelectronics
 * 05 JUN 2017
 * 
 * www.youtube.com/c/learnelectronics
 * arduino0169@gmail.com
*/
#include <ESP8266WiFi.h>
 
const char* ssid = "";
const char* password = "";
 
int ledPin = LED_BUILTIN;
WiFiServer server(80);
 
void setup() {
  Serial.begin(9600);
  delay(10);
 
 
  pinMode(ledPin, OUTPUT);
  digitalWrite(ledPin, HIGH);
 
  // Connect to WiFi network
  Serial.println();
  Serial.println();
  Serial.print("Connecting to ");
  Serial.println(ssid);
 
  WiFi.begin(ssid, password);
 
  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }
  Serial.println("");
  Serial.println("WiFi connected");
 
  // Start the server
  server.begin();
  Serial.println("Server started");
 
  // Print the IP address
  Serial.print("Use this URL : ");
  Serial.print("http://");
  Serial.print(WiFi.localIP());
  Serial.println("/");
 
}
 
void loop() {
  // Check if a client has connected
  WiFiClient client = server.available();
  if (!client) {
    return;
  }
 
  // Wait until the client sends some data
  Serial.println("new client");
  while(!client.available()){
    delay(1);
  }
 
  // Read the first line of the request
  String request = client.readStringUntil('\r');
  Serial.println(request);
  client.flush();
 
  // Match the request
 
  int value = HIGH;
  if (request.indexOf("/LED=ON") != -1) {
    digitalWrite(ledPin, LOW);
    value = LOW;
  } 
  if (request.indexOf("/LED=OFF") != -1){
    digitalWrite(ledPin, HIGH);
    value = HIGH;
  }
 
 
 
  // Return the response
  client.println("HTTP/1.1 200 OK");
  client.println("Content-Type: text/html");
  client.println(""); //  do not forget this one
  client.println("<!DOCTYPE HTML>");
  client.println("<html>");
 
  client.print("<h1>Led pin is now: ");
 
  if(value == LOW) {
    client.print("On");  
  } else {
    client.print("Off");
  }
  client.println("</h1><br><br>");
  client.println("<h1>Click <a href=\"/LED=ON\">here</a> Turn relay ON<br></h1>");
  client.println("<h1>Click <a href=\"/LED=OFF\">here</a> Turn relay OFF<br></h1>");
  client.println("</html>");
 
  delay(1);
  Serial.println("Client disconnected");
  Serial.println("");
 
}
