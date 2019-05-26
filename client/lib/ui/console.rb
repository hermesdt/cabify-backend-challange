module UI
  class Console
    def initialize(api_service: ApiService.instance, input: STDINT, output: STDOUT)
      @api_service = api_service
      @output = output
      @input = input
      @basket = nil
    end

    def start
      if @basket.nil?
        @basket = @api_service.create_basket
      end

      done = false
      while !done
        done = run
      end
    end

    def run
      begin
        print_message
        key = read_key
        case execute_key(key)
        when :quit then
          return true
        when :unknown then
          @output.puts "Incorrect option"
        end
      rescue StandardError => e
        @output.puts "Received exception #{e.to_s}"
        @output.puts e.backtrace.join("\n")
      end

      false
    end

    private

    def read_key
      begin
        @output.print "Choose one option: "
        k = @input.gets.chomp
        Integer(k)
      rescue ArgumentError => e
        @output.puts "Incorrect option"
        retry
      end
    end

    def execute_key(key)
      action_idx = key.to_i
      if action_idx < 0
        return :unknown
      end

      if action_idx < get_items.size
        @basket = @api_service.add_item(@basket.id, get_items[action_idx]["code"])
      end

      if action_idx == get_items.size
        @api_service.close_basket(@basket.id)
        return :quit
      end

      return :unknown
    end

    def print_message
      items_str = get_items.each_with_index.map do |item, idx|
        "\t#{idx}: #{item.name} (#{item.price}€)"
      end.join("\n")

      basket_str = @basket.items.
        group_by { |item| item.code }.
        map { |code, items| "\t- #{items[0].name} (#{items[0].price}€) x #{items.size}" }.
        join("\n")
      
      promotions_str = @basket.promotions&.map do |promo|
        "\t- #{promo.name}. Earned #{promo.total_discount}€"
      end.join("\n")


      message = <<-MESSAGE

Welcome! Your basket id is #{@basket.id}
* Add Item:
#{items_str}
#{get_items.size}: Close basket

Basket:
#{basket_str}
Promos:
#{promotions_str}
Total: #{@basket.total}€
MESSAGE

      @output.puts message
    end

    def get_items
      @items ||= @api_service.get_items.map { |item| Item.from_json(item) }
    end
  end
end
