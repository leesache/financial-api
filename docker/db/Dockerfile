# docker/db/Dockerfile
FROM postgres:15

# Install dependencies
RUN apt-get update && \
    apt-get install -y curl bash && \
    rm -rf /var/lib/apt/lists/*

RUN curl -fsSL \
    https://raw.githubusercontent.com/pressly/goose/master/install.sh |\
    sh 

# Copy migrations folder into the container
COPY migrations /migrations

# Copy entrypoint script
COPY docker/db/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Expose the default PostgreSQL port
EXPOSE 5432

# Run the entrypoint script when the container starts
ENTRYPOINT ["/entrypoint.sh"]