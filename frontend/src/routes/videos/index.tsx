import { useEffect } from 'react';

export default function VideosPage() {
  useEffect(() => {
    document.title = 'Videos - Music Search Engine';
  }, []);

  return (
    <main className="min-h-screen bg-base-200 p-4 md:p-8">
      <div className="max-w-6xl mx-auto">
        <h1 className="text-3xl font-bold mb-6">Videos</h1>
        <div className="text-center py-16 bg-base-100 rounded-lg">
          <span className="text-5xl mb-4 block">🎬</span>
          <p className="text-base-content/70">Video catalog coming soon</p>
        </div>
      </div>
    </main>
  );
}
