# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Sim < Formula
  desc "Straight-forward, fast, scalable API simulation."
  homepage "https://github.com/kitproj/sim"
  version "0.0.4"

  on_macos do
    if Hardware::CPU.arm?
      url "https://github.com/kitproj/sim/releases/download/v0.0.4/sim_0.0.4_Darwin_arm64.tar.gz"
      sha256 "d9a8e23afdfe9c8edb18d9033cc80152281733812260b4575980e081ff479516"

      def install
        bin.install "sim"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/sim/releases/download/v0.0.4/sim_0.0.4_Darwin_x86_64.tar.gz"
      sha256 "30f01c356c26398e3cb8aefc028fbada3934b6689ba825d938286af973ff15d3"

      def install
        bin.install "sim"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/sim/releases/download/v0.0.4/sim_0.0.4_Linux_x86_64.tar.gz"
      sha256 "9a6b487b511aa51af98e71dd37a2103f2af77d8029ad7e938d7e957aeb403f23"

      def install
        bin.install "sim"
      end
    end
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kitproj/sim/releases/download/v0.0.4/sim_0.0.4_Linux_arm64.tar.gz"
      sha256 "e138efe77085afa73d84247ea9d8475adaad715196b9e518cf470745173ad596"

      def install
        bin.install "sim"
      end
    end
  end
end