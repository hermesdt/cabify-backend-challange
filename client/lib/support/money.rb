class Money
  class << self
    def number_to_money(num)
      num / 100.0
    end
  end
end
