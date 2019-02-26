#include <stdio.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdbool.h>


void analizeTxt(char *logFile);
void addToStruct(char *type, char *line);
void processWord(char *p);
int currLine = 1;

struct dict{
    char *word;
    char *line[1000];
    int index;
};

struct dict all_words[10000];
int pos = 0;

int main(int argc, char **argv) {

    if (argc < 2) {
	printf("Usage:./cross-ref text.txt\n");
	return 1;
    }
    analizeTxt(argv[1]);

    return 0;
}

void processWord(char *word){
    if(strcmp(word, "the") == 0 || strcmp(word, "i") == 0 || strcmp(word, "of") == 0 || strcmp(word, "and") == 0 || strcmp(word, "to") == 0 || strcmp(word, "in") == 0 || strcmp(word, "with") == 0 || strcmp(word, "a") == 0){
	return;
    }
    char *ln = (char *)malloc(sizeof(char)*100);
    char *text[32];
    sprintf(text, "%d", currLine);
    ln = text;
    addToStruct(word, ln);
    
}


void addToStruct(char *word, char *info){
    for(int i = 0; i < pos; i++){
	if(strcmp(all_words[i].word, word) == 0){
	    all_words[i].line[all_words[i].index] = (char *)malloc(sizeof(char)*1000);
  	    strcpy(all_words[i].line[all_words[i].index], info);
    	    all_words[i].index++;
	    return;

	}
    }
    all_words[pos].word = (char *)malloc(sizeof(char)*1000);
    all_words[pos].line[0] = (char *)malloc(sizeof(char)*1000);
    strcpy(all_words[pos].word, word);
    strcpy(all_words[pos].line[0], info);
    all_words[pos].index = 1;
    pos++;

}

void analizeTxt(char *logFile) {

    FILE *fInput;
    int c;
    char *f = (char *)malloc(sizeof(char)*1);

    if((fInput = open(logFile, O_RDONLY)) < 0){
	printf("cannot open file");
    }

    char *line = (char *)malloc(sizeof(char)*100);


    while((c = read(fInput, f, 1)) > 0){
	*f = tolower(*f);
	if(*f == ',' || *f == '.' || *f == ';' || *f == '"' || *f == '\'' || *f == '-'
	|| *f == ')' || *f == '(' || *f == '*' || *f == '!'){
	    continue;
	}
	if(*f == ' '){
	    if(*line != NULL){
	        processWord(line);
	    }
	    *line = NULL;
	}
	else if( *f == '\n'){
	    if(*line != NULL){
	        processWord(line);
	    }
	    currLine = currLine + 1;
	    *line = NULL;
	    
	}else{
	    strcat(line, f);
	}

    }
    for(int k = 0; k < pos; k++){
	write(0, all_words[k].word, 100);
	write(0, "\n", 1);
	for(int j = 0; j < all_words[k].index; j++){
	    write(0, " ", 1);
	    write(0, all_words[k].line[j], 500);
	}
	write(0, "\n", 1);
    }


    close(fInput);

}
