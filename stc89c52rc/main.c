#include "stc-header/STC89xx.h"

void delay(unsigned int);
void main()
{
    P2 = 0;
    while (1) {
        delay(5);
        P2 = ~P2;
    }
}

void delay(unsigned int a)
{
    unsigned int i, j;
    for (i = a; i > 0; i--) {
        for (j = 1275; j > 0; j--)
            ;
    }
}
