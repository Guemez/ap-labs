#include <sys/inotify.h>
#include <limits.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
#include <string.h>
#include "logger.h"
#define _XOPEN_SOURCE 500   /* See feature_test_macros(7) */
#include <ftw.h>

int inotifyFd, wd;

int c1, c2;
int cook;
char *path;

static int
display_info(const char *fpath, const struct stat *sb,
             int tflag, struct FTW *ftwbuf)
{
    if( tflag == FTW_D){
        wd = inotify_add_watch(inotifyFd, fpath, IN_CREATE | IN_DELETE | IN_MOVED_FROM | IN_MOVED_TO );
        if (wd == -1)
            errorf("inotify_add_watch");
        //infof("Watching %s using wd %d\n", fpath, wd);
    }

    return 0;           /* To tell nftw() to continue */
}

int callnftw(){
    if (nftw(path, display_info, 20, 0)
            == -1) {
             errorf("nftw");
             exit(0);
            }
}

static void             /* Display information from inotify_event structure */
displayInotifyEvent(struct inotify_event *i)
{

    if (i->cookie > 0){
        //warnf("cookie =%4d; ", i->cookie);
    }
    if (i->mask & IN_CREATE){
        infof("CREATED");
        if (i->len > 0)
            infof("        name = %s\n", i->name);
    }
    if (i->mask & IN_DELETE){
        infof("DELETED");
        if (i->len > 0)
            infof("        name = %s\n", i->name);
    }
    if (i->mask & IN_MOVED_FROM){
	cook = i->cookie;
        //printf("%i", cook);
    }
    if (i->mask & IN_MOVED_TO){   
	if(cook == i->cookie){
	    //printf("%i", cook);
	    infof("RENAMED");
	cook = 0;
	if (i->len > 0)
            infof("        name = %s\n", i->name);
	}
    }
    callnftw();

    printf("\n");

    /*if (i->len > 0)
        printf("        name = %s\n", i->name);*/
}

#define BUF_LEN (10 * (sizeof(struct inotify_event) + NAME_MAX + 1))

int
main(int argc, char *argv[])
{
    
    int j;
    char buf[BUF_LEN] __attribute__ ((aligned(8)));
    ssize_t numRead;
    char *p;
    //char *fp;
    struct inotify_event *event;

    if (argc < 2 || strcmp(argv[1], "--help") == 0)
        errorf("%s pathname...\n", argv[0]);


    inotifyFd = inotify_init();                 /* Create inotify instance */
    if (inotifyFd == -1)
        errorf("inotify_init");

    path = argv[1];
	
        wd = inotify_add_watch(inotifyFd, path, IN_CREATE | IN_DELETE | IN_MOVED_FROM | IN_MOVED_TO );
        if (wd == -1)
            errorf("inotify_add_watch");

	callnftw();



    

    for (;;) {                                  /* Read events forever */
        numRead = read(inotifyFd, buf, BUF_LEN);
        if (numRead == 0)
            errorf("read() from inotify fd returned 0!");

        if (numRead == -1)
            errorf("read");

        //printf("Read %ld bytes from inotify fd\n", (long) numRead);


        /* Process all of the events in buffer returned by read() */

        for (p = buf; p < buf + numRead; ) {
            event = (struct inotify_event *) p;
            displayInotifyEvent(event);



            p += sizeof(struct inotify_event) + event->len;
        }
    }
    exit(0);
}
