# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Sim < Formula
  desc "Straight-forward, fast, scalable API simulation."
  homepage "https://github.com/kitproj/sim"
  version "0.0.19"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/sim/releases/download/v0.0.19/sim_0.0.19_Darwin_x86_64.tar.gz"
      sha256 "971098973efd1c0a8c884e2d912e8aad7280cdf59ebbcaa5127bee61ab3e6002"

      def install
        bin.install "sim"
      end
    end
    if Hardware::CPU.arm?
      url "https://github.com/kitproj/sim/releases/download/v0.0.19/sim_0.0.19_Darwin_arm64.tar.gz"
      sha256 "0391e62be1f3b1b5bea878fd7297890b715bdf034e0fee291d4322f841b2878c"

      def install
        bin.install "sim"
      end
    end
  end

  on_linux do
    if Hardware::CPU.arm? && Hardware::CPU.is_64_bit?
      url "https://github.com/kitproj/sim/releases/download/v0.0.19/sim_0.0.19_Linux_arm64.tar.gz"
      sha256 "2ec028b9d860bb73e94a53563943866ecfe63c0b60f164f407771333de1679e9"

      def install
        bin.install "sim"
      end
    end
    if Hardware::CPU.intel?
      url "https://github.com/kitproj/sim/releases/download/v0.0.19/sim_0.0.19_Linux_x86_64.tar.gz"
      sha256 "99d384048df914f75d650f71e6ecbd7d07138e6439481c8d1ed4dea93c21a79a"

      def install
        bin.install "sim"
      end
    end
  end
end
