FROM fnproject/hotwrap:latest as hotwrap

## Start of your normal docker file
FROM alpine:latest

# Install hotwrap binary in your container
COPY --from=hotwrap /hotwrap /hotwrap
# Install whatever you want
# just any old command
CMD wc -l


ENTRYPOINT ["/hotwrap"]
