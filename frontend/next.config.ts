import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  compiler: {
    styledComponents: true,
  },
  images: {
    remotePatterns: [
      {
        protocol: 'http',
        hostname: 'localhost',
        port: '',
        pathname: '/**',
        search: '',
      },
    ],
    unoptimized: true,
  },
};


module.exports = {
  async rewrites() {
    return [
    ]
  },
  output: 'standalone'
}

export default nextConfig;
