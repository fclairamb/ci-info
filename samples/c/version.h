#ifndef __VERSION_H__
#define __VERSION_H__

typedef struct {
    const char 
        *version, 
        *commit_hash, 
        *commit_date, 
        *commit_smart, 
        *build_date;
} build_info_t;

extern build_info_t build_info;

#endif
