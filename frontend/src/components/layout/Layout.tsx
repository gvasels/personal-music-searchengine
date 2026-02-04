import { ReactNode, useState, useEffect, useRef } from 'react';
import { useNavigate, useLocation } from '@tanstack/react-router';
import { Header } from './Header';
import { Sidebar } from './Sidebar';
import { MobileNav } from './MobileNav';
import { PlayerBar } from '../player/PlayerBar';
import { SimulationBanner } from '../admin/SimulationBanner';
import { useFeatureFlags } from '../../hooks/useFeatureFlags';
import { useAuth } from '../../hooks/useAuth';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  const [mobileNavOpen, setMobileNavOpen] = useState(false);
  const { role, isSimulating, isLoaded } = useFeatureFlags();
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  // Track previous simulation state to detect changes
  const prevSimulatingGuestRef = useRef<boolean | null>(null);

  // Handle role simulation navigation
  useEffect(() => {
    if (!isLoaded) return;

    const isSimulatingGuest = isSimulating && role === 'guest';
    const onPermissionDenied = location.pathname === '/permission-denied';
    const wasSimulatingGuest = prevSimulatingGuestRef.current;

    // When simulating guest role, redirect to permission-denied page
    if (isSimulatingGuest && !onPermissionDenied) {
      navigate({ to: '/permission-denied' });
    }
    // When stopped simulating guest (switched role or stopped simulation),
    // and we're still on permission-denied, redirect to home
    else if (wasSimulatingGuest === true && !isSimulatingGuest && onPermissionDenied && isAuthenticated) {
      navigate({ to: '/' });
    }

    // Update ref for next comparison
    prevSimulatingGuestRef.current = isSimulatingGuest;
  }, [isLoaded, isSimulating, role, location.pathname, navigate, isAuthenticated]);

  return (
    <div className="h-screen bg-base-200 flex flex-col overflow-hidden">
      <Header onMenuClick={() => setMobileNavOpen(true)} />
      <SimulationBanner />
      <div className="flex flex-1 min-h-0">
        {/* Desktop sidebar */}
        <Sidebar />
        <main className="flex-1 p-6 pb-28 overflow-auto">{children}</main>
      </div>
      <PlayerBar />

      {/* Mobile navigation overlay */}
      <MobileNav isOpen={mobileNavOpen} onClose={() => setMobileNavOpen(false)} />
    </div>
  );
}
