#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <assert.h>
#include <sys/stat.h>
#include <fcntl.h>
#include <unistd.h>
#include <errno.h>


char *reduce(char *input, int l) {
    int i = 0;
    for (; i < l;) {
        if (abs(input[i] - input[i + 1]) == 32) {
            memmove(input + i, input + i + 2, l - i);
            l -= 2;
            input[l] = '\0';
            if (i>0) --i;
            continue;
        }
        i++;
    }
    return input;
}


// modifies input bytes
char *removeChar(char *input, char lower) {
    if (lower <= 'Z') {
        lower += 32;
    }
    char upper = lower - 32;
    int i = 0;
    int l = strlen(input);
    for (; i < l - 1;) {
        char c = input[i];
        if (c == upper || c == lower) {
            // squanch it
            memmove(input + i, input + i + 1, l - i);
            l--;
            continue;
        }
        i++;
    }
    return input;
}


long readFile(const char *path, char **buf) {
    int fd = -1;
    FILE * fh = NULL;
    if ((fd = open(path, O_RDONLY)) < 0) {
        return -1;
    }
    fh = fdopen(fd, "r");
    if (fh == NULL) {
        close(fd);
        return -1;
    }

    struct stat st;
    fstat(fd, &st);
    int size = st.st_size;
    *buf = calloc(size+1, sizeof(char));

    if (fgets(*buf, size, fh) == NULL) {
        size = -1;
        free(*buf);
    }
    fclose(fh);
    return size;

}

int main() {

    const char *path="./input.txt";
    char *input;
    int inputL = readFile(path, &input);

    if (inputL < 0) {
        printf("error reading file %s: %s\n", path, strerror(errno));
        exit(1);
    }
    char *scratch = malloc(inputL);
    memcpy(scratch, input, inputL);
    char *reduced = reduce(scratch, inputL - 1);
    int part1len = strlen(reduced);
    printf("part 1 len: %d\n", part1len);
    int part2len = part1len;
    char *scratch2 = calloc(part1len+1, sizeof(char));
    for(char c = 'a'; c <= 'z'; ++c) {
        memcpy(scratch2, reduced, part1len);
        char *reduced2 = removeChar(scratch2, c);
        reduced2 = reduce(reduced2, strlen(reduced2));
        int l = strlen(reduced2);
        if (l < part2len) {
            part2len = l;
        }
    }
    printf("part 2: %d\n", part2len);
    free(scratch);
    free(scratch2);
    free(input);
}

