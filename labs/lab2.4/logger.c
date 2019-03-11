
#include <stdio.h>
#include <stdarg.h>
#include <execinfo.h>
#include <syslog.h>
#include <string.h>

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

int loggType = 0; // stdout = 0, syslog = 1;

int initLogger(char *logType){
   if(strcmp(logType, "stdout") == 0){
      loggType = 0;
   } else if (strcmp(logType, "syslog") == 0){
      loggType = 1;
   } else{
      printf("invalid log type\n");
   }
   return 0;
}

void textcolor(int attr, int fg, int bg)
{	char command[13];

	/* Command is the control command to the terminal */
	sprintf(command, "%c[%d;%d;%dm", 0x1B, attr, fg + 30, bg + 40);
	printf("%s", command);
}

void bt(void) {
    int c, i;
    void *addresses[10];
    char** names;
    c = backtrace(addresses, 10);
    names = backtrace_symbols(addresses, c);
    if(loggType == 0){
       printf("BEGIN\n", c);
       for(i = 0; i < c; i++) {
         //printf("%d: %X", i, (int)addresses[i]);
         printf("%s\n", names[i]);
       }
       printf("\nEND\n");
    } else{
       openlog ("LOG-PANIC", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);

       syslog (LOG_ERR, names[i]);

       closelog ();
    }
    
}


int infof(const char *format, ...){
   textcolor(BRIGHT, WHITE, BLACK);
   va_list arg;
   int done;
   va_start (arg, format);
   if(loggType == 0){
      done = vfprintf (stdout, format, arg);
   } else{
      openlog ("LOG-INFO", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_INFO, format, arg);
		closelog();
   }
   
   va_end (arg);
   printf("\033[0m");
   return done;
}

int warnf(const char *format, ...){
   textcolor(BRIGHT, YELLOW, BLACK);
   va_list arg;
   int done;
   va_start (arg, format);
   if(loggType == 0){
      done = vfprintf (stdout, format, arg);
   } else{
      openlog ("LOG-WARN", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_WARNING, format, arg);
		closelog();
   }
   
   va_end (arg);
   printf("\033[0m");
   return done;
}

int errorf(const char *format, ...){
   textcolor(BRIGHT, RED, BLACK);
   va_list arg;
   int done;
   va_start (arg, format);
   if(loggType == 0){
      done = vfprintf (stdout, format, arg);
   } else{
      openlog ("LOG-ERROR", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_ERR, format, arg);
		closelog();
   }
   va_end (arg);
   printf("\033[0m");
   return done;
}

int panicf(const char *format, ...){
   textcolor(BRIGHT, CYAN, BLACK);
   va_list arg;
   int done;
   va_start (arg, format);
   if(loggType == 0){
      done = vfprintf (stdout, format, arg);
   } else{
      openlog ("LOG-PANIC", LOG_CONS | LOG_PID | LOG_NDELAY, LOG_LOCAL1);
		vsyslog(LOG_ERR, format, arg);
		closelog();
   }
   va_end (arg);
   bt();
   printf("\033[0m");
   return done;
}
