'use client';

import { useRouter } from 'next/navigation';
import { useState } from 'react';
import Link from 'next/link';
import { PenLine, Search } from 'lucide-react';
import { Button } from '@/components/ui/button';
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar';
import { Bell, Settings, LogOut } from 'lucide-react';
import { useAuth } from '@/lib/auth/auth-context';

export function HeaderSearch() {
  const router = useRouter();
  const [query, setQuery] = useState('');

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (query.trim()) {
      router.push(`/search?q=${encodeURIComponent(query.trim())}`);
    }
  };

  return (
    <form onSubmit={handleSubmit} className="relative hidden md:block">
      <Search className="search-input-icon absolute top-1/2 h-4 w-4 -translate-y-1/2 text-muted-foreground ltr:left-3 rtl:right-3" />
      <input
        type="search"
        value={query}
        onChange={(e) => setQuery(e.target.value)}
        placeholder="Search"
        className="search-input h-9 w-52 rounded-full bg-secondary text-sm text-foreground placeholder:text-muted-foreground focus:outline-none focus:ring-2 focus:ring-primary/20 ltr:pl-9 ltr:pr-4 rtl:pr-9 rtl:pl-4"
      />
    </form>
  );
}

export function HeaderNav() {
  const { user, isAuthenticated, logout } = useAuth();

  return (
    <nav className="flex items-center gap-2">
      {isAuthenticated && (
        <Link
          href="/write"
          className="hidden items-center gap-2 text-sm text-muted-foreground hover:text-foreground sm:flex"
        >
          <PenLine className="h-5 w-5" />
          <span>Write</span>
        </Link>
      )}
      <Button variant="ghost" size="icon" className="text-muted-foreground hover:text-foreground">
        <Bell className="h-5 w-5" />
      </Button>

      {isAuthenticated && user ? (
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="ghost" className="relative h-8 w-8 rounded-full">
              <Avatar className="h-8 w-8">
                <AvatarImage src={user.avatar} alt={user.name} />
                <AvatarFallback>{user.name.charAt(0).toUpperCase()}</AvatarFallback>
              </Avatar>
            </Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent align="end" className="w-48">
            <DropdownMenuItem asChild>
              <Link href="/me/stories">My stories</Link>
            </DropdownMenuItem>
            <DropdownMenuItem asChild>
              <Link href="/teams">Teams</Link>
            </DropdownMenuItem>
            {user.role === 'super_admin' && (
              <DropdownMenuItem asChild>
                <Link href="/admin/verification">Admin</Link>
              </DropdownMenuItem>
            )}
            <DropdownMenuItem asChild>
              <Link href="/settings">Settings</Link>
            </DropdownMenuItem>
            <DropdownMenuSeparator />
            <DropdownMenuItem
              onClick={() => void logout()}
              className="text-destructive"
            >
              <LogOut className="me-2 h-4 w-4" />
              Sign out
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      ) : (
        <Button
          asChild
          size="sm"
          className="rounded-full bg-primary px-4 text-primary-foreground hover:bg-primary/90"
        >
          <Link href="/sign-in">Sign in</Link>
        </Button>
      )}
    </nav>
  );
}
