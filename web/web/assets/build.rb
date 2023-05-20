#!/usr/bin/env ruby
# frozen_string_literal: true

#
# build.rb: Build style.min.css and script.min.js.
#

# load libraries
require 'fileutils'
require 'json'
require 'logger'
require 'openssl'
require 'open-uri'
require 'uri'
require 'zip'

class Config
  attr :log_level, :zip_url, :zip_hash

  # bulma zip URL
  BULMA_ZIP_URL = 'https://github.com/jgthms/bulma/releases/download/0.9.4/bulma-0.9.4.zip'

  # expected sha256 hash of bulma zip
  BULMA_ZIP_HASH = '781fb4d195c83fb37090a4b1b66f2a8fc72ea993dc708e6e9bfdce30088fa635'

  def initialize
    @log_level = ENV.fetch('BUILD_LOG_LEVEL', 'info')
    @zip_url = ENV.fetch('BUILD_BULMA_ZIP_URL', BULMA_ZIP_URL)
    @zip_hash = ENV.fetch('BUILD_BULMA_ZIP_HASH', BULMA_ZIP_HASH)
  end
end

class Builder
  def self.run
    new(Config.new).run
  end

  def initialize(config)
    @config = config.freeze

    # init logger
    @log = ::Logger.new(STDERR)
    @log.level = config.log_level
  end

  def run
    # work in temporary directory
    Dir.mktmpdir do |tmp_dir|
      # switch to temp dir
      Dir.chdir(tmp_dir)

      # download bulma-0.9.4.zip
      zip_path = File.join(tmp_dir, File.basename(@config.zip_url))
      @log.debug('run') { "zip_path = #{zip_path}" }
      fetch(zip_path, @config.zip_url, @config.zip_hash)

      # create css working directory
      css_dir = File.join(tmp_dir, 'css')
      @log.debug('run') { "creating css dir #{css_dir}" }
      Dir.mkdir(css_dir)

      # expand 
      @log.debug('run') { "extracting sass files" }
      File.open(zip_path, 'rb') do |fh|
        Zip::File.open(fh) do |zip|
          # find sass entries
          entries = zip.entries.select { |e| e.name =~ /bulma\/sass/ }
          @log.debug('run') { JSON(entries.map { |e| e.name }) }

          # expand entries
          entries.each do |e|
            if e.directory?
              # create directory
              dir_path = File.join(css_dir, e.name)
              @log.debug('run') { "creating dir #{dir_path}" }
              FileUtils.mkdir_p(dir_path)
            else
              # uncompress file
              dst_path = File.join(css_dir, e.name)
              @log.debug('run') { "extracting #{e.name} to #{dst_path}" }
              IO.copy_stream(e.get_input_stream, dst_path)
            end
          end
        end
      end
      @log.debug('run') { "done extracting sass files" }

      # copy bookman.sass into css dir
      src_sass_path = File.join(__dir__, 'bookman.sass')
      dst_sass_path = File.join(css_dir, 'bulma', 'bookman.sass')
      @log.debug('run') { "copy #{src_sass_path} to #{dst_sass_path}" }
      IO.copy_stream(src_sass_path, dst_sass_path)

      # build path to style.min.css
      dst_style_path = File.join(__dir__, '../public/style.min.css')
      src_cmd = ['sass', '--sourcemap=none', dst_sass_path]
      dst_cmd = ['minify', '--mime', 'text/css', '-o', dst_style_path]
      @log.debug('run') { JSON({ src_cmd: src_cmd, dst_cmd: dst_cmd }) }

      # generate public/style.min.css
      IO.popen(src_cmd, 'r') do |src_io|
        IO.popen(dst_cmd, 'w') do |dst_io|
          IO.copy_stream(src_io, dst_io)
        end
      end

      # generate public/script.min.js
      src_script_path = File.join(__dir__, '../assets/script.js')
      dst_script_path = File.join(__dir__, '../public/script.min.js')
      sh 'minify', '-o', dst_script_path, src_script_path
    end
  end

  private

  #
  # fetch URL, check hash, write to destination path.
  # 
  def fetch(dst_path, src_url, src_hash)
    return if File.exists?(dst_path)
  
    # generate temporary name in destination directory
    # FIXME: probably not necessary now that entire dir is temporary
    tmp_path = File.join(
      File.expand_path(File.dirname(dst_path)),
      '_' + OpenSSL::Random.random_bytes(32).unpack('H*').first
    )
    @log.debug('fetch') { "tmp_path = #{tmp_path}" }
  
    # download to temp path
    URI.parse(src_url).open do |src|
      @log.debug('fetch') { "copying #{src_url} to #{tmp_path}" }
      IO.copy_stream(src, tmp_path)
    end
  
    # check hash
    got_hash = OpenSSL::Digest.file(tmp_path, 'sha256').hexdigest
    if got_hash != src_hash
      File.unlink(tmp_path)
      raise "#{tmp_path}: hash mismatch: (got #{got_hash}, exp #{src_hash}"
    end
  
    # rename to destination path
    @log.debug('fetch') { "rename #{tmp_path} to #{dst_path}" }
    File.rename(tmp_path, dst_path)
  end

  def sh(*cmd)
    # log command
    @log.debug('sh') { JSON(cmd) }

    # exec command
    system(*cmd)
  end
end

# allow cli invocation
Builder.run if __FILE__ == $0
