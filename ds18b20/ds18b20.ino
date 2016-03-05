// Function
//      DS18B20 temperature sensor
// Library
//      https://github.com/PaulStoffregen/OneWire
//      https://github.com/milesburton/Arduino-Temperature-Control-Library
// Reference
//      http://www.instructables.com/id/DS18B20-temperature-probe-with-LCD/?ALLSTEPS
// Connection
//      Temperature Sensor - DS18B20: Vcc->5V(Arduino);
//                                    Gnd->Gnd(Arduino);
//                                    Data->4.7KOhm->5V(Arduino);
//                                        |->D7.
// Tested
//      arduino uno/nano328
#include <OneWire.h>
#include <DallasTemperature.h>

#define pinDS18B20 7 // one wire bus
#define BAUDRATE 9600

OneWire oneWire(pinDS18B20);
DallasTemperature sensors(&oneWire);

void setup()
{
    Serial.begin(BAUDRATE);

    // Start up the library
    sensors.begin();
}

void loop()
{
    // Send the command to get temperatures
    sensors.requestTemperatures();
    Serial.print("Temperature for Device 1 is: ");
    // Why "byIndex"?
    // You can have more than one IC on the same bus.
    // 0 refers to the first IC on the wire
    Serial.println(sensors.getTempCByIndex(0));

    delay(1000);
}
