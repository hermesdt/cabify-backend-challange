module UI
  class Console
    def initialize(basket_id = nil, api_service = ApiService.instance)
      @api_service = api_service
      @basket_id = basket_id
      @total = nil
    end

    def run
      if @basket_id.nil?
        @basket_id = @api_service.create_basket
      end

      done = false
      while !done
        print_message
        key = read_key
        if execute_key(key) == :quit
          done = true
        end
      end
    end

    private

    def read_key
      begin
        print "Choose one option: "
        k = STDIN.gets.chomp
        Integer(k)
      rescue ArgumentError => e
        puts "Incorrect option"
        retry
      end
    end

    def execute_key(key)
      action_idx = key.to_i
      if action_idx < 0
        return
      end

      if action_idx < get_items.size
        @total = @api_service.add_item(@basket_id, get_items[action_idx]["code"])
      end

      if action_idx == get_items.size
        @api_service.close_basket(@basket_id)
        return :quit
      end
    end

    def print_message
      items_str = get_items.each_with_index.map do |item, idx|
        "\t#{idx}: #{item["name"]}"
      end.join("\n")

      message = <<-MESSAGE

Welcome! Your basket id is #{@basket_id}
* Add Item:
#{items_str}
#{get_items.size}: Close basket\n
MESSAGE

      puts message
    end

    def get_items
      @items ||= @api_service.get_items
    end
  end
end