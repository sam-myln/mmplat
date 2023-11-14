# Multimedia-platform
Allows you to host locally(globally) files in provided directories, 
supports format selection (only mp3, pdf etc) or format exclusion.
Requires two-factor setup (optional, enabled by default.)
## Format:
```
mmplat --bind|--addr localhost:8080 --folders dir... [-r|--recursive] [--disable-two-factor] [--file-formats formats...] [--log info|debug|error] [--login-data files...]
```
## Recomended usage
```
mmplat --bind localhost:8080 -r --folders $HOME/videos/movies-with-family --login-data authorize_users
```