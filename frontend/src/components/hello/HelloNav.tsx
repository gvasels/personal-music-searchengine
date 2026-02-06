/**
 * HelloNav Component
 *
 * Navigation bar for the Hello Music Search page.
 * Provides branding and a link back to the home page.
 */

export function HelloNav() {
  return (
    <nav className="navbar bg-base-100">
      <div className="flex-1">
        <span className="text-xl font-bold">Hello Music Search</span>
      </div>
      <div className="flex-none">
        <a href="/" className="btn btn-ghost">
          Home
        </a>
      </div>
    </nav>
  );
}
