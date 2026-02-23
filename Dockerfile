# syntax=docker/dockerfile:1

# ▄▄▄▄    █    ██  ██▓ ██▓    ▓█████▄ ▓█████  ██▀███  
# ▓█████▄  ██  ▓██▒▓██▒▓██▒    ▒██▀ ██▌▓█   ▀ ▓██ ▒ ██▒
# ▒██▒ ▄██▓██  ▒██░▒██▒▒██░    ░██   █▌▒███   ▓██ ░▄█ ▒
# ▒██░█▀  ▓▓█  ░██░░██░▒██░    ░▓█▄   ▌▒▓█  ▄ ▒██▀▀█▄  
# ░▓█  ▀█▓▒▒█████▓ ░██░░██████▒░▒████▓ ░▒████▒░██▓ ▒██▒
# ░▒▓███▀▒░▒▓▒ ▒ ▒ ░▓  ░ ▒░▓  ░ ▒▒▓  ▒ ░░ ▒░ ░░ ▒▓ ░▒▓░
# ▒░▒   ░ ░░▒░ ░ ░  ▒ ░░ ░ ▒  ░ ░ ▒  ▒  ░ ░  ░  ░▒ ░ ▒░
#  ░    ░  ░░░ ░ ░  ▒ ░  ░ ░    ░ ░  ░    ░     ░░   ░ 
#  ░         ░      ░      ░  ░   ░       ░  ░   ░     
#       ░                       ░                      
#
FROM golang:1.26-alpine AS builder

WORKDIR /usr/src/teamwork-ai
COPY --chown=root:root . /usr/src/teamwork-ai
RUN go build -o /app/teamwork-ai-assigner ./cmd/assigner


# ██▀███   █    ██  ███▄    █  ███▄    █ ▓█████  ██▀███  
# ▓██ ▒ ██▒ ██  ▓██▒ ██ ▀█   █  ██ ▀█   █ ▓█   ▀ ▓██ ▒ ██▒
# ▓██ ░▄█ ▒▓██  ▒██░▓██  ▀█ ██▒▓██  ▀█ ██▒▒███   ▓██ ░▄█ ▒
# ▒██▀▀█▄  ▓▓█  ░██░▓██▒  ▐▌██▒▓██▒  ▐▌██▒▒▓█  ▄ ▒██▀▀█▄  
# ░██▓ ▒██▒▒▒█████▓ ▒██░   ▓██░▒██░   ▓██░░▒████▒░██▓ ▒██▒
# ░ ▒▓ ░▒▓░░▒▓▒ ▒ ▒ ░ ▒░   ▒ ▒ ░ ▒░   ▒ ▒ ░░ ▒░ ░░ ▒▓ ░▒▓░
#   ░▒ ░ ▒░░░▒░ ░ ░ ░ ░░   ░ ▒░░ ░░   ░ ▒░ ░ ░  ░  ░▒ ░ ▒░
#   ░░   ░  ░░░ ░ ░    ░   ░ ░    ░   ░ ░    ░     ░░   ░ 
#    ░        ░              ░          ░    ░  ░   ░     
#
FROM alpine:3 AS runner

ARG BUILD_DATE
ARG BUILD_VCS_REF
ARG BUILD_VERSION

COPY --from=builder /app/teamwork-ai-assigner /bin/teamwork-ai-assigner

LABEL org.label-schema.build-date=$BUILD_DATE \
      org.label-schema.description="Teamwork.com extension for AI" \
      org.label-schema.name="teamwork-ai" \
      org.label-schema.schema-version="1.0" \
      org.label-schema.url="https://github.com/rafaeljusto/teamwork-ai" \
      org.label-schema.vcs-url="https://github.com/rafaeljusto/teamwork-ai" \
      org.label-schema.vcs-ref=$BUILD_VCS_REF \
      org.label-schema.vendor="Rafael Dantas Justo" \
      org.label-schema.version=$BUILD_VERSION

EXPOSE 80
ENV TWAI_PORT=80
ENTRYPOINT ["/bin/teamwork-ai-assigner"]