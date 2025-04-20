FROM ruby:3.2-alpine

RUN apk add --no-cache build-base
WORKDIR /app

COPY . /app
RUN gem install sinatra rackup puma

EXPOSE 4567
CMD ["ruby", "app.rb"]
