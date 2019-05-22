module UI
  class Console
    def initialize(api_service: ApiService.instance)
      @api_service = api_service
    end

    def banner
    end
  end
end