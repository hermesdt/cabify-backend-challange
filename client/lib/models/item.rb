class Item < Struct.new(:code, :name, :price)

  class << self
    def from_json json
      Item.new(json["code"], json["name"], json["price"])
    end
  end

  def price
    super.yield_self { |price| Integer(price)/100.0 }
  end
end
