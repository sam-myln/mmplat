# Multimedia-platform
Allows you to host locally(globally) files in provided directories, 
supports format selection (only mp3, pdf etc) or format exclusion.
Requires two-factor setup (optional, enabled by default.)
## Prerequisites
- download go requirements
- run ``` mkdir build ```
- run ``` touch .env ``` 
## Build
make build
## Format:
```
mmplat --bind|--addr localhost:8080 --folders dir... [-r|--recursive] [--file-formats formats...] [--log info|debug|error] [--login-data files...] OR [--allow-guest]
```
## Recomended usage
```
mmplat --bind localhost:8080 -r --folders $HOME/videos/movies-with-family --allow-guest
```
