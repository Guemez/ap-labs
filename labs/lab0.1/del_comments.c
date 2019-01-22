#include <stdio.h>
#include <stdlib.h>

#define t 1
#define f 0

int main(void)
{
    FILE *ft;
    char const *name = "test.c";
    int ch, prev;
    int quote = f;
    ft = fopen(name, "r+");
    if (ft == NULL)
    {
        fprintf(stderr, "cannot open target file %s\n", name);
        exit(1);
    }
    ch = prev = fgetc(ft);
    while (ch != EOF){
	prev = ch;
	
	if (ch == '"' && quote == f && prev != '\\'){
		quote = t;
	}
	else if (ch == '"' && quote == t && prev != '\\'){
		quote = f;
	}
	ch = fgetc(ft);
	printf("inside quote?: %d\n", quote);
	if (prev == '/' && ch == '*' && quote == f){
		printf("comment found!\n");
		fseek(ft, -2, SEEK_CUR);
            	fputc(' ',ft);
		fseek(ft, 0, SEEK_CUR);
		prev = fgetc(ft);
		ch = fgetc(ft);
		while ( prev != '*' || ch != '/'){
			fseek(ft, -2, SEEK_CUR);
            		fputc(' ',ft);
			fseek(ft, 0, SEEK_CUR);
			prev = fgetc(ft);
			ch = fgetc(ft);
		}
		fseek(ft, -2, SEEK_CUR);
            	fputc(' ',ft);
		fputc(' ',ft);
		fseek(ft, 0, SEEK_CUR);
		printf("done with comment\n");
	}
        /*if (ch == 'i')
        {
            fseek(ft, -1, SEEK_CUR);
            fputc('a',ft);
            fseek(ft, 0, SEEK_CUR);
        }*/
    }
    fclose(ft);
    return 0;
}
