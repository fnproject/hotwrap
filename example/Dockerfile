FROM fnproject/hotwrap:latest as hotwrap

## Start of your normal docker file
FROM alpine:latest

# Install hotwrap binary in your container
COPY --from=hotwrap /hotwrap /hotwrap

ENTRYPOINT ["/hotwrap", "awk", "'{print $1}'"]
