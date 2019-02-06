//strlib.c
#include <stdio.h>
int strlen(char *str){
    int count = 0;
    while(str[count] != '\0'){
        count++;
    }
    return count;
}

char* stradd(char *origin, char *addition){
    int l = strlen(origin) + strlen(addition) + 1;
    char *p = malloc(l);
    int i = 0;
    while (i < strlen(origin)){
        p[i] = origin[i];
        i++;
    }
    int j = 0;
    while(j < strlen(addition)){
	p[i] = addition[j];
	i++;
	j++;
    }
    p[i] = '\0';
    return p;
}

int strfind(char *origin, char *substr){
    int len1 = strlen(origin);
    int len2 = strlen(substr);
    //printf("str1: %s, str2: %s\n", origin, substr);
    //printf("l1: %d, l2: %d\n", len1, len2);
    int i = 0;
    int j = 0;
    int flag = 0;
    while(i < len1){
	if(*(origin + i) == *(substr + j)){
	    flag = 1;
	    i++;
	    j = i;
	    for (int k = 1; k < len2 - 1; k++){
		if(*(origin + i) != *(substr + k)){
		    flag = 0;
		    i = j;
		    break;
		}
	    }
	    
	}
	i++;
    }
    return flag;
}


