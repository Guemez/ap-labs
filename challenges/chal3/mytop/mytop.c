#include <stdio.h>
#include <string.h>
#include <dirent.h>
#include <ctype.h>
#include <sys/stat.h>
#include <sys/types.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>
#include <stdbool.h>
#include <signal.h>
#include <time.h>
#include <stdbool.h>

void analizeTxt(char *logFile);

void clear();
static void sigHandler(int sig);


struct process{

    char *pid;
    char *parent;
    char *name;
    char *state;
    int memory;
    char *threads;
    int open_files;
    char *ptm;

};

struct process all_p[5000];
int pos = 0;

int main(){


	if (signal(SIGINT, sigHandler) == SIG_ERR){
            perror("Error handling signals\n");
	}

	struct dirent *dir;

	DIR *d = opendir("/proc/");

	char path[30];
	int c = 1;

    while(1){
	strcpy(path, "/proc/");
	while((dir = readdir(d)) != NULL) {
		if(isdigit(dir->d_name[0])){
			strcat(path, strcat(dir->d_name, "/status"));
			analizeTxt(path);
			pos++;
			//printf("\n");
			strcpy(path, "/proc/");
		}

	}
	printf("|_______|_______|_______________|_______|\n");
	printf("|PID    | PPID  | STATUS        |THREADS|\n");
	printf("|_______|_______|_______________|_______|\n");
	for(int i = 0; i < pos; i++){
		//printf("%s\t%s\n", all_p[i].pid, all_p[i].parent);
		write(1, "|", 1);
		write(1, all_p[i].pid, 30);
		write(1, "\t", 1);
		write(1, "|", 1);
		write(1, all_p[i].parent, 30);
		write(1, "\t", 1);
		write(1, "|", 1);
		write(1, all_p[i].state, 30);
		write(1, "\t", 1);
		write(1, "|", 1);
		write(1, all_p[i].threads, 30);
		write(1, "\t", 1);
		write(1, "|", 1);
		write(1, "\n", 1);

	}
	printf("|_______|_______|_______________|_______|\n");
	printf("%i\n", c);
	sleep(5);
	clear();
	c++;
    }

	clear();
	return 0;
}



static void
sigHandler(int sig)
{
    char filename[150];
    strcpy(filename, "mytop-status-");
    
    time_t rawtime;
    struct tm * timeinfo;
    time ( &rawtime );
    timeinfo = localtime ( &rawtime );
    strcat(filename, asctime(timeinfo));
    strcat(filename, ".txt");

	FILE *fd = fopen(filename, "w");

	if(fd == NULL){
		perror("Can't open file");
		exit(1);
	}

	fprintf(fd, "|_______|_______|_______________|_______|\n");
	fprintf(fd, "|PID    | PPID  | STATUS        |THREADS|\n");
	fprintf(fd, "|_______|_______|_______________|_______|\n");
	for(int i = 0; i < pos; i++){
		fprintf(fd, "|%s\t|%s\t|%s\t|%s\t|\n", all_p[i].pid, all_p[i].parent, all_p[i].state, all_p[i].threads);
	}
	fprintf(fd, "|_______|_______|_______________|_______|\n");

	fclose(fd);
    
	printf("Saved in file %s\n", filename);
}

void clear() {
  printf("\e[1;1H\e[2J"); 
}



void processLine(char *line){
    
    char *t = (char *)malloc(sizeof(char) * 1); 
    char *buffer = (char *)malloc(sizeof(char) * 4);
    char *data = (char *)malloc(sizeof(char) * 50);
    bool flag = true;
    int p = 0;

    while((*t = *line) != '\0'){
	//printf("char: %s\n", t);
	if(*t == ':'){
	    flag = false;
	}

	if(flag){
	    if(p<4){
		strcat(buffer, t);
		p++;
	    }
	}

	if(flag == false){

	    if (strcmp(buffer, "pid") == 0 || strcmp(buffer, "ppid") == 0 ||
strcmp(buffer, "name") == 0 || strcmp(buffer, "stat") == 0 || strcmp(buffer, "thre") == 0){

	        if(*t != ' ' && *t != ':' && *t != '\t' && *t != '\n'){
		    //printf("char:%s\n", t);
		    strcat(data, t);
	        }	

	    }
	}

	++line;
    }

    //printf("buffer: %s, data: %s\n", buffer, data);

    if (strcmp(buffer, "pid") == 0){
	//printf("PID: %s\n",data);
	all_p[pos].pid = (char *)malloc(sizeof(char) * 50);
	strcpy(all_p[pos].pid, data);
	//printf("inside struct: %s\n", all_p[pos].pid);
    }else if (strcmp(buffer, "ppid") == 0){
	//printf("PPID: %s\n",data);
	all_p[pos].parent = (char *)malloc(sizeof(char) * 50);
	strcpy(all_p[pos].parent, data);
	//printf("inside struct: %s\n", all_p[pos].parent);
    }else if (strcmp(buffer, "name") == 0){
	all_p[pos].name = (char *)malloc(sizeof(char) * 50);
	strcpy(all_p[pos].name, data);
	//printf("NAME: %s\n",data);

	//printf("name in struct: %s\n", all_p[pos].name);
    }else if (strcmp(buffer, "stat") == 0){
	//printf("STATUS: %s\n",data);
	all_p[pos].state = (char *)malloc(sizeof(char) * 50);
	strcpy(all_p[pos].state, data);
	//printf("inside struct: %s\n", all_p[pos].state);
    }else if (strcmp(buffer, "thre") == 0){
	//printf("# THREADS: %s\n",data);
	all_p[pos].threads = (char *)malloc(sizeof(char) * 50);
	strcpy(all_p[pos].threads, data);
	//printf("inside struct: %s\n", all_p[pos].threads);
    }

}



void analizeTxt(char *logFile) {

    //FILE *fIn;
    int fp;
    int c;
    char *f = (char *)malloc(sizeof(char)*1);


    if((fp = open(logFile, O_RDONLY)) < 0){
	printf("cannot open file");
    }

    char* line = (char *)malloc(sizeof(char)*100);


    while((c = read(fp, f, 1)) > 0){
	*f = tolower(*f);

	if( *f == '\n'){
	    strcat(line, "\n\0");
	    processLine(line);
	    free(line);
	    line = (char *)malloc(sizeof(char)*100);
	    //*line = NULL;
	    
	}else{
	    strcat(line, f);
	}

    }

}


