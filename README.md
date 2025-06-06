# Noto

Noto is a tool for building and maintaining a personal knowledge base—without the overhead of writing full notes by hand.

It scans markdown files from my Obsidian vault. Notes with `noto` property will be added to a filtered list. It will feed the note into an LLM to generate a concise, structured summary for future reference.

For now targeted at technical documents, but with some adjustmets to the prompt should work with more general content.

Goal is to make it easier to take and maintaining notes. I can continue saving links and jotting quick thoughts, while still ending up with a durable, searchable, and well-organized archive of what I’ve learned.

## Future plans

- Once ArchiveBox supports http API I want to backup all linked content.
- Will keep expermenting with different prompts to find the best fit for summarization.
- Long term would like the tool to run as a server reacting to fontmatter properties automaticly
