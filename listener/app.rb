# app.rb
require 'sinatra'
require 'json'

before do
  request.body.rewind
  @request_payload = request.body.read
end

helpers do
  def log_request
    puts "[#{request.request_method}] #{request.path}"
    puts "Headers: #{request.env.select { |k, _| k.start_with? 'HTTP_' }}"
    puts "Query Params: #{params.to_json}"
    puts "Body: #{@request_payload}"
  end
end

['get', 'post', 'put', 'delete'].each do |method|
  send(method, '/*') do
    log_request
    content_type :json
    { message: "Handled #{request.request_method}", path: request.path }.to_json
  end
end

set :bind, '0.0.0.0'
set :port, 4567
