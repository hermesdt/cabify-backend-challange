require 'webrick'

RSpec.configure do |config|
  config.before(:suite) do
    FakeServer.start!
  end
end

module WEBrick
  module HTTPServlet
    class ProcHandler
      alias do_PUT do_POST
    end
  end
end

module FakeServer
  def start!
    if Thread.current[:fake_server]
      return
    end

    server = WEBrick::HTTPServer.new Port: 11334, BindAddress: '0.0.0.0'
    server_thread = Thread.new { server.start }

    Thread.current[:fake_server] = server
    Thread.current[:fake_server_thread] = server_thread
  end
  module_function :start!

  def url
    config = Thread.current[:fake_server].config
    "http://#{config[:BindAddress]}:#{config[:Port]}"
  end
  module_function :url

  def stub_endpoint path: nil, status: nil, body: nil, headers: []
    raise ArgumentError.new("path can't be nil") if path.nil?
    raise ArgumentError.new("status can't be nil") if status.nil?
    raise ArgumentError.new("body can't be nil") if body.nil?
    raise ArgumentError.new("headers must be an array") if !headers.is_a?(Array)

    server = Thread.current[:fake_server]

    server.mount_proc path do |req, res|
      res.status =  status
      res.body = body
      headers.each do |key, value|
        res.header[key] = value
      end
    end

    server
  end
  module_function :stub_endpoint
end
