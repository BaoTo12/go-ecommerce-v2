import type { Metadata } from 'next';
import './globals.css';
import Navigation from '../components/Navigation';

export const metadata: Metadata = {
  title: 'Shopee Clone - E-Commerce Platform',
  description: 'Shopee-style e-commerce platform with Flash Sale, Live Shopping, Gamification, and more.',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="vi">
      <body className="min-h-screen bg-[#F5F5F5]">
        <Navigation />
        <main>{children}</main>
        <footer className="bg-[#FBFBFB] border-t py-8 mt-8">
          <div className="container mx-auto px-4">
            <div className="grid gap-8 md:grid-cols-4 text-sm">
              <div>
                <h3 className="font-bold text-[#EE4D2D] mb-4 uppercase">Chăm sóc khách hàng</h3>
                <ul className="space-y-2 text-gray-600">
                  <li><a href="#" className="hover:text-[#EE4D2D]">Trung tâm trợ giúp</a></li>
                  <li><a href="#" className="hover:text-[#EE4D2D]">Hướng dẫn mua hàng</a></li>
                  <li><a href="#" className="hover:text-[#EE4D2D]">Thanh toán</a></li>
                  <li><a href="#" className="hover:text-[#EE4D2D]">Vận chuyển</a></li>
                </ul>
              </div>
              <div>
                <h3 className="font-bold text-[#EE4D2D] mb-4 uppercase">Về Shopee</h3>
                <ul className="space-y-2 text-gray-600">
                  <li><a href="#" className="hover:text-[#EE4D2D]">Giới thiệu</a></li>
                  <li><a href="#" className="hover:text-[#EE4D2D]">Tuyển dụng</a></li>
                  <li><a href="#" className="hover:text-[#EE4D2D]">Điều khoản</a></li>
                  <li><a href="#" className="hover:text-[#EE4D2D]">Chính sách bảo mật</a></li>
                </ul>
              </div>
              <div>
                <h3 className="font-bold text-[#EE4D2D] mb-4 uppercase">Thanh toán</h3>
                <div className="flex flex-wrap gap-2">
                  {['💳 Visa', '💳 Master', '🏧 ATM', '💵 COD'].map(p => (
                    <span key={p} className="bg-white px-2 py-1 rounded border text-xs">{p}</span>
                  ))}
                </div>
                <h3 className="font-bold text-[#EE4D2D] mb-4 mt-6 uppercase">Vận chuyển</h3>
                <div className="flex flex-wrap gap-2">
                  {['🚚 Express', '📦 Standard', '✈️ Fast'].map(s => (
                    <span key={s} className="bg-white px-2 py-1 rounded border text-xs">{s}</span>
                  ))}
                </div>
              </div>
              <div>
                <h3 className="font-bold text-[#EE4D2D] mb-4 uppercase">Theo dõi</h3>
                <div className="flex gap-3 text-2xl">
                  <a href="#" className="hover:opacity-70">📘</a>
                  <a href="#" className="hover:opacity-70">📸</a>
                  <a href="#" className="hover:opacity-70">🐦</a>
                </div>
                <h3 className="font-bold text-[#EE4D2D] mb-4 mt-6 uppercase">Tải ứng dụng</h3>
                <div className="flex gap-2">
                  <span className="bg-black text-white px-3 py-1 rounded text-xs">📱 App Store</span>
                  <span className="bg-black text-white px-3 py-1 rounded text-xs">🤖 Play Store</span>
                </div>
              </div>
            </div>
            <div className="mt-8 pt-6 border-t text-center text-gray-500 text-xs">
              © 2024 Shopee Clone. Hyperscale E-Commerce Platform.
            </div>
          </div>
        </footer>
      </body>
    </html>
  );
}
