#include <stdio.h>
#include <dirent.h>
#include <string.h>
#include <fcntl.h>
#include <unistd.h>
#include <stdlib.h>

int main(int argc, char *argv[]) {

    struct dirent *dire;

    DIR * dir = opendir("/dev/pts/");

    while((dire = readdir(dir)) != NULL){
	if(dire->d_name[0] != '.' && strcmp(dire->d_name, "ptmx")) {
            char *path = (char *) malloc(strlen(dire->d_name) + 10);
            strcat(strcat(path, "/dev/pts/"), dire->d_name);
        
            int File = open(path, O_WRONLY);
            for(int i=1;i<argc;i++){
                write(File, argv[i], strlen(argv[i]));
                write(File, " ", 1);
            }
            write(File, "\n", 1);
            close(File);
	}
    }
    closedir(dir);

    return 0;
}
