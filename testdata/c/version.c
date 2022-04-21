#include "version.h"

build_info_t build_info = {
    .version = "{{ .Version }}",
    .commitHash = "{{ .CommitHash }}", 
    .commitDate = "{{ .CommitDate }}", 
    .commitTagOrBranch = "{{ .CommitTagOrBranch }}",
    .buildDate = "{{ .BuildDate }}",
};
