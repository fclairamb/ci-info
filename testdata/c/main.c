#include <stdio.h>

#include "version.h"

int main(int argc, char **argv) {
    printf("Version:           %s\n", build_info.version);
    printf("CommitHash:        %s\n", build_info.commitHash);
    printf("CommitDate:        %s\n", build_info.commitDate);
    printf("CommitTagOrBranch: %s\n", build_info.commitTagOrBranch);
    printf("Date:              %s\n", build_info.buildDate);
    return 0;
}
