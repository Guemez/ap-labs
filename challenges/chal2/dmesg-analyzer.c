#include <stdio.h>
#include <string.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdbool.h>

#define REPORT_FILE "report.txt"



void analizeLog(char *logFile, char *report);
void addToStruct(char *type, char *info);
void processLine(char *p);

struct logFile{
    char *type;
    char *info[100];
    int index;
};

struct logFile logs[1000];
int pos = 0;

int main(int argc, char **argv) {

    if (argc < 2) {
	printf("Usage:./dmesg-analizer.o logfile.txt\n");
	return 1;
    }

    analizeLog(argv[1], REPORT_FILE);

    return 0;
}

void processLine(char *p){
    
    char *t = (char *)malloc(sizeof(char) * 1); 
    char *aux = (char *)malloc(sizeof(char) * 1);
    char *x = (char *)malloc(sizeof(char) * 1);
    char *type = (char *)malloc(sizeof(char) * 1000);
    char *limbo = (char *)malloc(sizeof(char) * 1000);
    char *info = (char *)malloc(sizeof(char) * 1000);
    bool flag_info = true;
    bool found_type = false;
    bool already_checked = false;
    
    for(int i = 0; i < 15; i++){
        (*x = *p);
        strcat(info, x);
        p++;
    }    

    while((*t = *p) != '\0') {
        *aux = *++p; 
        --p;
        //strcat(type, t);
        if(*t == ':' && *aux == ' ' && !already_checked){
            ++p;
            *t = *p;
            strcpy(type, limbo);
            flag_info = false;
            found_type = true;
            already_checked = true;
        }
        if(flag_info){
            strcat(limbo, t);
        }else{
            strcat(info, t);
        }
        
        ++p;  
    }
    if(!found_type){
        strcat(type, "General");
        strcat(info, limbo);
    }
    addToStruct(type, info);

    
}

void addToStruct(char *type, char *info){

    for(int i = 0; i < pos; i++){
	//printf("COMPARING %s TO %s\n", logs[i].type, type);
	if(*logs[i].type == *type){
	    //printf("IT ALREADY EXISTS\n");
	    logs[i].info[logs[i].index] = (char *)malloc(sizeof(char)*1000);
  	    strcpy(logs[i].info[logs[i].index], info);
    	    logs[i].index++;
	    return;

	}
    }
    logs[pos].type = (char *)malloc(sizeof(char)*1000);
    logs[pos].info[0] = (char *)malloc(sizeof(char)*1000);
    strcpy(logs[pos].type, type);
    strcpy(logs[pos].info[0], info);
    logs[pos].index = 1;
    pos++;

}

void analizeLog(char *logFile, char *report) {
    printf("Generating Report from: [%s] log file\n", logFile);

    // Implement your solution here.
    FILE *fInput, *fOutput;
    int c;
    char *f = (char *)malloc(sizeof(char)*1);

    if((fInput = open(logFile, O_RDONLY)) < 0){
	printf("cannot open file");
    }

    /*if((fOutput = open(report, O_CREAT | O_WRONLY | O_RDONLY, 0640)) < 0){
	printf("cannot open file");
    }*/
    fOutput = fopen(report, "w");



    char *line = (char *)malloc(sizeof(char)*1000);

    while((c = read(fInput, f, 1)) > 0){
	if( *f == '\n'){
	    strcat(line, "\n\0");
	    processLine(line);
	    *line = NULL;
	    
	} else{
	    strcat(line, f);
	}


	//write(1, f, c);
	//write(fOutput, f, c);
    }
    for(int k = 0; k < pos; k++){
    	//printf("%s\n", logs[k].type);
	fprintf(fOutput, logs[k].type);
	fprintf(fOutput, "\n");
	write(1, logs[k].type, 100);
	write(1, "\n", 1);
	for(int j = 0; j < logs[k].index; j++){
	    fprintf(fOutput, "\t");
	    fprintf(fOutput, logs[k].info[j]);
	    //printf("\t%s", logs[k].info[j]);
	    write(1, "\t", 1);
	    write(1, logs[k].info[j], 500);
	}
    }
    //fflush(fOutput);


    close(fInput);
    fclose(fOutput);

    printf("Report is generated at: [%s]\n", report);
}
