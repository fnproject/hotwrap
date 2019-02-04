# Fn HotWrap tool 

HotWrap is a beta tool that lets you create "Hot" Fn functions based on conventional unix command line tools (like shell commands or anything else you can invoke in a terminal) while also taking advantage of Fn's streaming event model inside your container.
 

HotWrap implements the Fn FDK contract via a command wrapper `hotwrap` this command wrapper then invokes a command for each event your function receives. 

HotWrap sends the body of incoming events to your command via STDIN and reads the response from STDOUT 

# Using Hotwrap 



HotWrap works best if you use a `docker` type function: 


suppose you have a Dockerfile for a command that works on the CLI: 

```
FROM alpine:latest

# just any old command 
COMMAND /usr/bin/wc -l   

```



Add HotWrap to your container as follows: 

Dockerfile:
```
# Pull the HotWrap container  as a build dependency 
FROM fnproject/hotwrap:latest as hotwrap

## Start of your normal docker file 
FROM alpine:latest

# just any old command 
CMD /usr/bin/wc -l   

# Install hotwrap binary in your container 
COPY --from=hotwrap /hotwrap /hotwrap 
ENTRYPOINT ["/hotwrap"]
```


func.yaml:
```
schema_version: 20180708
name: example
version: 0.0.1
```

Deploy the function to an Fn server with app name `hotdemo` and invoke the function.

```bash

echo $'some\nlines\nof\ntext' | fn invoke hotdemo example 

4
```
 
You can pass value to the application using HTTP headers. 

```bash

curl -H "abc:123" -H "xyz:456" http://localhost:8080/t/app/func

# The function will be able to consume those values by using the correspsonding environment variables
# FN_HEADER_Abc=123
# FN_HEADER_Xyz=456

```




Hotwrap is a portable statically linked binary that should work in any linux container. 