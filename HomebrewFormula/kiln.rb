# typed: false
# frozen_string_literal: true

# This file was generated by GoReleaser. DO NOT EDIT.
class Kiln < Formula
  desc ""
  homepage ""
  version "0.63.0-rc.1"

  on_macos do
    if Hardware::CPU.intel?
      url "https://github.com/pivotal-cf/kiln/releases/download/0.63.0-rc.1/kiln-darwin-0.63.0-rc.1.tar.gz"
      sha256 "b610e54ad9426c93b802f351a3d5895a7bef7d0a71fdae1d944cd1e1353ff175"

      def install
        bin.install "kiln"
      end
    end
  end

  on_linux do
    if Hardware::CPU.intel?
      url "https://github.com/pivotal-cf/kiln/releases/download/0.63.0-rc.1/kiln-linux-0.63.0-rc.1.tar.gz"
      sha256 "116ed0b303fd6cdf6c5d327d805422c1a27b0ae839a46a2bb0e48e7cba7f3b5b"

      def install
        bin.install "kiln"
      end
    end
  end

  test do
    system "#{bin}/kiln --version"
  end
end
