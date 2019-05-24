require 'ostruct'

class Promotion < Struct.new(:name, :total_discount)
  class << self
    def from_json json
      Promotion.new(json["name"], json["total_discount"])
    end
  end

  def total_discount
    @total_discount ||= Money.number_to_money(super)
  end
end
