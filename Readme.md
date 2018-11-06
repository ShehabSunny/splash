## `gosplash` is a simple cli tool for changing the wallpaper dynamically and easily from hundrands of thousands of images of unsplash.com.

### Prerequisite:
You need to have a developer account of unsplash.com in order to use this tool.  
1. Go to https://unsplash.com/developers, sign up, create an application and get the access key. 
2. Create a file at `$HOME/.gosplash/.env` and add the line `ACCESS_KEY=<access_key>` (replace with your access key) in the file.


### Build and Install:
Pull the source code with
`go get github.com/ShehabSunny/gosplash`  

Build with:
```
make build
```
or 
```
go build .
``` 

Install with:
```
go install
```
---
  

### Usage: 
This command will pull a random image from the API and set as desktop wallpaper.
```
gosplash
```

This command will pull a random image from the API searched with the topic provided with -t flag and set it as desktop wallpaper.
```
gosplash -t sunset
```
or 
```
gosplash -t "new york"
```
