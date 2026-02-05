import { Link } from '@tanstack/react-router';

export function HelloNav() {
  return (
    <nav>
      <Link to="/hello-search" className="btn btn-ghost">
        Hello Search
      </Link>
    </nav>
  );
}
