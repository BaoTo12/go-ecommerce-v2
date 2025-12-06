/** @type {import('next').NextConfig} */
const nextConfig = {
    reactStrictMode: true,
    // Module Federation will be configured here for micro-frontends
    experimental: {
        appDir: true,
    },
    images: {
        remotePatterns: [
            {
                protocol: 'https',
                hostname: 'images.unsplash.com',
            },
            {
                protocol: 'https',
                hostname: 'ui-avatars.com',
            },
            {
                protocol: 'https',
                hostname: 'via.placeholder.com',
            },
        ],
    },
}

module.exports = nextConfig
