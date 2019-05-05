FROM node:10.11-alpine as dev

COPY ./yarn.lock ./package.json /app/
WORKDIR /app/
RUN npm install -g yarn && \
    yarn install

COPY . /app/

RUN yarn build

ENTRYPOINT yarn serve

FROM nginx:1.15.5-alpine as prod

COPY --from=Dev /app/dist/ /usr/share/nginx/html/
