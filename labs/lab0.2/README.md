## Compile
To compile use the following commands:
```
gcc -c str-main.c -o str-main.o
gcc -c strlib.c -o strlib.o
gcc str-main.o strlib.o -o str-final.o
```

## Run
To run type the following:
```./str-final.o str1 str2 str3```

Note: The substring is looked for between str1 and str3
