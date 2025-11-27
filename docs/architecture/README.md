# Architecture diagrams

Files:

- `microservices.md` — Mermaid diagram showing a draft microservices layout (API Gateway, Auth, User, Student, Course, Billing, Notification, Search, datastores, message broker, monitoring and CI/CD).

How to view:

- In VS Code: open `microservices.md` and use Markdown Preview (press `Ctrl+Shift+V` or `Ctrl+K V`).
- Online: copy the Mermaid code block into https://mermaid.live to render and export PNG/SVG.
- CLI: use `mmdc` (Mermaid CLI) to render locally. Example:

```bash
# Install mermaid-cli (requires Node.js)
npm install -g @mermaid-js/mermaid-cli

# Render to PNG
mmdc -i docs/architecture/microservices.md -o docs/architecture/microservices.png
```

Next steps:

- Tell me which services or interactions you want adjusted (names, databases, sync vs async patterns, external integrations).
- I can generate a PlantUML or SVG export if you prefer a different format.
