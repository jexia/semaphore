FROM jexiacom/maestro:latest

WORKDIR /app

ENV LOG_LEVEL=debug

COPY *.hcl ./
COPY proto/ ./proto/
COPY annotations/ ./annotations/

ENTRYPOINT ["maestro", "run", "-f", "config.hcl"]