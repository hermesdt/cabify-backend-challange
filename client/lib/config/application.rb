Bundler.require(:default, ENV.fetch("RUBY_ENV", :development).to_sym)

lib_files = [File.expand_path(File.dirname(__FILE__)), "..", "**", "*.rb"].join(File::SEPARATOR)
Dir[lib_files].each { |f| require f }