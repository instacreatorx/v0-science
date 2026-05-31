import Image from "next/image";
import Link from "next/link";
import { Button } from "@/components/ui/button";

const topics = [
  "Technology",
  "Self Improvement",
  "Writing",
  "Relationships",
  "Machine Learning",
  "Productivity",
  "Politics",
];

const recommendedAuthors = [
  {
    id: "1",
    name: "Sarah Chen",
    bio: "Tech writer & entrepreneur",
    avatar: "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=100&h=100&fit=crop&crop=face",
    followers: "42K",
  },
  {
    id: "2",
    name: "Marcus Williams",
    bio: "Exploring the future of AI",
    avatar: "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face",
    followers: "28K",
  },
  {
    id: "3",
    name: "Emma Rodriguez",
    bio: "Designer & Creative Director",
    avatar: "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100&h=100&fit=crop&crop=face",
    followers: "35K",
  },
];

export function Sidebar() {
  return (
    <aside className="hidden lg:block lg:w-80 lg:shrink-0">
      <div className="sticky top-20 space-y-8">
        {/* Membership CTA */}
        <div className="rounded-lg border border-border bg-card p-5">
          <h3 className="mb-2 font-serif text-lg font-bold text-card-foreground">
            Get unlimited access
          </h3>
          <p className="mb-4 text-sm text-muted-foreground">
            Read member-only stories and support writers you love.
          </p>
          <Button className="w-full rounded-full bg-primary text-primary-foreground hover:bg-primary/90">
            Become a member
          </Button>
        </div>

        {/* Recommended Topics */}
        <div>
          <h3 className="mb-4 text-sm font-bold uppercase tracking-wider text-foreground">
            Recommended Topics
          </h3>
          <div className="flex flex-wrap gap-2">
            {topics.map((topic) => (
              <Link
                key={topic}
                href={`/topic/${topic.toLowerCase().replace(" ", "-")}`}
                className="rounded-full bg-secondary px-4 py-2 text-sm text-secondary-foreground transition-colors hover:bg-muted"
              >
                {topic}
              </Link>
            ))}
          </div>
        </div>

        {/* Who to Follow */}
        <div>
          <h3 className="mb-4 text-sm font-bold uppercase tracking-wider text-foreground">
            Who to Follow
          </h3>
          <div className="space-y-4">
            {recommendedAuthors.map((author) => (
              <div key={author.id} className="flex items-start gap-3">
                <div className="relative h-10 w-10 shrink-0 overflow-hidden rounded-full bg-muted">
                  <Image
                    src={author.avatar}
                    alt={author.name}
                    fill
                    className="object-cover"
                  />
                </div>
                <div className="min-w-0 flex-1">
                  <Link
                    href={`/author/${author.id}`}
                    className="block truncate font-medium text-foreground hover:underline"
                  >
                    {author.name}
                  </Link>
                  <p className="truncate text-sm text-muted-foreground">{author.bio}</p>
                </div>
                <Button
                  variant="outline"
                  size="sm"
                  className="shrink-0 rounded-full text-xs"
                >
                  Follow
                </Button>
              </div>
            ))}
          </div>
          <Link
            href="/explore/authors"
            className="mt-4 block text-sm font-medium text-primary hover:underline"
          >
            See more suggestions
          </Link>
        </div>

        {/* Footer Links */}
        <div className="flex flex-wrap gap-x-4 gap-y-1 text-xs text-muted-foreground">
          <Link href="/help" className="hover:text-foreground">
            Help
          </Link>
          <Link href="/status" className="hover:text-foreground">
            Status
          </Link>
          <Link href="/about" className="hover:text-foreground">
            About
          </Link>
          <Link href="/careers" className="hover:text-foreground">
            Careers
          </Link>
          <Link href="/blog" className="hover:text-foreground">
            Blog
          </Link>
          <Link href="/privacy" className="hover:text-foreground">
            Privacy
          </Link>
          <Link href="/terms" className="hover:text-foreground">
            Terms
          </Link>
        </div>
      </div>
    </aside>
  );
}
