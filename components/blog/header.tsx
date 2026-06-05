'use client';

import Link from "next/link";
import { useTranslations } from "next-intl";
import { HeaderNav, HeaderSearch } from "@/components/blog/header-client";

export function Header() {
  const tCommon = useTranslations("common");

  return (
    <header className="sticky top-0 z-50 border-b border-border bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/60">
      <div className="mx-auto flex h-14 max-w-5xl items-center justify-between px-4">
        <div className="flex items-center gap-6">
          <Link href="/" className="flex items-center gap-2">
            <span className="font-serif text-2xl font-bold tracking-tight text-foreground">
              {tCommon("brand")}
            </span>
          </Link>
          <HeaderSearch />
        </div>
        <HeaderNav />
      </div>
    </header>
  );
}
