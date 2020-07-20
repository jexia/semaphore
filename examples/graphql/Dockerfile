FROM jexiacom/semaphore:latest

WORKDIR /app

ENV LOG_LEVEL=debug

COPY *.hcl ./
COPY proto/ ./proto/
COPY annotations/ ./annotations/

ENTRYPOINT ["semaphore", "daemon", "-f", "config.hcl"]