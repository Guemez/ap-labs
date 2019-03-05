
#include <stdio.h>
#include <stdarg.h>
#include <stdlib.h>

#define RESET		0
#define BRIGHT 		1
#define DIM		2
#define UNDERLINE 	3
#define BLINK		4
#define REVERSE		7
#define HIDDEN		8

#define BLACK 		0
#define RED		1
#define GREEN		2
#define YELLOW		3
#define BLUE		4
#define MAGENTA		5
#define CYAN		6
#define	WHITE		7

void textcolor(int attr, int fg, int bg)
{	char command[13];

	/* Command is the control command to the terminal */
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	printf("%s", command);
}

void bt(void) {
    int c, i;
    void *addresses[10];


    c = backtrace(addresses, 10);
    printf("BEGIN\n", c);
    for(i = 0; i < c; i++) {
        printf("%d: %X\n", i, (int)addresses[i]);
    }
    printf("\nEND\n");
}


int infof(const char *format, ...){
   textcolor(BRIGHT, WHITE, BLACK);
   va_list arg;
   int done;
   va_start (arg, format);
   done = vfprintf (stdout, format, arg);
   va_end (arg);
   printf("\033[0m");
   return done;
}

int warnf(const char *format, ...){
   textcolor(BRIGHT, YELLOW, BLACK);
   va_list arg;
   int done;
   va_start (arg, format);
   done = vfprintf (stdout, format, arg);
   va_end (arg);
   printf("\033[0m");
   return done;
}

int errorf(const char *format, ...){
   textcolor(BRIGHT, RED, BLACK);
   va_list arg;
   int done;
   va_start (arg, format);
   done = vfprintf (stdout, format, arg);
   va_end (arg);
   printf("\033[0m");
   return done;
}

int panicf(const char *format, ...){
   textcolor(BRIGHT, CYAN, BLACK);
   va_list arg;
   int done;
   va_start (arg, format);
   done = vfprintf (stdout, format, arg);
   va_end (arg);
   bt();
   printf("\033[0m");
   return done;
}
