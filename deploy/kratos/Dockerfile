FROM oryd/kratos:v25.4.0

USER root

# Copy configuration files
COPY kratos.prod.yml /etc/config/kratos/kratos.prod.yml
COPY identity.schema.json /etc/config/kratos/identity.schema.json

# Set ownership to ory user (default user in kratos image)
RUN chown -R ory:ory /etc/config/kratos

USER ory
