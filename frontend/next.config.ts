import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  compiler: {
    styledComponents: true
  },
  images: {
   remotePatterns: [
   {
      protocol: "http",
      hostname: "**",
    },
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
