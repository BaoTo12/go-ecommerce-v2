import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  title: 'Titan Commerce - Hyperscale E-Commerce Platform',
  description: 'Modern e-commerce platform with live streaming, gamification, and social shopping',
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <header className="sticky top-0 z-50 w-full border-b bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
          <div className="container flex h-16 items-center">
            <div className="mr-4 flex">
              <a className="mr-6 flex items-center space-x-2" href="/">
                <span className="hidden font-bold sm:inline-block">
                  Titan Commerce
                </span>
              </a>
            </div>
            <nav className="flex items-center space-x-6 text-sm font-medium">
              <a href="/products">Products</a>
              <a href="/live">Live Streaming</a>
              <a href="/deals">Flash Sales</a>
              <a href="/rewards">Rewards</a>
            </nav>
            <div className="flex flex-1 items-center justify-end space-x-4">
              <a href="/cart">Cart</a>
              <a href="/login">Login</a>
            </div>
          </div>
        </header>
        <main>{children}</main>
        <footer className="border-t">
          <div className="container py-8">
            <p className="text-center text-sm text-muted-foreground">
              © 2025 Titan Commerce. Built with Next.js 15, React 19, and Tailwind CSS.
            </p>
          </div>
        </footer>
      </body>
    </html>
  )
}
