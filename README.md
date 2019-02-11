# Fn HotWrap

HotWrap is a beta tool that lets you create Fn functions based on conventional unix command line tools  (like shell commands or anything else you can invoke in a terminal) while also taking advantage of Fn's streaming event model inside your container.
 

HotWrap implements the Fn FDK contract via a command wrapper `hotwrap` this command wrapper then invokes a command for each event your function receives. 

HotWrap sends the body of incoming events to your command via STDIN and reads the response from STDOUT 

# Using HotWrap 

HotWrap works using Fn's existing docker support:  


suppose you have a Dockerfile for a command that works on the CLI: 

```
FROM alpine:latest

# just any old command 
COMMAND /usr/bin/wc -l   

```



Add HotWrap to your container as follows: 

Dockerfile:
```
## Start of your normal docker file 
FROM alpine:latest

# Install hotwrap binary in your container 
COPY --from=fnproject/hotwrap:latest  /hotwrap /hotwrap 

# just any old command 
CMD /usr/bin/wc -l   

# update entrypoint to use hotwrap, this will wrap your command 
ENTRYPOINT ["/hotwrap"]
```

Create a func.yaml as follows: 
```
schema_version: 20180708
name: example
version: 0.0.1
```

Deploy the function to an Fn server with app name `hotdemo`: 

```bash
fn deploy --app hotdemo

```

Invoke the function: 


```bash

echo $'some\nlines\nof\ntext' | fn invoke hotdemo example 

4
```
 
The Input passed to the function will be passed on stdin and any output that the code returns on stdout will be returned as function output. 

As with other functions anything sent to stderr will be passed to the functions logs. 

HotWrap is a portable  statically linked binary that should work in any linux container.  It invokes commands in a shell and requires at least "/bin/sh". 
 
 
 # Accessing headers from function call 
 
 You an receive invocation headers from function calls as environment variables. 
 
 All incoming function heaers are transposed into environment variables using the following rules: 
 
 * Must start with Fn- (this includes http trigger/gateway headers (see below))
 * Capitalized 
 * s/-/_/ 
 * Disambiguated by taking the first matching value where multiple headers are present or two headers resolve to the same variable. 
 

 For HTTP request headers and details received from triggers these are mapped as follows: 
 
 * FN_HTTP_H_<Header Name> 
 * FN_HTTP_METHOD 
 * FN_HTTP_REQUEST_URL
 
 e.g. for a call to a trigger: 
 
 ```
 GET /my/trigger HTTP/1.1 
 My-Header: foo 
 Accept: * 
 Accept: application/xml 

```

you function will see the following environment 

```
FN_HTTP_METHOD=GET
FN_HTTP_REQUEST_URL=http://tld.com/my/trigger
FN_HTTP_H_MY_HEADER=foo
FN_HTTP_H_ACCEPT=*

``` 