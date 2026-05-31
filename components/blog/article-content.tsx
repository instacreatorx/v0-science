interface ArticleContentProps {
  content: string;
}

export function ArticleContent({ content }: ArticleContentProps) {
  return (
    <article 
      className="prose prose-lg prose-neutral mx-auto max-w-2xl px-4 
        prose-headings:font-serif prose-headings:font-bold prose-headings:text-foreground
        prose-h2:mt-10 prose-h2:mb-4 prose-h2:text-2xl md:prose-h2:text-3xl
        prose-p:text-foreground prose-p:leading-relaxed prose-p:mb-6
        prose-a:text-primary prose-a:no-underline hover:prose-a:underline
        prose-blockquote:border-l-4 prose-blockquote:border-primary prose-blockquote:bg-muted/50
        prose-blockquote:py-4 prose-blockquote:px-6 prose-blockquote:not-italic
        prose-blockquote:text-foreground prose-blockquote:font-serif prose-blockquote:text-xl
        prose-blockquote:my-8 prose-blockquote:rounded-r
        prose-strong:text-foreground prose-strong:font-semibold
        prose-ul:my-6 prose-ul:pl-6
        prose-li:text-foreground prose-li:mb-2 prose-li:marker:text-muted-foreground"
      dangerouslySetInnerHTML={{ __html: content }}
    />
  );
}
