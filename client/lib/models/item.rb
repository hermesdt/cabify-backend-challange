require 'ostruct'

class Item < Struct.new(:code, :name, :price)

  class << self
    def from_json json
      Item.new(json["code"], json["name"], json["price"])
    end
  end

  def price
    @price ||= Money.number_to_money(super)
  end
end
