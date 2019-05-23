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

      id = subject.create_basket
      expect(id).to eq(Basket.new("asdf"))
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

      basket = subject.add_item(basket_id, code)
      expect(basket).to eq(Basket.new(nil, nil, 40.5))
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

      basket = subject.close_basket(basket_id)
      expect(basket).to eq(Basket.new(nil, nil, 20))
    end
  end

  describe "#get_items" do
    it "returns the list of items" do
      body = JSON.dump([
        {"code" => "VOUCHER", "name" => "voucher", "price" => 1.00}
      ])
      FakeServer.stub_endpoint(
        path: "/items",
        status: 200,
        body: body)

      items = subject.get_items
      expect(items).to eq([
        {"code"=>"VOUCHER", "name"=>"voucher", "price"=>1.0}
      ])
    end
  end
end
