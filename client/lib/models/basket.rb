require 'ostruct'

class Basket < Struct.new(:id, :items, :total)
  def items
    @items ||= super&.yield_self do |items|
      items.map { |item| Item.from_json(item) }
    end || []
  end

  def total
    @total ||= super&.yield_self do |total|
      Integer(total) / 100.0
    end
  end
end
