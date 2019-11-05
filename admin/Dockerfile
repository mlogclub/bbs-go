# builder stage
FROM node:10.16-slim AS builder

# Install build dependencies
RUN apt-get -qq update && \
    apt-get -qq install -y --no-install-recommends \
      build-essential \
      git \
      openssh-client \
      locales \
    && rm -rf /var/lib/apt/lists/*

# Set locale: https://daten-und-bass.io/blog/fixing-missing-locale-setting-in-ubuntu-docker-image/
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend noninteractive locales && \
    update-locale LANG=en_US.UTF-8
ENV LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8

# Update timezone
ENV TZ=Asia/Shanghai

WORKDIR /app

# Install dependencies
COPY package.json* ./
RUN npm i

FROM builder AS assets

# Copy application code
COPY . .

# Build project
RUN npm run build

# Remove unused dependencies
#ARG NPM_PRUNE_FLAGS
#RUN npm prune ${NPM_PRUNE_FLAGS}

# application stage
FROM node:10.16-slim AS application

# Install runtime dependencies
RUN apt-get -qq update && \
    apt-get -qq install -y --no-install-recommends \
      curl \
      locales \
    && rm -rf /var/lib/apt/lists/*

# Set locale: https://daten-und-bass.io/blog/fixing-missing-locale-setting-in-ubuntu-docker-image/
RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend noninteractive locales && \
    update-locale LANG=en_US.UTF-8
ENV LANG=en_US.UTF-8 LC_ALL=en_US.UTF-8

# Update timezone
ENV TZ=Asia/Shanghai

EXPOSE 8080

HEALTHCHECK --start-period=10s \
            --interval=15s \
            --timeout=5s \
            --retries=3 \
            CMD curl -sSf http://localhost:8080/heartbeat || exit 1

WORKDIR /app

# Install serve
RUN npm i -g serve

# node_modules from former assets stage
COPY --from=assets /app/dist /app/dist
#CMD ["serve", "-s", "dist", "-l", "8080"] # dont know why it dont not work,

COPY --from=assets /app/node_modules /app/node_modules
COPY . .
CMD ["npm", "run", "serve"]
