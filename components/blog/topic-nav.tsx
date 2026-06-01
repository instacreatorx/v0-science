"use client";

import Link from "next/link";
import { useRef, useState } from "react";
import { useTranslations } from "next-intl";
import { ChevronLeft, ChevronRight, Plus } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";

const topicIds = [
  { id: "for-you", key: "forYou", href: "/" },
  { id: "following", key: "following", href: "/following" },
  { id: "featured", key: "featured", href: "/featured" },
  { id: "technology", key: "technology", href: "/topic/technology" },
  { id: "programming", key: "programming", href: "/topic/programming" },
  { id: "data-science", key: "dataScience", href: "/topic/data-science" },
  { id: "self-improvement", key: "selfImprovement", href: "/topic/self-improvement" },
  { id: "writing", key: "writing", href: "/topic/writing" },
  { id: "productivity", key: "productivity", href: "/topic/productivity" },
  { id: "design", key: "design", href: "/topic/design" },
] as const;

export function TopicNav() {
  const t = useTranslations("topics");
  const scrollRef = useRef<HTMLDivElement>(null);
  const [activeId, setActiveId] = useState("for-you");
  const [showLeftArrow, setShowLeftArrow] = useState(false);
  const [showRightArrow, setShowRightArrow] = useState(true);

  const scroll = (direction: "left" | "right") => {
    if (scrollRef.current) {
      const scrollAmount = 200;
      const isRtl = document.documentElement.dir === "rtl";
      const multiplier = isRtl ? -1 : 1;
      scrollRef.current.scrollBy({
        left: direction === "left" ? -scrollAmount * multiplier : scrollAmount * multiplier,
        behavior: "smooth",
      });
    }
  };

  const handleScroll = () => {
    if (scrollRef.current) {
      const { scrollLeft, scrollWidth, clientWidth } = scrollRef.current;
      const absScroll = Math.abs(scrollLeft);
      setShowLeftArrow(absScroll > 0);
      setShowRightArrow(absScroll < scrollWidth - clientWidth - 10);
    }
  };

  return (
    <div className="relative border-b border-border bg-background">
      <div className="mx-auto flex max-w-5xl items-center px-4">
        <Button
          variant="ghost"
          size="icon"
          className="me-2 h-8 w-8 shrink-0 text-muted-foreground hover:text-foreground"
        >
          <Plus className="h-4 w-4" />
        </Button>

        {showLeftArrow && (
          <div className="absolute z-10 flex items-center ltr:left-12 rtl:right-12">
            <div className="h-full w-12 bg-gradient-to-r from-background to-transparent ltr:block rtl:hidden" />
            <div className="hidden h-full w-12 bg-gradient-to-l from-background to-transparent ltr:hidden rtl:block" />
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 bg-background text-muted-foreground hover:text-foreground"
              onClick={() => scroll("left")}
            >
              <ChevronLeft className="h-4 w-4 rtl:rotate-180" />
            </Button>
          </div>
        )}

        <div
          ref={scrollRef}
          onScroll={handleScroll}
          className="scrollbar-hide flex-1 overflow-x-auto"
          style={{ scrollbarWidth: "none", msOverflowStyle: "none" }}
        >
          <nav className="flex items-center gap-1 py-2">
            {topicIds.map((topic) => (
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
                {t(topic.key)}
              </Link>
            ))}
          </nav>
        </div>

        {showRightArrow && (
          <div className="absolute z-10 flex items-center ltr:right-4 rtl:left-4">
            <Button
              variant="ghost"
              size="icon"
              className="h-8 w-8 bg-background text-muted-foreground hover:text-foreground"
              onClick={() => scroll("right")}
            >
              <ChevronRight className="h-4 w-4 rtl:rotate-180" />
            </Button>
            <div className="h-full w-12 bg-gradient-to-l from-background to-transparent ltr:block rtl:hidden" />
            <div className="hidden h-full w-12 bg-gradient-to-r from-background to-transparent ltr:hidden rtl:block" />
          </div>
        )}
      </div>
    </div>
  );
}
