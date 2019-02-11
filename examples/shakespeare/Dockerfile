## Start of your normal docker file or use an existing image with code in
FROM alpine:latest

# Install hotwrap binary in your container
COPY --from=fnproject/hotwrap:latest /hotwrap /hotwrap


ADD sonnets.txt /
ADD hamlet.txt /
ADD shakespeare.sh  /
CMD  "/shakespeare.sh"

ENTRYPOINT ["/hotwrap"]
