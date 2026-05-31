import { Article } from "@/components/blog/article-card";

export interface ArticleWithContent extends Article {
  content: string;
  author: Article["author"] & {
    bio?: string;
    followers?: string;
  };
}

export const articleContent: Record<string, string> = {
  "1": `
<p>In an industry that celebrates "move fast and break things," I've come to believe we've lost something important: the craft of building software meant to endure.</p>

<h2>The Problem with Disposable Software</h2>

<p>After fifteen years of writing code professionally, I've watched countless projects rise and fall. Some failed for obvious reasons—bad market fit, poor timing, lack of funding. But many failed for a reason that's rarely discussed: they were built to be temporary.</p>

<p>We've become so obsessed with shipping fast that we've forgotten how to ship well. We reach for the newest framework, the trendiest architecture, the most exciting paradigm—without asking whether these choices serve our users or just our egos.</p>

<blockquote>"The best code is code that doesn't need to be changed. The second best code is code that's easy to change. Everything else is technical debt waiting to happen."</blockquote>

<h2>What Lasting Software Looks Like</h2>

<p>Software that lasts shares certain characteristics:</p>

<ul>
<li><strong>Boring technology choices:</strong> It uses proven tools rather than experimental ones. PostgreSQL over the database-of-the-month. Server-side rendering over whatever JavaScript framework launched last Tuesday.</li>
<li><strong>Deep domain understanding:</strong> The developers took time to understand the problem before jumping to solutions.</li>
<li><strong>Respect for the user:</strong> Every feature exists because users need it, not because developers wanted to build it.</li>
<li><strong>Operational excellence:</strong> Logging, monitoring, and debugging were considered from day one, not bolted on after the first production incident.</li>
</ul>

<h2>The Path Forward</h2>

<p>I'm not suggesting we abandon innovation or resist change. But I am suggesting we ask better questions before we start typing. Questions like: Will this choice still make sense in five years? Am I reaching for this tool because it's right, or because it's new? What would the simplest possible solution look like?</p>

<p>The software industry desperately needs more craftspeople—developers who take pride not just in shipping, but in building things that last. Things that work reliably, that users can depend on, that don't need to be rewritten every eighteen months.</p>

<p>It's harder than building disposable software. It requires more thought, more patience, more discipline. But it's also more rewarding. There's a deep satisfaction in looking at code you wrote five years ago and realizing it still works, still serves its purpose, still helps people.</p>

<p>That's the kind of software I want to build. That's the kind of craftsperson I want to be.</p>
`,
  "2": `
<p>The day I quit, I made $847 before lunch. My Stripe dashboard showed recurring revenue climbing steadily. My newsletter had 12,000 subscribers. And yet, when I typed my resignation letter, my hands were shaking.</p>

<h2>The Golden Handcuffs</h2>

<p>At Big Tech, the compensation was extraordinary. Base salary, stock grants, bonuses—it added up to numbers that would have seemed absurd to my younger self. I had colleagues who'd been there for a decade, not because they loved the work, but because leaving meant walking away from hundreds of thousands in unvested equity.</p>

<p>That's what golden handcuffs do. They make the exit cost so high that staying feels like the only rational choice.</p>

<blockquote>"I realized I was trading years of my life for money I didn't need, doing work I didn't care about, to impress people I didn't like."</blockquote>

<h2>The Breaking Point</h2>

<p>It happened in a meeting. The forty-seventh meeting that week, discussing the same initiative we'd been discussing for months. I looked around the room at smart, capable people spending their lives on something that didn't matter to any of us.</p>

<p>And I thought: Is this really how I want to spend my finite time on Earth?</p>

<h2>The Leap</h2>

<p>I'd been writing on the side for three years. Building an audience slowly, learning what resonated, developing my voice. When I finally looked at the numbers, I realized I could survive on my writing income—not luxuriously, but comfortably.</p>

<p>More importantly, I could wake up every day excited about my work. I could create things that mattered to me. I could be present with my family instead of checking Slack during dinner.</p>

<h2>What I've Learned</h2>

<p>Six months into full-time writing, here's what I know:</p>

<ul>
<li><strong>Money isn't everything:</strong> But financial runway is. I couldn't have made this leap without savings and side income.</li>
<li><strong>Status is a trap:</strong> Nobody at the coffee shop cares where you used to work.</li>
<li><strong>Freedom requires discipline:</strong> Without a boss, you have to become your own accountability system.</li>
<li><strong>Impact is possible anywhere:</strong> I've helped more people with one viral article than I did in a decade at Big Tech.</li>
</ul>

<p>I don't regret my time in tech. It taught me to think at scale, to ship relentlessly, to work with brilliant people. But I'm grateful I found the courage to leave before those golden handcuffs became permanent.</p>
`,
  "3": `
<p>You've read the articles: wake up at 5 AM, meditate, exercise, journal, take a cold shower, review your goals—all before breakfast. The perfect morning routine, they say, will transform your life.</p>

<p>But what if that's not quite right?</p>

<h2>What the Science Actually Shows</h2>

<p>Recent research from Stanford's Sleep Lab tells a different story. The studies reveal that optimal morning routines are far more individual than the productivity gurus suggest—and that many popular practices might actually be counterproductive.</p>

<blockquote>"The most effective morning routine is the one that aligns with your chronotype, not the one that looks best on social media."</blockquote>

<h2>The Chronotype Factor</h2>

<p>Chronotypes—your natural inclination toward being a morning person or night owl—are largely genetic. Forcing yourself to wake at 5 AM when your biology prefers 8 AM doesn't just feel bad; it can impair cognitive function for hours.</p>

<p>Dr. Matthew Walker's research shows that fighting your chronotype leads to:</p>

<ul>
<li>Reduced creative problem-solving ability</li>
<li>Impaired memory consolidation</li>
<li>Increased cortisol levels throughout the day</li>
<li>Higher risk of depression and anxiety</li>
</ul>

<h2>Rethinking the Cold Shower</h2>

<p>Cold exposure has real benefits—improved circulation, mood enhancement, increased alertness. But timing matters. For many people, a cold shower immediately upon waking can spike cortisol at a time when it's already naturally elevated, leading to increased anxiety and decreased immune function.</p>

<h2>What Actually Works</h2>

<p>Based on current neuroscience, here's what we know about effective mornings:</p>

<ul>
<li><strong>Light exposure:</strong> Getting bright light (ideally sunlight) within 30 minutes of waking is one of the most powerful tools for regulating your circadian rhythm.</li>
<li><strong>Hydration:</strong> You lose about a liter of water while sleeping. Drinking water upon waking helps cognitive function more than coffee does.</li>
<li><strong>Movement:</strong> Any movement helps—it doesn't need to be a full workout. Even a 10-minute walk significantly improves focus.</li>
<li><strong>Delayed caffeine:</strong> Waiting 90-120 minutes before coffee prevents the afternoon crash and improves sleep quality.</li>
</ul>

<h2>The Real Secret</h2>

<p>The most effective morning routine is one you'll actually stick to. Complexity is the enemy of consistency. Start with one or two evidence-based practices, do them reliably for a month, then consider adding more.</p>

<p>Your morning doesn't need to be Instagram-worthy. It needs to work for you.</p>
`,
  "4": `
<p>Large Language Models have transformed from research curiosities into tools that millions use daily. But how do they actually work? Let's build an intuition, step by step.</p>

<h2>The Foundation: Tokens and Embeddings</h2>

<p>Before an LLM can do anything with text, it needs to convert words into numbers. This happens in two steps.</p>

<p>First, tokenization breaks text into pieces. These aren't always whole words—"understanding" might become "under" + "stand" + "ing". This allows the model to handle any word, even ones it's never seen before.</p>

<p>Then, each token gets mapped to an embedding: a list of hundreds or thousands of numbers that represent its meaning. Similar words have similar embeddings—"cat" and "kitten" are close together in this mathematical space, while "cat" and "democracy" are far apart.</p>

<blockquote>"An embedding is like GPS coordinates for meaning. Words that mean similar things are in similar locations."</blockquote>

<h2>The Magic: Attention Mechanisms</h2>

<p>The breakthrough behind modern LLMs is the attention mechanism. It allows the model to look at all the words in a sentence simultaneously and figure out which ones are related.</p>

<p>Consider: "The cat sat on the mat because it was tired."</p>

<p>What does "it" refer to? The cat, obviously—but how does a computer know? Attention allows the model to assign weights to each word when processing "it". The cat gets high attention; "mat" gets low attention.</p>

<h2>Transformers: The Architecture</h2>

<p>A transformer stacks multiple attention layers on top of each other. Each layer builds on the previous one's understanding:</p>

<ul>
<li><strong>Early layers:</strong> Learn basic patterns—grammar, syntax, common phrases</li>
<li><strong>Middle layers:</strong> Capture relationships—subjects and verbs, pronouns and referents</li>
<li><strong>Later layers:</strong> Handle abstract concepts—tone, intent, factual relationships</li>
</ul>

<h2>Training: Learning from Text</h2>

<p>LLMs learn by prediction. Given "The quick brown fox jumps over the", the model tries to predict the next word. When it guesses wrong, its parameters are adjusted slightly.</p>

<p>This happens billions of times across terabytes of text. Gradually, the model learns not just which words follow other words, but the deeper patterns of language—and the knowledge encoded within it.</p>

<h2>The Emergent Magic</h2>

<p>Something remarkable happens at scale. When transformers get large enough and see enough data, they develop capabilities they were never explicitly taught: reasoning, summarization, translation, coding.</p>

<p>We don't fully understand why this happens. It's one of the most fascinating open questions in AI research—and it's why these models feel almost magical when you use them.</p>
`,
  "featured-1": `
<p>The future didn't arrive with a bang. It crept in gradually—a better autocomplete here, a surprisingly good image there, a code suggestion that saved an hour of work. And then, one day, we looked around and realized everything had changed.</p>

<h2>The New Partnership</h2>

<p>I spent the last year interviewing people whose jobs have been transformed by AI. Not replaced—transformed. What I found challenged my assumptions about what human-AI collaboration actually looks like.</p>

<p>A architect told me she designs buildings 40% faster now. Not because AI does the creative work—it doesn't—but because it handles the tedious calculations that used to consume her afternoons. She spends more time on what she loves: the creative vision, the human experience of space.</p>

<blockquote>"AI didn't replace my creativity. It replaced the busywork that was burying it."</blockquote>

<h2>The Skills That Matter Now</h2>

<p>The professionals thriving in this new landscape share certain characteristics:</p>

<ul>
<li><strong>Judgment over execution:</strong> They know when AI output is good enough and when it needs human refinement.</li>
<li><strong>Creative direction:</strong> They can articulate what they want clearly enough for AI to help achieve it.</li>
<li><strong>Critical evaluation:</strong> They catch AI mistakes that would fool a novice.</li>
<li><strong>Synthesis:</strong> They combine AI capabilities with human insight in ways neither could achieve alone.</li>
</ul>

<h2>The Amplification Effect</h2>

<p>Here's what surprised me most: AI amplifies existing inequalities in skill and motivation. People who were already good at their jobs became dramatically more productive. People who were coasting found AI wasn't good enough to coast on.</p>

<p>The gap between top performers and average performers is widening. AI is a multiplier, and multiplying by zero still equals zero.</p>

<h2>What Comes Next</h2>

<p>The AI systems we have today are primitive compared to what's coming. Current models hallucinate, miss context, fail at basic reasoning. And yet they're already transforming industries.</p>

<p>The question isn't whether AI will change your field—it will. The question is whether you'll shape that change or be shaped by it.</p>

<p>The future of human-AI collaboration is already here. It's just not evenly distributed yet.</p>
`,
};

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
