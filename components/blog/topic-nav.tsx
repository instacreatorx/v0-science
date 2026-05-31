"use client";

import Link from "next/link";
import { useRef, useState } from "react";
import { ChevronLeft, ChevronRight, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const topics = [
  { id: "for-you", label: "For You", href: "/" },
  { id: "following", label: "Following", href: "/following" },
  { id: "featured", label: "Featured", href: "/featured" },
  { id: "technology", label: "Technology", href: "/topic/technology" },
  { id: "programming", label: "Programming", href: "/topic/programming" },
  { id: "data-science", label: "Data Science", href: "/topic/data-science" },
  { id: "self-improvement", label: "Self Improvement", href: "/topic/self-improvement" },
  { id: "writing", label: "Writing", href: "/topic/writing" },
  { id: "productivity", label: "Productivity", href: "/topic/productivity" },
  { id: "design", label: "Design", href: "/topic/design" },
];

export function TopicNav() {
  const scrollRef = useRef<HTMLDivElement>(null);
  const [activeId, setActiveId] = useState("for-you");
  const [showLeftArrow, setShowLeftArrow] = useState(false);
  const [showRightArrow, setShowRightArrow] = useState(true);

  const scroll = (direction: "left" | "right") => {
    if (scrollRef.current) {
      const scrollAmount = 200;
      scrollRef.current.scrollBy({
        left: direction === "left" ? -scrollAmount : scrollAmount,
        behavior: "smooth",
      });
    }
  };

  const handleScroll = () => {
    if (scrollRef.current) {
      const { scrollLeft, scrollWidth, clientWidth } = scrollRef.current;
      setShowLeftArrow(scrollLeft > 0);
      setShowRightArrow(scrollLeft < scrollWidth - clientWidth - 10);
    }
  };

  return (
    <div className="relative border-b border-border bg-background">
      <div className="mx-auto flex max-w-5xl items-center px-4">
        <Button
          variant="ghost"
          size="icon"
          className="mr-2 h-8 w-8 shrink-0 text-muted-foreground hover:text-foreground"
        >
          <Plus className="h-4 w-4" />
        </Button>

        {showLeftArrow && (
          <div className="absolute left-12 z-10 flex items-center">
            <div className="h-full w-12 bg-gradient-to-r from-background to-transparent" />
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 bg-background text-muted-foreground hover:text-foreground"
              onClick={() => scroll("left")}
            >
              <ChevronLeft className="h-4 w-4" />
            </Button>
          </div>
        )}

        <div
          ref={scrollRef}
          onScroll={handleScroll}
          className="flex-1 overflow-x-auto scrollbar-hide"
          style={{ scrollbarWidth: "none", msOverflowStyle: "none" }}
        >
          <nav className="flex items-center gap-1 py-2">
            {topics.map((topic) => (
              <Link
                key={topic.id}
                href={topic.href}
                onClick={() => setActiveId(topic.id)}
                className={cn(
                  "shrink-0 whitespace-nowrap px-3 py-2 text-sm transition-colors",
                  activeId === topic.id
                    ? "border-b-2 border-foreground font-medium text-foreground"
                    : "text-muted-foreground hover:text-foreground"
                )}
              >
                {topic.label}
              </Link>
            ))}
          </nav>
        </div>

        {showRightArrow && (
          <div className="absolute right-4 z-10 flex items-center">
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 bg-background text-muted-foreground hover:text-foreground"
              onClick={() => scroll("right")}
            >
              <ChevronRight className="h-4 w-4" />
            </Button>
            <div className="h-full w-12 bg-gradient-to-l from-background to-transparent" />
          </div>
        )}
      </div>
    </div>
  );
}
