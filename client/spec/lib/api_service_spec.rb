require 'json'
require 'spec_helper'

RSpec.describe ApiService do
  subject { ApiService.instance }

  describe "#create_cart" do
    let(:basket_id) { "asdf" }
    
    it "returns the the id of the new basket" do
      body = JSON.dump({"id" => basket_id})
      FakeServer.setup_endpoint(path: '/baskets', status: 201, body: body)

      id = ApiService.instance.create_basket
      expect(id).to eq(basket_id)
    end
  end

  "#add_item"
  "#close_cart"
end
