import { ReactNode, useState } from 'react';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { MobileNav } from './MobileNav';
import { PlayerBar } from '../player/PlayerBar';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const [mobileNavOpen, setMobileNavOpen] = useState(false);

  return (
    <div className="min-h-screen bg-base-200 flex flex-col">
      <Header onMenuClick={() => setMobileNavOpen(true)} />
      <div className="flex flex-1">
        {/* Desktop sidebar */}
        <Sidebar />
        <main className="flex-1 p-6 pb-24">{children}</main>
      </div>
      <PlayerBar />

      {/* Mobile navigation overlay */}
      <MobileNav isOpen={mobileNavOpen} onClose={() => setMobileNavOpen(false)} />
    </div>
  );
}
