import { notFound } from 'next/navigation';
import Image from 'next/image';
import { Header } from '@/components/blog/header';
import { VerifiedBadge } from '@/components/blog/verified-badge';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

async function fetchTeam(slug: string) {
  try {
    const res = await fetch(`${API_BASE}/api/teams/slug/${encodeURIComponent(slug)}`, {
      cache: 'no-store',
    });
    if (!res.ok) return null;
    return res.json();
  } catch {
    return null;
  }
}

interface TeamPageProps {
  params: Promise<{ slug: string }>;
}

export default async function TeamPage({ params }: TeamPageProps) {
  const { slug } = await params;
  const team = await fetchTeam(slug);

  if (!team) {
    notFound();
  }

  return (
    <div className="min-h-screen bg-background">
      <Header />
      <main className="mx-auto max-w-2xl px-4 py-12">
        <div className="flex items-start gap-6">
          {team.avatar && (
            <div className="relative h-20 w-20 shrink-0 overflow-hidden rounded-full bg-muted">
              <Image src={team.avatar} alt={team.name} fill className="object-cover" />
            </div>
          )}
          <div>
            <h1 className="flex items-center gap-2 font-serif text-3xl font-bold">
              {team.name}
              <VerifiedBadge verified={!!team.is_verified} className="h-6 w-6 text-sky-500" />
            </h1>
            <p className="text-muted-foreground">@{team.slug}</p>
            {team.bio && <p className="mt-4 text-foreground">{team.bio}</p>}
          </div>
        </div>
      </main>
    </div>
  );
}
