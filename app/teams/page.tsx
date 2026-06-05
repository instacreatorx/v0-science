'use client';

import { useEffect, useState } from 'react';
import Link from 'next/link';
import { toast } from 'sonner';
import { Header } from '@/components/blog/header';
import { ProtectedRoute } from '@/components/auth/protected-route';
import { VerifiedBadge } from '@/components/blog/verified-badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Textarea } from '@/components/ui/textarea';
import { teamsApi, type Team } from '@/lib/api';

export default function TeamsPage() {
  const [teams, setTeams] = useState<Team[]>([]);
  const [name, setName] = useState('');
  const [slug, setSlug] = useState('');
  const [bio, setBio] = useState('');
  const [creating, setCreating] = useState(false);

  const load = async () => {
    try {
      const data = await teamsApi.mine();
      setTeams(data);
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to load teams');
    }
  };

  useEffect(() => {
    void load();
  }, []);

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setCreating(true);
    try {
      await teamsApi.create({ name, slug, bio });
      setName('');
      setSlug('');
      setBio('');
      toast.success('Team created');
      void load();
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to create team');
    } finally {
      setCreating(false);
    }
  };

  const handleVerifyRequest = async (teamId: number) => {
    const proofText = window.prompt('Why should this team be verified?');
    if (!proofText?.trim()) return;
    try {
      await teamsApi.submitVerification(teamId, { proof_text: proofText.trim() });
      toast.success('Verification request submitted');
    } catch (error) {
      toast.error(error instanceof Error ? error.message : 'Failed to submit request');
    }
  };

  return (
    <ProtectedRoute>
      <div className="min-h-screen bg-background">
        <Header />
        <main className="mx-auto max-w-2xl px-4 py-8">
          <h1 className="mb-8 font-serif text-3xl font-bold">Your teams</h1>

          <form onSubmit={handleCreate} className="mb-10 space-y-4 rounded-lg border border-border p-6">
            <h2 className="font-semibold">Create a team publication</h2>
            <div className="space-y-2">
              <Label htmlFor="name">Name</Label>
              <Input id="name" value={name} onChange={(e) => setName(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="slug">Slug</Label>
              <Input id="slug" value={slug} onChange={(e) => setSlug(e.target.value)} required />
            </div>
            <div className="space-y-2">
              <Label htmlFor="bio">Bio</Label>
              <Textarea id="bio" value={bio} onChange={(e) => setBio(e.target.value)} rows={3} />
            </div>
            <Button type="submit" disabled={creating}>
              {creating ? 'Creating...' : 'Create team'}
            </Button>
          </form>

          <ul className="space-y-4">
            {teams.map((team) => (
              <li key={team.id} className="flex items-center justify-between rounded-lg border border-border p-4">
                <div>
                  <Link href={`/team/${team.slug}`} className="flex items-center gap-2 font-semibold hover:underline">
                    {team.name}
                    <VerifiedBadge verified={!!team.is_verified} />
                  </Link>
                  <p className="text-sm text-muted-foreground">/{team.slug}</p>
                </div>
                {!team.is_verified && (
                  <Button size="sm" variant="outline" onClick={() => void handleVerifyRequest(team.id)}>
                    Request verification
                  </Button>
                )}
              </li>
            ))}
          </ul>
        </main>
      </div>
    </ProtectedRoute>
  );
}
