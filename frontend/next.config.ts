import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  compiler: {
    styledComponents: true
  }
};

module.exports = {
  async rewrites() {
    return [
    ]
  }
}

export default nextConfig;
