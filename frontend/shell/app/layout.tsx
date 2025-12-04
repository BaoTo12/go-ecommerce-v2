import type { Metadata } from 'next';
import './globals.css';
import Navigation from '../components/Navigation';

export const metadata: Metadata = {
  title: 'Titan Commerce - Hyperscale E-Commerce Platform',
  description: 'Next-generation e-commerce platform with 50M DAU capacity, featuring live shopping, flash sales, gamification, and AI-powered features.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="min-h-screen bg-gray-50">
        <Navigation />
        <main>{children}</main>
        <footer className="border-t bg-white py-8 mt-16">
          <div className="container mx-auto px-4">
            <div className="grid gap-8 md:grid-cols-4">
              <div>
                <h3 className="mb-4 font-bold">🚀 Titan Commerce</h3>
                <p className="text-sm text-muted-foreground">
                  Hyperscale e-commerce platform serving 50M+ daily active users across Southeast Asia.
                </p>
              </div>
              <div>
                <h4 className="mb-3 font-semibold">Shopping</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li><a href="/products" className="hover:text-foreground">All Products</a></li>
                  <li><a href="/deals/flash-sale" className="hover:text-foreground">Flash Sales</a></li>
                  <li><a href="/live" className="hover:text-foreground">Live Shopping</a></li>
                  <li><a href="/deals/coupons" className="hover:text-foreground">Coupons</a></li>
                </ul>
              </div>
              <div>
                <h4 className="mb-3 font-semibold">Rewards</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li><a href="/rewards" className="hover:text-foreground">Daily Check-in</a></li>
                  <li><a href="/rewards" className="hover:text-foreground">Lucky Draw</a></li>
                  <li><a href="/rewards" className="hover:text-foreground">Missions</a></li>
                </ul>
              </div>
              <div>
                <h4 className="mb-3 font-semibold">Platform</h4>
                <ul className="space-y-2 text-sm text-muted-foreground">
                  <li>🏗️ Cell-Based Architecture</li>
                  <li>⚡ 50ms P99 Latency</li>
                  <li>🛡️ Fraud Detection</li>
                  <li>📊 Real-time Analytics</li>
                </ul>
              </div>
            </div>
            <div className="mt-8 border-t pt-6 text-center text-sm text-muted-foreground">
              © 2024 Titan Commerce. Built for hyperscale.
            </div>
          </div>
        </footer>
      </body>
    </html>
  );
}
