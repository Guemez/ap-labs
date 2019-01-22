#include <stdio.h>
#include <stdlib.h>

#define t 1
#define f 0

int main(void)
{
    FILE *ft;
    char file_name[50];
    printf("Enter the name of the file: ");
    scanf("%[^\n]%*c", file_name);
    char const *name = file_name;
    int ch, prev;
    int quote = f;
    ft = fopen(name, "r+");
    if (ft == NULL)
    {
        fprintf(stderr, "cannot open target file %s\n", name);
        exit(1);
    }
    ch = prev = fgetc(ft);
    while (ch != EOF || prev != EOF){
	prev = ch;
	if (ch == '"' && quote == f && prev != '\\'){
		quote = t;
	}
	else if (ch == '"' && quote == t && prev != '\\'){
		quote = f;
	}
	ch = fgetc(ft);

	if (prev == '/' && ch == '*' && quote == f){
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
	}

	if (prev == '/' && ch == '/' && quote == f){
		fseek(ft, -2, SEEK_CUR);
            	fputc(' ',ft);
		fseek(ft, 0, SEEK_CUR);
		prev = fgetc(ft);
		ch = fgetc(ft);
		while ( prev != '\n' || ch != '\n'){
			fseek(ft, -2, SEEK_CUR);
            		fputc(' ',ft);
			fseek(ft, 0, SEEK_CUR);
			prev = fgetc(ft);
			ch = fgetc(ft);
		}
		fseek(ft, -2, SEEK_CUR);
            	fputc(' ',ft);
		fputc('\n',ft);
		fseek(ft, 0, SEEK_CUR);
	}

    }
    fclose(ft);
    return 0;
}
