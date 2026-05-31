import Image from "next/image";
import Link from "next/link";
import { Article } from "@/components/blog/article-card";
import { Button } from "@/components/ui/button";

interface AuthorCardProps {
  article: Article;
}

const authorBios: Record<string, { bio: string; followers: string }> = {
  "Sarah Chen": {
    bio: "Principal Engineer with 15 years of experience building scalable systems. I write about software craftsmanship, architecture, and the human side of tech.",
    followers: "12.4K",
  },
  "Marcus Williams": {
    bio: "Former Big Tech engineer turned full-time writer. I help people escape the corporate ladder and build lives they love.",
    followers: "28.1K",
  },
  "Dr. Emma Rodriguez": {
    bio: "Neuroscientist and productivity researcher. Translating brain science into practical advice for ambitious humans.",
    followers: "45.2K",
  },
  "Alex Thompson": {
    bio: "ML Engineer at a stealth startup. Making AI concepts accessible through visual explanations.",
    followers: "8.7K",
  },
  "Dr. James Liu": {
    bio: "AI researcher and futurist. Exploring the intersection of human creativity and machine intelligence.",
    followers: "67.3K",
  },
};

export function AuthorCard({ article }: AuthorCardProps) {
  const authorInfo = authorBios[article.author.name] || {
    bio: "Writer and thinker exploring ideas that matter.",
    followers: "5.2K",
  };

  return (
    <div className="mt-12 border-t border-border pt-8">
      <div className="flex gap-4">
        <Link href="#" className="shrink-0">
          <div className="relative h-16 w-16 overflow-hidden rounded-full bg-muted">
            <Image
              src={article.author.avatar}
              alt={article.author.name}
              fill
              className="object-cover"
            />
          </div>
        </Link>
        <div className="flex-1">
          <div className="mb-1 flex items-center justify-between">
            <div>
              <p className="text-xs font-medium uppercase tracking-wide text-muted-foreground">
                Written by
              </p>
              <Link 
                href="#" 
                className="text-lg font-bold text-foreground hover:underline"
              >
                {article.author.name}
              </Link>
            </div>
            <Button 
              className="bg-primary text-primary-foreground hover:bg-primary/90"
              size="sm"
            >
              Follow
            </Button>
          </div>
          <p className="mb-2 text-sm text-muted-foreground">
            {authorInfo.followers} Followers
          </p>
          <p className="text-foreground">
            {authorInfo.bio}
          </p>
        </div>
      </div>
    </div>
  );
}
