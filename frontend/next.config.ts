import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  compiler: {
    styledComponents: true
  },
  images: {
   remotePatterns: [
    new URL('http://localhost/**')
   ],
},
};

module.exports = {
  async rewrites() {
    return [
    ]
  }
}

export default nextConfig;
