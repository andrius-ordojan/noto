# Noto

Noto is a tool for building and maintaining a personal knowledge base—without the overhead of writing full notes by hand.

It scans markdown files from my Obsidian vault, detects saved links, and pulls in the content behind them. It combines that source material with any notes or directions I’ve added, and feeds everything into an LLM to generate a concise, structured summary for future reference.

Whether it's technical documentation, cooking tips, philosophy articles, or blog posts - Noto distills the key takeaways into clean, readable notes.

To protect against link rot, all referenced web pages are archived using [ArchiveBox](https://archivebox.io/), ensuring the original context is always preserved—even if the source disappears.

This way, I can continue saving links and jotting quick thoughts, while still ending up with a durable, searchable, and well-organized archive of what I’ve learned.

🛠 Note: Noto is still under active development and not yet ready for general use.

## To Do

- read notes
  - [ ] parse markdown files
    - [ ] find tags
    - [ ] extract links
      - [ ] separate articles from other websites
    - [ ] define format for directives for llm
    - [ ] define format for "chapters"
    - [ ] extract content text
- archive links
  - [ ] submit links to archivebox
  - [ ] download the readability extracted content
  - [ ] retrieve the link to the archived link
- generate summary
  - [ ] create prompt for llm
  - [ ] generate summary using llm
  - [ ] rewrite note
