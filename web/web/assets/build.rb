#!/usr/bin/env ruby
# frozen_string_literal: true

#
# build.rb: Build packed web assets.
#
# This script generates the following packed files in the `public/`
# directory:
#
# - public/index.html
# - public/script.min.js
# - public/style.min.css
#
# This script uses the following command-line tools:
#
# - esbuild (https://esbuild.github.io/)
# - minify (https://github.com/tdewolf/minify)
# - sass (https://sass-lang.com/)
#

# load libraries
require 'fileutils'
require 'json'
require 'logger'
require 'openssl'
require 'open-uri'
require 'uri'
require 'zip'

#
# Configuration object.
#
# Reads configuration from environment variables.
#
class Config
  attr :log_level, :out_dir, :zip_url, :zip_hash, :esbuild_cmd, :minify_cmd, :sass_cmd

  # bulma zip URL
  BULMA_ZIP_URL = 'https://github.com/jgthms/bulma/releases/download/0.9.4/bulma-0.9.4.zip'

  # expected sha256 hash of bulma zip
  BULMA_ZIP_HASH = '781fb4d195c83fb37090a4b1b66f2a8fc72ea993dc708e6e9bfdce30088fa635'

  # default output directory
  DEFAULT_OUTPUT_DIR = File.join(__dir__, '..', 'public')

  #
  # Create Config instance from environment variables.
  #
  # Reads from the following environment variables:
  #
  # * `BUILD_LOG_LEVEL`: log level.  Defaults to `info` if unspecified.
  # * `BUILD_OUTPUT_DIR`: Path to output `public/` directory.
  # * `BUILD_BULMA_ZIP_URL`: URL to `bulma-0.9.4.zip`.
  # * `BUILD_BULMA_ZIP_HASH`: Hex-encoded SHA256 hash of bulma zip file.
  # * `BUILD_ESBUILD_CMD`: `esbuild` command.
  # * `BUILD_MINIFY_CMD`: `minify` command.
  # * `BUILD_SASS_CMD`: `sass` command.
  #
  def initialize
    @log_level = ENV.fetch('BUILD_LOG_LEVEL', 'info')
    @out_dir = ENV.fetch('BUILD_OUTPUT_DIR', DEFAULT_OUTPUT_DIR)
    @zip_url = ENV.fetch('BUILD_BULMA_ZIP_URL', BULMA_ZIP_URL)
    @zip_hash = ENV.fetch('BUILD_BULMA_ZIP_HASH', BULMA_ZIP_HASH)

    # get commands from environment
    @esbuild_cmd = ENV.fetch('BUILD_ESBUILD_CMD', 'esbuild')
    @minify_cmd = ENV.fetch('BUILD_MINIFY_CMD', 'minify')
    @sass_cmd = ENV.fetch('BUILD_SASS_CMD', 'sass')
  end
end

#
# Class which generates the following files in the `public/` directory
# based on configuration stored in environment variables:
#
# - `public/index.html`
# - `public/script.min.js`
# - `public/style.min.css`
#
class Builder
  #
  # Singleton method to create and run builder from environment.
  #
  def self.run
    new(Config.new).run
  end

  #
  # Create Builder instance.
  #
  def initialize(config)
    @config = config.freeze

    # init logger
    @log = ::Logger.new(STDERR)
    @log.level = config.log_level
  end

  #
  # Run builder and generate the following files in the `public/`
  # directory:
  #
  # - `public/index.html`
  # - `public/script.min.js`
  # - `public/style.min.css`
  #
  def run
    dur = time do
      # work in temporary directory
      Dir.mktmpdir do |tmp_dir|
        # switch to temp dir
        Dir.chdir(tmp_dir)

        # start/join build threads
        [
          Thread.new(tmp_dir) { |tmp_dir| build_style(tmp_dir) },
          Thread.new { build_script },
          Thread.new { build_index },
        ].each { |t| t.join }
      end
    end

    @log.debug(__method__) { 'done (%f seconds)' % [dur] }
  end

  private

  # path to unpacked sass file
  SRC_SASS_PATH = File.join(__dir__, 'bookman.sass')

  #
  # Build `public/style.min.css` by doing the following:
  #
  # 1. Download bulma-0.9.4.zip.
  # 2. Verify the SHA256 hash of the zip file.
  # 3. Extract `bulma/sass` from the zip file.
  # 4. Copy `assets/bookman.sass` to `bulma/`.
  # 5. Run `sass` on `bulma/bookman.sass`.
  # 6. Pipe the `sass` output to `esbuild`.
  # 7. Write the `esbuild` output to `public/style.min.css`.
  #
  def build_style(tmp_dir)
    log_build('public/style.min.css') do
      # download bulma-0.9.4.zip
      zip_path = File.join(tmp_dir, File.basename(@config.zip_url))
      @log.debug(__method__) { "zip_path = #{zip_path}" }
      fetch(zip_path, @config.zip_url, @config.zip_hash)

      # create css working directory
      css_dir = File.join(tmp_dir, 'css')
      @log.debug(__method__) { "creating css dir #{css_dir}" }
      Dir.mkdir(css_dir)

      # expand bulma.zip
      @log.debug(__method__) { "extracting sass files" }
      File.open(zip_path, 'rb') do |fh|
        Zip::File.open(fh) do |zip|
          # find sass entries
          entries = zip.entries.select { |e| e.name =~ /bulma\/sass/ }
          @log.debug(__method__) { JSON(entries.map { |e| e.name }) }

          # expand entries
          entries.each do |e|
            # check for path traversal
            raise "invalid path: #{e.name}" if e.name.match(/\.\./)

            if e.directory?
              # create directory
              dir_path = File.join(css_dir, e.name)
              @log.debug(__method__) { "creating dir #{dir_path}" }
              FileUtils.mkdir_p(dir_path)
            else
              # uncompress file
              dst_path = File.join(css_dir, e.name)
              @log.debug(__method__) { "extracting #{e.name} to #{dst_path}" }
              IO.copy_stream(e.get_input_stream, dst_path)
            end
          end
        end
      end
      @log.debug(__method__) { "done extracting sass files" }

      # build sass destination path
      dst_sass_path = File.join(css_dir, 'bulma', 'bookman.sass')

      # copy bookman.sass into css dir
      @log.debug(__method__) { "copy #{SRC_SASS_PATH} to #{dst_sass_path}" }
      IO.copy_stream(SRC_SASS_PATH, dst_sass_path)

      # build style.min.css path and commands
      dst_style_path = File.join(@config.out_dir, 'style.min.css')
      src_cmd = [@config.sass_cmd, '--sourcemap=none', dst_sass_path]
      dst_cmd = [@config.esbuild_cmd, '--minify', '--loader=css']
      @log.debug(__method__) { JSON({ src_cmd: src_cmd, dst_cmd: dst_cmd }) }

      # generate public/style.min.css by doing the following:
      #
      # 1. run sass on input bookman.sass and bulma directory.
      # 2. pipe the sass output into esbuild.
      # 3. pipe the esbuild output into `public/style.min.css`.
      IO.popen(src_cmd, 'r') do |src_io|
        IO.popen(dst_cmd, 'w', out: dst_style_path) do |dst_io|
          IO.copy_stream(src_io, dst_io)
        end
      end
    end
  end

  # path to unpacked script.js
  SRC_SCRIPT_PATH = File.join(__dir__, 'script.js')

  #
  # Build `public/script.min.js`.
  #
  def build_script
    log_build('public/script.min.js') do
      # build destination path and command
      dst_path = File.join(@config.out_dir, 'script.min.js')
      dst_cmd = [@config.esbuild_cmd, '--minify']

      # run command
      IO.popen(dst_cmd, 'w', out: dst_path) do |dst_io|
        IO.copy_stream(SRC_SCRIPT_PATH, dst_io)
      end
    end
  end

  # path to unpacked index.html
  SRC_INDEX_PATH = File.join(__dir__, '../assets/index.html')

  #
  # Build `public/index.html`.
  #
  def build_index
    log_build('public/index.html') do
      # build destination path
      dst_path = File.join(@config.out_dir, 'index.html')

      # minify index.html
      sh @config.minify_cmd, '-o', dst_path, SRC_INDEX_PATH
    end
  end

  def log_build(name, &block)
    @log.debug(__method__) { "generating #{name}" }
    dur = time(&block)
    @log.debug(__method__) do
      "done generating %s (%f seconds)" % [name, dur]
    end
  end

  #
  # fetch URL, check hash, write to destination path.
  #
  def fetch(dst_path, src_url, src_hash)
    return if File.exists?(dst_path)

    dur = time do
      # generate temporary name in destination directory
      # FIXME: probably not necessary now that entire dir is temporary
      tmp_path = File.join(
        File.expand_path(File.dirname(dst_path)),
        '_' + OpenSSL::Random.random_bytes(32).unpack('H*').first
      )
      @log.debug(__method__) { "tmp_path = #{tmp_path}" }

      # download to temp path
      URI.parse(src_url).open do |src|
        @log.debug(__method__) { "copying #{src_url} to #{tmp_path}" }
        IO.copy_stream(src, tmp_path)
      end

      # check hash
      got_hash = OpenSSL::Digest.file(tmp_path, 'sha256').hexdigest
      if got_hash != src_hash
        File.unlink(tmp_path)
        raise "#{tmp_path}: hash mismatch: (got #{got_hash}, exp #{src_hash}"
      end

      # rename to destination path
      @log.debug(__method__) { "rename #{tmp_path} to #{dst_path}" }
      File.rename(tmp_path, dst_path)
    end

    @log.debug(__method__) { 'done (%f seconds)' % [dur] }
  end

  #
  # Return duration of given block, in seconds.
  #
  def time
    t0 = Time.now
    yield
    Time.now - t0
  end

  #
  # Execute command.
  #
  def sh(*cmd)
    # log command
    @log.debug('sh') { JSON(cmd) }

    # exec command
    system(*cmd)
  end
end

# allow CLI invocation
Builder.run if __FILE__ == $0
