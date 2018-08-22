# start from scratch
FROM scratch

# copy binary
COPY wakeUp config.json template/status.html /
COPY template/status.html template/
# COPY template/favicon.ico template/

# tell exposing port
EXPOSE 8000

# run binary
CMD ["/wakeUp", "-v"]

# Thankfully inspired by Daniël van Gils
# https://blog.cloud66.com/how-to-create-the-smallest-possible-docker-image-for-your-golang-application/