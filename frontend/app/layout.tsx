import './globals.css';
import { Inter } from 'next/font/google';

const inter = Inter({ 
  subsets: ['latin', 'cyrillic'],
  weight: ['300', '400', '500', '600', '700', '800', '900'],
});

export default function Layout({ children }: { children: React.ReactNode }) {
  return (
    <html>
      <body className={inter.className}>
        {children}
      </body>
    </html>
  )
}