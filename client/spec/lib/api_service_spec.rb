require 'json'
require 'spec_helper'

RSpec.describe ApiService do
  subject { ApiService.instance }

  before(:each) do
    stub_const("ApiService::BASE_URL", FakeServer.url)
  end

  describe "#create_cart" do
    let(:basket_id) { "asdf" }
    
    it "returns the id of the new basket" do
      body = JSON.dump({"id" => basket_id})
      FakeServer.stub_endpoint(path: '/baskets', status: 201, body: body)

      id = ApiService.instance.create_basket
      expect(id).to eq(basket_id)
    end
  end

  describe "#add_item" do
    let(:basket_id) { "asdf" }
    let(:code) { "VOUCHER" }
    let(:total) { 40.5 }
    
    it "returns the total after addig the item" do
      body = JSON.dump({"total" => total})
      FakeServer.stub_endpoint(
        path: "/baskets/#{basket_id}/items",
        status: 200,
        body: body)

      basket_total = ApiService.instance.add_item(basket_id, code)
      expect(basket_total).to eq(total)
    end
  end

  describe "#close_basket" do
    let(:basket_id) { "asdf" }
    let(:total) { 20 }
    
    it "returns the total after closing the basket" do
      body = JSON.dump({"total" => total})
      FakeServer.stub_endpoint(
        path: "/baskets/#{basket_id}",
        status: 200,
        body: body)

      basket_total = ApiService.instance.close_basket(basket_id)
      expect(basket_total).to eq(total)
    end
  end
end
