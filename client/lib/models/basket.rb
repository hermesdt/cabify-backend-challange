require 'ostruct'

class Basket < Struct.new(:id, :items, :promotions, :total)
  def items
    @items ||= super&.yield_self do |items|
      items.map { |item| Item.from_json(item) }
    end || []
  end

  def promotions
    @promotions ||= super&.yield_self do |promotions|
      promotions.map { |promo| Promotion.from_json(promo) }
    end || []
  end

  def total
    @total ||= Money.number_to_money(super)
  end
end
