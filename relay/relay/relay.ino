// Function
//      Using 2-way relay module
// Connection
//      Relay module->Arduino: Gnd->Gnd; In1->D5; In2->D6; Vcc->5V.
//      Relay module->Power switcher: NO1->AC PowerA1; NO2->AC PowerA2.
//      Relay module->Power switcher: COM1->COM2->AC PowerB.
void setup(){
    pinMode(2, OUTPUT);
    pinMode(3, OUTPUT);
    digitalWrite(2, LOW); // HIGH -> open in common-open of relay
    digitalWrite(3, LOW); // HIGH -> open in common-open of relay
}

void loop(){
    // delay(5000);
    // digitalWrite(2, LOW);
    // delay(5000);
    // digitalWrite(2, HIGH);
}
