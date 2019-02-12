#include<sys/types.h>
#include<sys/stat.h>
#include <fcntl.h>  
#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>

int main(int argc, char *argv[])
{
    int file, r;
    char *prog = argv[0];
    char *f;
    f = (char *) malloc(1);

    if (argc == 1){ 
        //filecopy(0, 1);
    }else{
        while (--argc > 0)
	        if ((file = open(*++argv, O_RDONLY)) < 0) {
                fprintf(stderr, "%s: can′t open %s\n",
			prog, *argv);
                return 1;
            } else {
                while((r = read(file,f,1)) > 0){
		    write(1,f,r);
		}
                close(file);
            }
    }
    if (ferror(stdout)) {
        fprintf(stderr, "%s: error writing stdout\n", prog);
        return 2;
    }
    
    return 0;
}
