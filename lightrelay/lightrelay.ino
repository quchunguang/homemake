// Function
//      Control LEDs and 2-way relay with Light sensor
// Reference
//      http://www.geek-workshop.com/thread-1340-1-1.html
// Connection
//      Light sensor: Gnd->10KOhm->LS->5V;
//                               |->A0.
//      LED: Gnd->330Ohm->pinLedL->D2;
//           Gnd->330Ohm->pinLedM->D2;
//           Gnd->330Ohm->pinLedH->D2.
//      Relay module->Arduino: Gnd->Gnd; In1->D5; In2->D6; Vcc->5V.
//      Relay module->Power switcher: NO1->AC PowerA1; NO2->AC PowerA2.
//      Relay module->Power switcher: COM1->COM2->AC PowerB.
// Tested
//      arduino uno/nano328

int sensorValue;
int baudrate = 9600;

int pinLedL = 2;
int pinLedM = 3;
int pinLedH = 4;
int pinRelay1 = 5;  // HIGH ON
int pinRelay2 = 6;  // HIGH OFF
int pinReadLight = A0;

bool R1OnR2Off = false;

void setup()
{
    Serial.begin(baudrate);

    pinMode(pinLedL, OUTPUT);
    pinMode(pinLedM, OUTPUT);
    pinMode(pinLedH, OUTPUT);
    pinMode(pinRelay1, OUTPUT);
    pinMode(pinRelay2, OUTPUT);

    digitalWrite(pinLedL, LOW);
    digitalWrite(pinLedM, LOW);
    digitalWrite(pinLedH, LOW);
    digitalWrite(pinRelay1, LOW);
    digitalWrite(pinRelay2, LOW);
}

void loop()
{
    sensorValue = analogRead(pinReadLight);

    Serial.println(sensorValue);

    if (sensorValue > 950) {
        digitalWrite(pinLedL, HIGH);
        digitalWrite(pinLedM, HIGH);
        digitalWrite(pinLedH, HIGH);

        if (!R1OnR2Off) {
            digitalWrite(pinRelay1, HIGH);
            digitalWrite(pinRelay2, LOW);
            R1OnR2Off = true;
            Serial.println("[X] Relay 1 On, Relay 2 Off");
        }
    } else if (sensorValue > 800) {
        digitalWrite(pinLedL, HIGH);
        digitalWrite(pinLedM, HIGH);
        digitalWrite(pinLedH, LOW);
    } else if (sensorValue > 700) {
        digitalWrite(pinLedL, HIGH);
        digitalWrite(pinLedM, LOW);
        digitalWrite(pinLedH, LOW);
    } else {
        digitalWrite(pinLedL, LOW);
        digitalWrite(pinLedM, LOW);
        digitalWrite(pinLedH, LOW);

        if (R1OnR2Off) {
            digitalWrite(pinRelay1, LOW);
            digitalWrite(pinRelay2, HIGH);
            R1OnR2Off = false;
            Serial.println("[X] Relay 1 Off, Relay 2 On");
        }
    }

    delay(1000);
}
