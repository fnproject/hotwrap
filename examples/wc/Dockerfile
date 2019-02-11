## Start of your normal docker file or use an existing image with code in
FROM alpine:latest

# Install hotwrap binary in your container
COPY --from=fnproject/hotwrap:latest /hotwrap /hotwrap

# Copy what you need
# Any old command
CMD  "/usr/bin/wc -l "

ENTRYPOINT ["/hotwrap"]
