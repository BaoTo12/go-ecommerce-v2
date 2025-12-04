/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    // Module Federation will be configured here for micro-frontends
    experimental: {
        appDir: true,
    },
}

module.exports = nextConfig
