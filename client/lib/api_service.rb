require 'json'
require 'singleton'
require 'net/http'

class ApiService
  include Singleton

  FlexibleConfig.use :api_service do |config|
    BASE_URL = config.fetch(:base_url, 'http://localhost:3000').freeze
  end

  def create_basket
    data = request(:post, "/baskets", {})
    Basket.new(data["id"], data["items"], data["promotions"], data["total"])
  end

  def add_item(basket_id, code)
    data = request(:put, "/baskets/#{basket_id}/items", {"code" => code})
    Basket.new(data["id"], data["items"], data["promotions"], data["total"])
  end

  def close_basket(basket_id)
    request(:put, "/baskets/#{basket_id}", {})
  end

  def get_items
    request(:get, "/items")
  end

  private

  def request(method, path, payload=nil)
    uri = URI.parse(BASE_URL + path)
    http = Net::HTTP.new(uri.host, uri.port)

    case method
    when :get then
      request = Net::HTTP::Get.new(uri.request_uri)
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
