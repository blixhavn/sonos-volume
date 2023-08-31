#include <DigiUSB.h>

const int pinA = 2; // CLK pin
const int pinB = 0; // DT pin
const int buttonPin = 1; // SW pin

int aState;
int aLastState;
int buttonState = 0;
int lastButtonState = 0;

void setup() {
  pinMode (pinA, INPUT);
  pinMode (pinB, INPUT);
  pinMode (buttonPin, INPUT_PULLUP);
  
  DigiUSB.begin();
  
  aLastState = digitalRead(pinA);
}

void loop() {
  aState = digitalRead(pinA);
  buttonState = digitalRead(buttonPin);

  // If the button is pressed
  if (buttonState == LOW && lastButtonState == HIGH) {
    DigiUSB.write('P'); // Send 'P' for Play/Pause
  }
  lastButtonState = buttonState;

  // If the rotary encoder is turned
  if (aState != aLastState) {
    if (digitalRead(pinB) != aState) {
      DigiUSB.write('+'); // Send '+' for volume up
    } else {
      DigiUSB.write('-'); // Send '-' for volume down
    }
  }
  aLastState = aState;
  
  DigiUSB.refresh();
}
