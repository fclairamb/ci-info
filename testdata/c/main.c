#include <stdio.h>

#include "version.h"

int main(int argc, char **argv) {
    printf("Version:           %s\n", build_info.version);
    printf("CommitHash:        %s\n", build_info.commit_hash);
    printf("CommitDate:        %s\n", build_info.commit_date);
    printf("CommitSmart:       %s\n", build_info.commit_smart);
    printf("Date:              %s\n", build_info.build_date);
    return 0;
}
