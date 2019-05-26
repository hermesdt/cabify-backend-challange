require 'spec_helper'

RSpec.describe UI::Console do
  let(:output) { StringIO.new }
  let(:items) do
    [
      {"code" => "VOUCHER", "name" => "Cabify voucher", "price" => 100},
      {"code" => "TSHIRT", "name" => "Cabify tshirt", "price" => 2000}
    ]
  end
  let(:promotions) do
    [
      {"name" => "Promo 1", "total_discount" => 250}
    ]
  end
  let(:input) { double(:input, gets: "0") }
  let(:basket) do
    Basket.new("basket-id", items, promotions, 1000)
  end
  let(:avaiable_items) do
    [
      Item.new("VOUCHER", "Cabify voucher", 100),
      Item.new("TSHIRT", "Cabify tshirt", 2000)
    ]
  end

  let(:api_service) do
    double(:api_service,
      get_items: avaiable_items,
      create_basket: basket)
  end

  subject do
    UI::Console.new(api_service: api_service, input: input, output: output)
  end

  describe "#run" do
    before(:each) do
      subject.instance_variable_set("@basket", basket)
    end

    it "prints the message" do
      expect(api_service).to receive(:close_basket)
      expect(subject).to receive(:read_key).and_return(2)
      subject.run
      expect(output.string).to eq(
      "\nWelcome! Your basket id is basket-id\n" \
      "* Add Item:\n" \
      "\t0: Cabify voucher (1.0€)\n" \
      "\t1: Cabify tshirt (20.0€)\n" \
      "2: Close basket\n" \
      "\n" \
      "Basket:\n" \
      "\t- Cabify voucher (1.0€) x 1\n"+
      "\t- Cabify tshirt (20.0€) x 1\n"+
      "Promos:\n" \
      "\t- Promo 1. Earned 2.5€\n" \
      "Total: 10.0€\n")
    end

    it "adds item when selecting an item index" do
      expect(api_service).to receive(:add_item).with(basket.id, avaiable_items[1].code)
      expect(input).to receive(:gets).and_return(1.to_s)
      subject.run
    end

    it "closes basket when selecting items.size + 1 option" do
      expect(api_service).to receive(:close_basket)
      expect(input).to receive(:gets).and_return(avaiable_items.size.to_s)
      subject.run
    end

    it "prints message for incorrect input" do
      expect(input).to receive(:gets).ordered.and_return("asdf")
      expect(input).to receive(:gets).ordered.and_return("1")
      allow(api_service).to receive(:add_item)
      subject.run
      expect(output.string).to include("Choose one option: Incorrect option")
    end
  end
end
