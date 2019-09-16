FROM node:12

WORKDIR /app

COPY package.json /app
COPY yarn.lock /app
RUN yarn install

COPY . /app
RUN yarn compile

ENTRYPOINT ["yarn"]
CMD ["start"]
