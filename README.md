# Fn HotWrap tool 

HotWrap is a beta tool that lets you create "Hot" Fn functions based on conventional unix command line tools 
(like shell commands or anything else you can invoke in a terminal) while also taking advantage of Fn's streaming event model inside your container.
 

Hot wrap implements the Fn FDK contract via a command wrapper `hotwrap` this command wrapper then invokes a command for each event your function receives. 

HotWrap sends the body of incoming events to your command via STDIN and reads the response from STDOUT 

# Using Hotwrap 



Hotwrap works best if you use a `docker` type function: 


suppose you have a Dockerfile for a command that works on the CLI: 

```
FROM ubuntu:latest

# just any old command 
COMMAND /usr/bin/wc -l   

```



Add Hotwrap to your container as follows: 

Dockerfile:
```
# Pull the hotwrap container  as a build dependency 
FROM fnproject/hotwrap:latest as hotwrap

## Start of your normal docker file 
FROM ubuntu:latest

# just any old command 
COMMAND /usr/bin/wc -l   

# Install hotwrap binary in your container 
COPY --from hotwrap /hotwrap /hotwrap 
ENTRYPOINT ["/hotwrap"]
```


func.yaml:
```
name: example
version: 0.0.1
format: json
```


```bash

echo "some\nlines\nof\ntext" | fn run 

4
```
 

Hotwrap is a portable  statically linked binary that should work in any linux container. 