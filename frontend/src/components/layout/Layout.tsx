import { ReactNode, useState, useEffect } from 'react';
import { useNavigate, useLocation } from '@tanstack/react-router';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { MobileNav } from './MobileNav';
import { PlayerBar } from '../player/PlayerBar';
import { SimulationBanner } from '../admin/SimulationBanner';
import { useFeatureFlags } from '../../hooks/useFeatureFlags';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const [mobileNavOpen, setMobileNavOpen] = useState(false);
  const { role, isSimulating, isLoaded } = useFeatureFlags();
  const navigate = useNavigate();
  const location = useLocation();

  // When simulating guest role, redirect to permission-denied page
  useEffect(() => {
    if (isLoaded && isSimulating && role === 'guest' && location.pathname !== '/permission-denied') {
      navigate({ to: '/permission-denied' });
    }
  }, [isLoaded, isSimulating, role, location.pathname, navigate]);

  return (
    <div className="min-h-screen bg-base-200 flex flex-col">
      <Header onMenuClick={() => setMobileNavOpen(true)} />
      <SimulationBanner />
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
