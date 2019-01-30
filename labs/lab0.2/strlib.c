//strlib.c
#include <stdio.h>
int strlen(char *str){
    int count = 0;
    while(str[count] != '\0'){
        count++;
    }
    return count;
}

/**char stradd(char *origin, char *addition){

}*/

int strfind(char *origin, char *substr){

}

int main( int argc, char *argv[]){
    if(argc == 1){
        printf("no parameters provided");
	return 0;
    }
    int i = strlen(argv[1]);
    printf("len: %d\n", i);
}
