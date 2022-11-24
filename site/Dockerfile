FROM node:16.15.0-alpine

ENV HOST 0.0.0.0
ENV APP_HOME /app/bbs-go

WORKDIR "$APP_HOME"
COPY . ./
COPY nuxt.config.docker.js ./nuxt.config.js

# install
RUN npm config set registry https://registry.npm.taobao.org \
  && npm config set sass_binary_site https://npm.taobao.org/mirrors/node-sass
RUN npm install

# build
ENV NODE_ENV production
RUN npm run build

EXPOSE 3000
CMD ["npm", "run", "start"]
