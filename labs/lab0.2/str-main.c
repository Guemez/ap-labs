#include <stdio.h>

int strlen(char *str);
char* stradd(char *origin, char *addition);
int strfind(char *origin, char *substr);


int main( int argc, char *argv[]){
    if(argc == 1){
        printf("no parameters provided");
	return 0;
    }
    int i = strlen(argv[1]);
    printf("Initial lenght: %d\n", i);
    char *s = stradd(argv[1], argv[2]);
    printf("New string: %s\n", s);
    int j = strfind(argv[1], argv[3]);
    printf("Substring  was found: ");
    if(j == 1){
	printf("yes\n");
    }
    else{
	printf("no\n");
    }
}
