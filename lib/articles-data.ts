import { Article } from "@/components/blog/article-card";

export const articles: Article[] = [
  {
    id: "1",
    title: "The Art of Building Software That Lasts",
    excerpt:
      "In a world obsessed with speed and iteration, we&apos;ve forgotten the craft of building software meant to endure. Here&apos;s what I learned after 15 years in the industry.",
    author: {
      name: "Sarah Chen",
      avatar:
        "https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=100&h=100&fit=crop&crop=face",
    },
    publication: "The Startup",
    date: "May 28",
    readTime: "8 min read",
    image:
      "https://images.unsplash.com/photo-1461749280684-dccba630e2f6?w=800&h=600&fit=crop",
    tags: ["Technology", "Programming"],
    isMemberOnly: true,
  },
  {
    id: "2",
    title: "Why I Quit My $300K Job at Big Tech to Write Full-Time",
    excerpt:
      "The golden handcuffs were real, but so was the emptiness. After a decade of climbing the corporate ladder, I finally jumped off—and found something better.",
    author: {
      name: "Marcus Williams",
      avatar:
        "https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=100&h=100&fit=crop&crop=face",
    },
    date: "May 27",
    readTime: "12 min read",
    image:
      "https://images.unsplash.com/photo-1499750310107-5fef28a66643?w=800&h=600&fit=crop",
    tags: ["Career", "Self Improvement"],
  },
  {
    id: "3",
    title: "The Surprising Science Behind Morning Routines",
    excerpt:
      "New research reveals that your morning routine might be doing more harm than good. What neuroscience tells us about optimal daily rituals.",
    author: {
      name: "Dr. Emma Rodriguez",
      avatar:
        "https://images.unsplash.com/photo-1438761681033-6461ffad8d80?w=100&h=100&fit=crop&crop=face",
    },
    publication: "Better Humans",
    date: "May 26",
    readTime: "6 min read",
    image:
      "https://images.unsplash.com/photo-1484627147104-f5197bcd6651?w=800&h=600&fit=crop",
    tags: ["Productivity", "Science"],
    isMemberOnly: true,
  },
  {
    id: "4",
    title: "Understanding Large Language Models: A Visual Guide",
    excerpt:
      "Transformers, attention mechanisms, and embeddings explained through intuitive visualizations. No PhD required.",
    author: {
      name: "Alex Thompson",
      avatar:
        "https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=100&h=100&fit=crop&crop=face",
    },
    publication: "Towards AI",
    date: "May 25",
    readTime: "15 min read",
    image:
      "https://images.unsplash.com/photo-1677442136019-21780ecad995?w=800&h=600&fit=crop",
    tags: ["Machine Learning", "AI"],
  },
  {
    id: "5",
    title: "The Hidden Cost of Hustle Culture Nobody Talks About",
    excerpt:
      "We glorify the grind, but at what cost? A deep dive into the mental health crisis affecting ambitious professionals.",
    author: {
      name: "Jordan Park",
      avatar:
        "https://images.unsplash.com/photo-1500648767791-00dcc994a43e?w=100&h=100&fit=crop&crop=face",
    },
    date: "May 24",
    readTime: "9 min read",
    image:
      "https://images.unsplash.com/photo-1488190211105-8b0e65b80b4e?w=800&h=600&fit=crop",
    tags: ["Mental Health", "Work"],
    isMemberOnly: true,
  },
  {
    id: "6",
    title: "Building Your First AI Agent: A Step-by-Step Tutorial",
    excerpt:
      "From concept to deployment: everything you need to know to build an AI agent that actually works in production.",
    author: {
      name: "Rachel Kim",
      avatar:
        "https://images.unsplash.com/photo-1534528741775-53994a69daeb?w=100&h=100&fit=crop&crop=face",
    },
    publication: "Level Up Coding",
    date: "May 23",
    readTime: "18 min read",
    image:
      "https://images.unsplash.com/photo-1555255707-c07966088b7b?w=800&h=600&fit=crop",
    tags: ["Programming", "AI"],
  },
  {
    id: "7",
    title: "What 10 Years of Remote Work Taught Me About Productivity",
    excerpt:
      "Long before the pandemic, I was working remotely. Here are the hard-won lessons that most remote work guides miss.",
    author: {
      name: "David Miller",
      avatar:
        "https://images.unsplash.com/photo-1519085360753-af0119f7cbe7?w=100&h=100&fit=crop&crop=face",
    },
    date: "May 22",
    readTime: "7 min read",
    image:
      "https://images.unsplash.com/photo-1522071820081-009f0129c71c?w=800&h=600&fit=crop",
    tags: ["Remote Work", "Productivity"],
  },
  {
    id: "8",
    title: "The Minimalist Approach to System Design",
    excerpt:
      "Why the best architectures are often the simplest ones. A guide to avoiding over-engineering in modern software development.",
    author: {
      name: "Nina Patel",
      avatar:
        "https://images.unsplash.com/photo-1580489944761-15a19d654956?w=100&h=100&fit=crop&crop=face",
    },
    publication: "System Design Weekly",
    date: "May 21",
    readTime: "11 min read",
    image:
      "https://images.unsplash.com/photo-1517694712202-14dd9538aa97?w=800&h=600&fit=crop",
    tags: ["Architecture", "Engineering"],
    isMemberOnly: true,
  },
];

export const featuredArticle: Article = {
  id: "featured-1",
  title: "The Future of Human-AI Collaboration Is Already Here",
  excerpt:
    "We stand at an inflection point. AI is no longer a distant future—it&apos;s reshaping how we work, create, and think today. A deep exploration of what comes next.",
  author: {
    name: "Dr. James Liu",
    avatar:
      "https://images.unsplash.com/photo-1560250097-0b93528c311a?w=100&h=100&fit=crop&crop=face",
  },
  publication: "Future Vision",
  date: "May 29",
  readTime: "20 min read",
  image:
    "https://images.unsplash.com/photo-1485827404703-89b55fcc595e?w=1600&h=900&fit=crop",
  tags: ["AI", "Future", "Technology"],
  isMemberOnly: true,
};
