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
#include <stdlib.h>

void analizeTxt(char *logFile);

void clear();
static void sigHandler(int sig);
void showTable();
void checkOpenFiles();

int algo = 0;


struct process{

    char pid[50];
    char parent[50];
    char name[50];
    char state[50];
    char memory[50];
    char threads[50];
    int open_files;

};

struct process all_p[5000];
int pos = 0;

int main(){


	if (signal(SIGINT, sigHandler) == SIG_ERR){
            perror("Error handling signals\n");
	}

	struct dirent *dir;

	DIR *d = opendir("/proc/");

	char path[30], fpath[30];
	int c = 1;
    while(1){
	strcpy(path, "/proc/");
	strcpy(fpath, "/proc/");
	while((dir = readdir(d)) != NULL) {
		if(isdigit(dir->d_name[0])){
			strcat(path, dir->d_name);
			strcat(path, "/status");
			//printf("path: %s\n", path);
			strcpy(all_p[pos].memory, "0");
			analizeTxt(path);
			pos++;
			//printf("\n");
			strcpy(path, "/proc/");
			strcat(fpath, dir->d_name);
			strcat(fpath, "/fd");
			checkOpenFiles(fpath);
			algo++;
		}

	}
	showTable();
	sleep(3);
	clear();
	c++;
    }
    clear();
    return 0;
}

void checkOpenFiles(char *fpath){

		int openFiles;
		DIR *fdd = opendir(fpath);
		struct dirent *fd_dir;
		while((fd_dir = readdir(fdd)) != NULL) {
			openFiles++;
		}

		closedir(fdd);
		all_p[pos].open_files = openFiles - 2;
		strcpy(fpath, "/proc/");

}

void showTable(){
	float memory;
	printf("|________|________|_______________|_______|__________|__________|________________________________________|\n");
	printf("|PID     | PPID   | STATUS        |THREADS| MEMORY   |OPEN FILES|NAME                                    |\n");
	printf("|________|________|_______________|_______|__________|__________|________________________________________|\n");
	for(int i = 0; i < pos; i++){
		if(all_p[i].pid == 0){
			continue;
		}
		memory = atof(all_p[i].memory) / 1000;
		
		printf("|%8s|%8s|%15s|%7s|%8.4f M|%10i|%40s|\n", all_p[i].pid, all_p[i].parent, all_p[i].state, all_p[i].threads, memory, all_p[i].open_files, all_p[i].name);	

	}
	printf("|________|________|_______________|_______|__________|__________|________________________________________|\n");
	return;
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
    float memory;

	FILE *fd = fopen(filename, "w");

	if(fd == NULL){
		perror("Can't open file");
		exit(1);
	}
	fprintf(fd, "|________|________|_______________|_______|__________|__________|________________________________________|\n");
	fprintf(fd, "|PID     | PPID   | STATUS        |THREADS| MEMORY   |OPEN FILES|NAME                                    |\n");
	fprintf(fd, "|________|________|_______________|_______|__________|__________|________________________________________|\n");
	for(int i = 0; i < pos; i++){
		if(all_p[i].pid == 0){
			continue;
		}
		memory = atof(all_p[i].memory) / 1000;
		
		fprintf(fd, "|%8s|%8s|%15s|%7s|%8.4f M|%10i|%40s|\n", all_p[i].pid, all_p[i].parent, all_p[i].state, all_p[i].threads, memory, all_p[i].open_files, all_p[i].name);	

	}
	fprintf(fd, "|________|________|_______________|_______|__________|__________|________________________________________|\n");

	fclose(fd);
    
	printf("Saved in file %s\n", filename);
}

void clear() {
  printf("\e[1;1H\e[2J"); 
}



void processLine(char *line){
    
    //char *t = (char *)malloc(sizeof(char) * 1); 
    char *buffer = (char *)malloc(sizeof(char) * 4);
    char *data = (char *)malloc(sizeof(char) * 80);
    bool flag = true;
    int p = 0;

	int j = 0;

	for(int i = 0; line[i] != '\0'; i++){
	
		if(line[i] == ':'){
	    flag = false;
	}

	if(flag){
	    if(p<4){
		//printf("%s", t);
		buffer[i] = line[i];
		p++;
	    }
	}

	if(flag == false){

	    if (strcmp(buffer, "pid") == 0 || strcmp(buffer, "ppid") == 0 ||
strcmp(buffer, "name") == 0 || strcmp(buffer, "stat") == 0 || strcmp(buffer, "thre") == 0){

	        if(line[i] != ' ' && line[i] != ':' && line[i] != '\t' && line[i] != '\n'){
		    //printf("char:%s\n", t);
		    data[j] = line[i];
			j++;
	        }	

	    }
	    if (strcmp(buffer, "vmrs") == 0){

	        if(line[i] != ' ' && line[i] != ':' && line[i] != '\t' && line[i] != '\n' && line[i] != 'k' && line[i] != 'b'){
		    data[j] = line[i];
			j++;
	        }	

	    }
		
	}

	}


    if (strcmp(buffer, "pid") == 0){
	strcpy(all_p[pos].pid, data);
    }else if (strcmp(buffer, "ppid") == 0){
	strcpy(all_p[pos].parent, data);
    }else if (strcmp(buffer, "name") == 0){
	strcpy(all_p[pos].name, data);
    }else if (strcmp(buffer, "stat") == 0){
	strcpy(all_p[pos].state, data);
    }else if (strcmp(buffer, "thre") == 0){
	strcpy(all_p[pos].threads, data);
    }else if (strcmp(buffer, "vmrs") == 0){
	strcpy(all_p[pos].memory, data);
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

    char line[500]; 
    memset(line, 0, 500);


    while((c = read(fp, f, 1)) > 0){
	*f = tolower(*f);

	if( *f == '\n'){
	    //printf("%s\n", line);
	    strcat(line, "\n\0");
	    processLine(line);
	    memset(line, 0, 500);
	    
	    
	}else{
	    strcat(line, f);
	}

    }
    close(fp);

}


