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


struct process{

    char *pid;
    char *parent;
    char *name;
    char *state;
    char *memory;
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

	char path[30], fpath[30];
	int c = 1, openFiles = 0;
    while(1){
	strcpy(path, "/proc/");
	strcpy(fpath, "/proc");
	while((dir = readdir(d)) != NULL) {
		if(isdigit(dir->d_name[0])){
			strcat(path, strcat(dir->d_name, "/status"));
			all_p[pos].memory = "0";
			analizeTxt(path);
			pos++;
			//printf("\n");
			strcpy(path, "/proc/");
		}

		DIR *fdd = opendir(fpath);
		struct dirent *fd_dir;
		/*while((fd_dir = readdir(fdd)) != NULL) {
			openFiles++;
		}*/

		closedir(fdd);
		all_p[pos].open_files = openFiles - 2;
		strcpy(fpath, "/proc/");
		//checkOpenFiles(fpath);


	}
	showTable();
	sleep(5);
	clear();
	c++;
    }

	clear();
	return 0;
}

/*void checkOpenFiles(char *fpath){

		int openFiles;
		DIR *fdd = opendir(fpath);
		struct dirent *fd_dir;
		while((fd_dir = readdir(fdd)) != NULL) {
			openFiles++;
		}

		closedir(fdd);
		all_p[pos].open_files = openFiles - 2;
		strcpy(fpath, "/proc/");

}*/

void showTable(){
	float memory;
	printf("|________|________|_______________|_______|____________|__________|\n");
	printf("|PID     | PPID   | STATUS        |THREADS| MEMORY     |OPEN FILES|\n");
	printf("|________|________|_______________|_______|____________|__________|\n");
	for(int i = 0; i < pos; i++){
		//printf("%s\t%s\n", all_p[i].pid, all_p[i].parent);
		memory = atof(all_p[i].memory) / 1000;
		
		printf("|%8s|%8s|%15s|%7s|%12f|%10i|\n", all_p[i].pid, all_p[i].parent, all_p[i].state, all_p[i].threads, memory, all_p[i].open_files);	

	}
	printf("|________|________|_______________|_______|____________|__________|\n");
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
	fprintf(fd, "|________|________|_______________|_______|____________|__________|\n");
	fprintf(fd, "|PID     | PPID   | STATUS        |THREADS| MEMORY     |OPEN FILES|\n");
	fprintf(fd, "|________|________|_______________|_______|____________|__________|\n");
	for(int i = 0; i < pos; i++){
		//printf("%s\t%s\n", all_p[i].pid, all_p[i].parent);
		memory = atof(all_p[i].memory) / 1000;
		
		fprintf(fd, "|%8s|%8s|%15s|%7s|%12f|\n", all_p[i].pid, all_p[i].parent, all_p[i].state, all_p[i].threads, memory);	

	}
	fprintf(fd, "|________|________|_______________|_______|____________|__________|\n");

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
	    if (strcmp(buffer, "vmrs") == 0){

	        if(*t != ' ' && *t != ':' && *t != '\t' && *t != '\n' && *t != 'k' && *t != 'b'){
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
    else if (strcmp(buffer, "vmrs") == 0){
	//printf("# THREADS: %s\n",data);
	all_p[pos].memory = (char *)malloc(sizeof(char) * 50);
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
    close(fp);

}


