#ifndef __VERSION_H__
#define __VERSION_H__

typedef struct {
    const char 
        *version, 
        *commitHash, 
        *commitDate, 
        *commitTagOrBranch, 
        *buildDate;
} build_info_t;

extern build_info_t build_info;

#endif
