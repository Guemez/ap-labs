#include <stdio.h>
#include <stdlib.h>

int main(void)
{
    FILE *ft;
    char const *name = "test.c";
    int ch;
    ft = fopen(name, "r+");
    if (ft == NULL)
    {
        fprintf(stderr, "cannot open target file %s\n", name);
        exit(1);
    }
    while ((ch = fgetc(ft)) != EOF)
    {
	fprintf("%d", ch);
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
