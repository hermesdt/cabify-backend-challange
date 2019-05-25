require 'webrick'

RSpec.configure do |config|
  config.before(:suite) do
    FakeServer.instance.start!
  end

  config.after(:suite) do
    FakeServer.instance.stop!
  end

  config.before(:each) do
    FakeServer.instance.reset
  end
end

module WEBrick
  module HTTPServlet
    class ProcHandler
      alias do_PUT do_POST
    end
  end
end

class FakeServer
  include Singleton

  attr_reader :last_request, :server

  def start!
    @server = WEBrick::HTTPServer.new Port: 11334, BindAddress: '0.0.0.0'
    @server_thread = Thread.new { server.start }
  end

  def stop!
    @server.stop
    @server_thread.join
  end

  def reset
    @last_request = nil
  end

  class << self
    def url
      config = FakeServer.instance.server.config
      "http://#{config[:BindAddress]}:#{config[:Port]}"
    end

    def stub_endpoint path: nil, status: nil, body: nil, headers: []
      raise ArgumentError.new("path can't be nil") if path.nil?
      raise ArgumentError.new("status can't be nil") if status.nil?
      raise ArgumentError.new("body can't be nil") if body.nil?
      raise ArgumentError.new("headers must be an array") if !headers.is_a?(Array)

      server = FakeServer.instance.server

      server.mount_proc path do |req, res|
        res.status =  status
        res.body = body
        headers.each do |key, value|
          res.header[key] = value
        end

        FakeServer.instance.instance_eval {
          @last_request = req
        }
      end
    end

    def last_request
      FakeServer.instance.last_request
    end
  end
end
