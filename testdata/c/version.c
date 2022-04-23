#include "version.h"

build_info_t build_info = {
    .version = "{{ .Version }}",
    .commit_hash = "{{ .CommitHash }}", 
    .commit_date = "{{ .CommitDate }}", 
    .commit_smart = "{{ .CommitSmart }}",
    .build_date = "{{ .BuildDate }}",
};
