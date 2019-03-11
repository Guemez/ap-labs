#include "logger.h"

int main(int argc, char const *argv[])
{	

	    // default logging
    infof("INFO Message %d", 1);
    warnf("WARN Message %d", 2);
    errorf("ERROR Message %d", 2);

    // stdout logging
    initLogger("stdout");
    infof("INFO Message %d", 1);
    warnf("WARN Message %d", 2);
    errorf("ERROR Message %d", 2);

    // syslog logging
    initLogger("syslog");
    infof("INFO Message %d", 1);
    warnf("WARN Message %d", 2);
    errorf("ERROR Message %d", 2);
	/*int myInt = 10;
    double myDouble = 15.2;
    float myFloat = 12.4f;
    char myChar = 'H';
    char *myString = "Hello World";

	initLogger("syslog");

	infof("%cello, this is general information\n", myChar);	
	warnf("%cello, this is a warning with code: %f\n", myChar, myDouble);
	errorf("%cello, this is an error with code: %d\n", myChar, myInt);
	panicf("%s, this is a panic: %f\n", myString, myFloat);*/
	
	return 0;
}