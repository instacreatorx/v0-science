import { BadgeCheck } from "lucide-react";

interface VerifiedBadgeProps {
  verified?: boolean;
  className?: string;
}

export function VerifiedBadge({ verified, className = "h-3.5 w-3.5 text-sky-500" }: VerifiedBadgeProps) {
  if (!verified) return null;
  return <BadgeCheck className={className} aria-label="Verified" />;
}
