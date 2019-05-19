require 'json'
require 'singleton'
require 'net/http'

class ApiService
  include Singleton

  BASE_URL = "http://localhost:3000".freeze

  def create_basket
    request(:post, "/baskets", {})["id"]
  end

  def add_item(basket_id, code)
    request(:put, "/baskets/#{basket_id}/items", {"code" => code})["total"]
  end

  def close_basket(basket_id)
    request(:put, "/baskets/#{basket_id}", {})["total"]
  end

  private

  def request(method, path, payload)
    uri = URI.parse(BASE_URL + path)
    http = Net::HTTP.new(uri.host, uri.port)

    case method
    when :post then
      request = Net::HTTP::Post.new(uri.request_uri)
      request.body = JSON.dump(payload)
    when :put then
      request = Net::HTTP::Put.new(uri.request_uri)
      request.body = JSON.dump(payload)
    else
      raise ArgumentError.new("verb #{method} not implemented")
    end

    response = http.request(request)
    handle_response(response)
  end

  def handle_response(response)
    response.value
    JSON.parse(response.body)
  end
end
