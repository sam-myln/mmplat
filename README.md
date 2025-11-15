# Multimedia-platform
Allows you to host locally( or globally) files in provided directories, 
supports format selection (only mp3, pdf etc) or format exclusion.

## Build and install
```
make build
```
## SYNOPSIS
```bash
mmplat [OPTION]... [PATH]...

OPTIONS:
  --test
        enable test mode
  --port=PORT
        set service port (default: 8080)
  -r, --recursive
        process directories recursively
  --file-formats=FORMAT[,FORMAT]...
        specify allowed file formats, default to all non-executable of appropiate mime
  --log=LEVEL
        set logging level (debug, info, warn, error)
```
## DESCRIPTION

Process files at the specified path with various configuration options.
! File path canNOT be above the executable root.
```bash
# good
.
├── folder1
├── folder2
│   └── subfolder1
└── mmplat


# Not good
.
└── above
    ├── folder1
    ├── folder2
    │   └── subfolder1
    └── mmplay
```
## OPTIONS

**--test**
: Run in localhost

**--port=PORT**
: Specify port number, default :8080

**-r**, **--recursive**
: Process directories recursively

**--file-formats=FORMAT,...**
: Comma-separated file formats: mp4, pdf, etc

**--log-level=LEVEL**
: Log level: info, debug, or error

## ARGUMENTS

**<files_paths>...**
: Target directory or directories

## EXAMPLES

```bash
mmplat --test
mmplat --test -r dir1 dir2
mmplat --port 8080 -r --file-formats txt,json
mmplat --log-level debug media
```