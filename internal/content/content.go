package content

// Section represents a navigable section of the portfolio.
type Section struct {
	ID          string // URL slug: "technical-notes", "projects", etc.
	Title       string
	Subtitle    string
	NavIcon     string // SVG markup for the homepage nav card
	Cards       []Card
}

// Card represents a single content card within a section.
type Card struct {
	ID          string // URL slug within section
	Title       string
	Subtitle    string
	Description string // ≤60 chars, shown on card
	IconSVG     string // SVG for card image area
	Detail      string // HTML content for detail view
}

// Sections returns all portfolio sections with their cards.
func Sections() []Section {
	return []Section{
		technicalNotes(),
		projects(),
		musings(),
		theBullshitters(),
	}
}

// SectionByID looks up a section by its slug.
func SectionByID(id string) (Section, bool) {
	for _, s := range Sections() {
		if s.ID == id {
			return s, true
		}
	}
	return Section{}, false
}

// CardByID looks up a card within a section.
func CardByID(sectionID, cardID string) (Section, Card, bool) {
	s, ok := SectionByID(sectionID)
	if !ok {
		return Section{}, Card{}, false
	}
	for _, c := range s.Cards {
		if c.ID == cardID {
			return s, c, true
		}
	}
	return s, Card{}, false
}

func technicalNotes() Section {
	return Section{
		ID:       "technical-notes",
		Title:    "Technical Notes",
		Subtitle: "Lessons learned, patterns documented and things worth remembering.",
		NavIcon: `<svg viewBox="0 0 64 64" fill="none"><rect x="12" y="8" width="36" height="46" rx="2" stroke="#1a1a1a" stroke-width="2" fill="none"/><rect x="16" y="12" width="36" height="46" rx="2" stroke="#1a1a1a" stroke-width="2" fill="#fafafa"/><line x1="24" y1="24" x2="44" y2="24" stroke="#1a1a1a" stroke-width="1.5"/><line x1="24" y1="31" x2="44" y2="31" stroke="#1a1a1a" stroke-width="1.5"/><line x1="24" y1="38" x2="38" y2="38" stroke="#1a1a1a" stroke-width="1.5"/><circle cx="22" cy="48" r="3" stroke="#c44" stroke-width="1.5" fill="none"/></svg>`,
		Cards: []Card{
			{
				ID:          "go-std-lib",
				Title:       "Building with Go's Standard Library",
				Subtitle:    "Go · Architecture",
				Description: "Why Go's standard library is enough for most web apps.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="25" y="15" width="70" height="60" rx="3" stroke="#fafafa" stroke-width="2" fill="none"/><text x="38" y="52" font-family="monospace" font-size="18" fill="#c44">{ }</text><line x1="30" y1="25" x2="50" y2="25" stroke="#fafafa" stroke-width="1.5"/></svg>`,
				Detail: `<p>Go's standard library is one of the most underrated tools in a developer's arsenal. While the ecosystem is full of frameworks, the standard library alone provides everything you need to build production-grade web applications.</p>
<h2>Why Standard Library?</h2>
<p>The <code>net/http</code> package gives you a fully capable HTTP server. The <code>html/template</code> package provides safe, composable HTML templating. The <code>encoding/json</code> package handles serialisation. The <code>database/sql</code> package provides a clean interface to any SQL database.</p>
<h2>Project Structure</h2>
<ul><li>A <code>cmd/</code> directory for entry points</li><li>An <code>internal/</code> directory for business logic</li><li>A <code>views/</code> directory for templates</li><li>A <code>static/</code> directory for assets</li></ul>
<h2>The Trade-Off</h2>
<p>You write slightly more boilerplate upfront. In exchange, you get zero hidden magic, complete control over your request lifecycle, and a codebase any Go developer can read without learning a framework first.</p>`,
			},
			{
				ID:          "htmx-patterns",
				Title:       "HTMX Patterns for Server-Driven UIs",
				Subtitle:    "HTMX · Frontend",
				Description: "Partial swaps, lazy loading and progressive enhancement.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><text x="22" y="55" font-family="monospace" font-size="14" fill="#fafafa">&lt;/&gt;</text><line x1="60" y1="20" x2="60" y2="70" stroke="#fafafa" stroke-width="1" stroke-dasharray="4 3"/><rect x="65" y="30" width="35" height="30" rx="3" stroke="#c44" stroke-width="1.5" fill="none"/></svg>`,
				Detail: `<p>HTMX lets you build dynamic interfaces by returning HTML from the server instead of JSON. It's a return to the architecture the web was designed for, with modern capabilities layered on top.</p>
<h2>Core Patterns</h2>
<ul><li>Use <code>hx-get</code> and <code>hx-swap</code> for partial page updates</li><li>Use <code>hx-trigger="revealed"</code> for lazy loading</li><li>Use <code>hx-push-url</code> to maintain browser history</li><li>Use <code>hx-indicator</code> for loading states</li></ul>
<h2>Progressive Enhancement</h2>
<p>The best HTMX applications work without JavaScript enabled. Every link and form should function as standard HTML first. HTMX then enhances the experience. If the JavaScript fails to load, the page still works.</p>`,
			},
			{
				ID:          "k8s",
				Title:       "Kubernetes from the Ground Up",
				Subtitle:    "Cloud · Infrastructure",
				Description: "Container orchestration fundamentals demystified.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><circle cx="60" cy="45" r="25" stroke="#fafafa" stroke-width="2" fill="none"/><circle cx="60" cy="45" r="10" stroke="#c44" stroke-width="1.5" fill="none"/><line x1="60" y1="20" x2="60" y2="30" stroke="#fafafa" stroke-width="1.5"/><line x1="60" y1="60" x2="60" y2="70" stroke="#fafafa" stroke-width="1.5"/><line x1="35" y1="45" x2="45" y2="45" stroke="#fafafa" stroke-width="1.5"/><line x1="75" y1="45" x2="85" y2="45" stroke="#fafafa" stroke-width="1.5"/></svg>`,
				Detail: `<p>These are my notes from studying for the Kubernetes Cloud Native Associate certification — written for clarity, not certification prep.</p>
<h2>The Mental Model</h2>
<p>Think of Kubernetes as an operating system for your infrastructure. You declare the desired state, and Kubernetes works to make reality match.</p>
<h2>Key Concepts</h2>
<ul><li>Pods — one or more containers sharing a network namespace</li><li>Deployments — manage scaling, rolling updates, rollbacks</li><li>Services — stable network endpoints for your pods</li><li>ConfigMaps and Secrets — separate configuration from code</li></ul>
<h2>What I Learned</h2>
<p>The biggest lesson was understanding when Kubernetes is the right tool and when it's over-engineering. Not every application needs container orchestration.</p>`,
			},
			{
				ID:          "sql-patterns",
				Title:       "SQL Patterns That Scale",
				Subtitle:    "Databases · Performance",
				Description: "Indexing, query patterns and schema design lessons.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><ellipse cx="60" cy="30" rx="35" ry="12" stroke="#fafafa" stroke-width="2" fill="none"/><line x1="25" y1="30" x2="25" y2="60" stroke="#fafafa" stroke-width="2"/><line x1="95" y1="30" x2="95" y2="60" stroke="#fafafa" stroke-width="2"/><ellipse cx="60" cy="60" rx="35" ry="12" stroke="#fafafa" stroke-width="2" fill="none"/><ellipse cx="60" cy="45" rx="35" ry="12" stroke="#c44" stroke-width="1" fill="none" stroke-dasharray="4 3"/></svg>`,
				Detail: `<p>Writing SQL that works is easy. Writing SQL that scales is a discipline.</p>
<h2>Indexing Strategy</h2>
<p>Indexes are not free. Every index speeds up reads but slows down writes. Index the columns in your WHERE clauses and JOIN conditions.</p>
<h2>Query Patterns</h2>
<ul><li>Use EXISTS instead of IN for subqueries</li><li>Avoid SELECT * in production</li><li>Use EXPLAIN ANALYZE to understand query plans</li><li>Paginate with keyset pagination for large datasets</li></ul>
<h2>Schema Design</h2>
<p>Normalise first, denormalise intentionally. Start with a clean relational model, then introduce denormalisation only with measured evidence.</p>`,
			},
			{
				ID:          "error-handling",
				Title:       "Error Handling as a First-Class Concern",
				Subtitle:    "Go · Reliability",
				Description: "Why explicit error handling beats try-catch.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="20" y="15" width="80" height="60" rx="4" stroke="#fafafa" stroke-width="2" fill="none"/><line x1="20" y1="35" x2="100" y2="35" stroke="#fafafa" stroke-width="1"/><text x="30" y="29" font-family="monospace" font-size="9" fill="#c44">err !=</text><text x="68" y="29" font-family="monospace" font-size="9" fill="#fafafa">nil</text><line x1="30" y1="48" x2="75" y2="48" stroke="#fafafa" stroke-width="1.5"/><line x1="30" y1="56" x2="60" y2="56" stroke="#fafafa" stroke-width="1.5"/><line x1="30" y1="64" x2="50" y2="64" stroke="#fafafa" stroke-width="1"/></svg>`,
				Detail: `<p>In most languages, error handling is an afterthought. Go takes a different approach: errors are values, and handling them is part of writing the code.</p>
<h2>Why Go Gets This Right</h2>
<p>The <code>if err != nil</code> pattern is verbose. That's the point. It forces you to think about what happens when things go wrong at every step.</p>
<h2>Patterns That Work</h2>
<ul><li>Wrap errors with context using <code>fmt.Errorf("doing X: %w", err)</code></li><li>Define sentinel errors for conditions callers need to check</li><li>Use custom error types for structured information</li><li>Log at the boundary, handle where you have context</li></ul>
<h2>The Bigger Principle</h2>
<p>Reliable software fails predictably, communicates clearly, and recovers gracefully. Treating error handling as first-class is what separates production code from prototypes.</p>`,
			},
		},
	}
}

func projects() Section {
	return Section{
		ID:       "projects",
		Title:    "Projects",
		Subtitle: "Things I've built — from systems architecture to full-stack applications.",
		NavIcon: `<svg viewBox="0 0 64 64" fill="none"><rect x="8" y="8" width="48" height="48" rx="2" stroke="#1a1a1a" stroke-width="2" fill="none"/><rect x="14" y="14" width="16" height="16" rx="1" stroke="#1a1a1a" stroke-width="1.5" fill="none"/><rect x="34" y="14" width="16" height="8" rx="1" stroke="#1a1a1a" stroke-width="1.5" fill="none"/><rect x="34" y="26" width="16" height="4" rx="1" stroke="#c44" stroke-width="1.5" fill="none"/><rect x="14" y="34" width="36" height="16" rx="1" stroke="#1a1a1a" stroke-width="1.5" fill="none"/><line x1="14" y1="42" x2="50" y2="42" stroke="#1a1a1a" stroke-width="1"/></svg>`,
		Cards: []Card{
			{
				ID: "sendit", Title: "Sendit", Subtitle: "Full-Stack",
				Description: "A rapid-delivery courier service with tracking.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="20" y="20" width="80" height="50" rx="4" stroke="#fafafa" stroke-width="2" fill="none"/><polyline points="20,20 60,50 100,20" stroke="#c44" stroke-width="2" fill="none"/><line x1="20" y1="70" x2="45" y2="50" stroke="#fafafa" stroke-width="1.5"/><line x1="100" y1="70" x2="75" y2="50" stroke="#fafafa" stroke-width="1.5"/></svg>`,
				Detail:      `<p>Sendit is a rapid-delivery courier service built to demonstrate full-stack capabilities with real-time user experience and reliable backend processing.</p><h2>Architecture</h2><p>Client-server architecture with a React frontend and RESTful API. Users can create delivery orders, track parcels in real time, and manage delivery history.</p><h2>Technical Highlights</h2><ul><li>Real-time parcel tracking with status updates</li><li>Role-based access control</li><li>Responsive design across all devices</li><li>RESTful API with proper error handling</li></ul><h2>What I Learned</h2><p>The importance of thinking through the entire user journey before writing code.</p>`,
			},
			{
				ID: "andika", Title: "Andika", Subtitle: "Back-end",
				Description: "A notes management service built for clarity.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="30" y="10" width="60" height="70" rx="3" stroke="#fafafa" stroke-width="2" fill="none"/><line x1="40" y1="28" x2="80" y2="28" stroke="#fafafa" stroke-width="1.5"/><line x1="40" y1="38" x2="80" y2="38" stroke="#fafafa" stroke-width="1.5"/><line x1="40" y1="48" x2="65" y2="48" stroke="#c44" stroke-width="1.5"/><line x1="40" y1="58" x2="75" y2="58" stroke="#fafafa" stroke-width="1.5"/></svg>`,
				Detail:      `<p>Andika — Swahili for 'write' — is a notes management service designed around speed and simplicity.</p><h2>Design Decisions</h2><p>Notes are plain text with minimal metadata. No folders, no tags. You write, you save, you search.</p><h2>Technical Stack</h2><ul><li>RESTful API with CRUD operations</li><li>Full-text search across all notes</li><li>User authentication and authorisation</li><li>Efficient pagination for large collections</li></ul>`,
			},
			{
				ID: "o-sipital", Title: "O-Sipital", Subtitle: "Command Line",
				Description: "Hospital management from the command line.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="15" y="25" width="90" height="45" rx="3" stroke="#fafafa" stroke-width="2" fill="none"/><line x1="15" y1="38" x2="105" y2="38" stroke="#fafafa" stroke-width="1"/><text x="25" y="34" font-family="monospace" font-size="8" fill="#c44">&gt;_</text><line x1="30" y1="50" x2="70" y2="50" stroke="#fafafa" stroke-width="1.5"/><line x1="30" y1="58" x2="55" y2="58" stroke="#fafafa" stroke-width="1.5"/></svg>`,
				Detail:      `<p>O-Sipital is a hospital management system built entirely for the command line.</p><h2>Features</h2><ul><li>Patient registration and record management</li><li>Appointment scheduling with conflict detection</li><li>Doctor and department management</li><li>Secure role-based permissions</li></ul><h2>Why CLI?</h2><p>The command line works over SSH, on low-bandwidth connections, on any terminal. For environments where reliability matters more than aesthetics, a well-designed CLI can be more practical than a web application.</p>`,
			},
			{
				ID: "lawnbull", Title: "Lawnbull", Subtitle: "Front-End",
				Description: "A digital marketing service for brand identity.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><circle cx="60" cy="45" r="30" stroke="#fafafa" stroke-width="2" fill="none"/><circle cx="60" cy="45" r="18" stroke="#fafafa" stroke-width="1.5" fill="none"/><circle cx="60" cy="45" r="6" stroke="#c44" stroke-width="2" fill="none"/></svg>`,
				Detail:      `<p>Lawnbull is a digital marketing service website showcasing frontend capabilities — clean visual design, responsive layout and strong brand identity.</p><h2>Technical Highlights</h2><ul><li>Fully responsive design</li><li>Semantic HTML with accessibility considerations</li><li>CSS animations for micro-interactions</li><li>Performance-optimised asset loading</li></ul>`,
			},
			{
				ID: "mwaniki-dev", Title: "mwaniki.dev", Subtitle: "Go + HTMX + CSS",
				Description: "This portfolio — Go, HTMX and pure CSS.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="20" y="15" width="80" height="60" rx="4" stroke="#fafafa" stroke-width="2" fill="none"/><rect x="20" y="15" width="80" height="14" rx="4" stroke="#fafafa" stroke-width="1.5" fill="none"/><circle cx="30" cy="22" r="2.5" fill="#c44"/><circle cx="38" cy="22" r="2.5" fill="#fafafa"/><circle cx="46" cy="22" r="2.5" fill="#fafafa"/></svg>`,
				Detail:      `<p>This portfolio is a complete overhaul built with Go 1.24 standard library, HTMX and pure CSS. No frameworks, no bundlers, no build step.</p><h2>Design Principles</h2><ul><li>Minimalism — every element earns its place</li><li>Reusability — components that work across pages</li><li>Inclusivity — accessible by default</li><li>Clarity — easy to understand over clever</li></ul><h2>Architecture</h2><p>A single Go binary using <code>net/http</code> for routing and <code>html/template</code> for rendering. HTMX handles partial page swaps. CSS handles all layout, animation and responsive behaviour.</p>`,
			},
		},
	}
}

func musings() Section {
	return Section{
		ID:       "musings",
		Title:    "Musings",
		Subtitle: "Thoughts on software, design, work and everything in between.",
		NavIcon: `<svg viewBox="0 0 64 64" fill="none"><circle cx="32" cy="28" r="18" stroke="#1a1a1a" stroke-width="2" fill="none"/><circle cx="22" cy="50" r="4" stroke="#1a1a1a" stroke-width="1.5" fill="none"/><circle cx="14" cy="56" r="2" stroke="#1a1a1a" stroke-width="1.5" fill="none"/><line x1="24" y1="24" x2="40" y2="24" stroke="#1a1a1a" stroke-width="1.5"/><line x1="24" y1="30" x2="36" y2="30" stroke="#c44" stroke-width="1.5"/></svg>`,
		Cards: []Card{
			{
				ID: "simplicity", Title: "The Case for Simplicity", Subtitle: "Design · Philosophy",
				Description: "The most impactful systems have the fewest parts.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><circle cx="60" cy="45" r="28" stroke="#fafafa" stroke-width="2" fill="none"/><circle cx="60" cy="45" r="4" fill="#c44"/></svg>`,
				Detail:      `<p>There's a recurring pattern in software: we reach for complexity when simplicity would serve better.</p><p>Simplicity is not the absence of capability. It's the discipline of choosing the right amount.</p><h2>What Simplicity Looks Like</h2><ul><li>A function that does one thing with a clear name</li><li>A data model with no redundant fields</li><li>An API with consistent conventions</li><li>A deployment you can explain in three sentences</li></ul><h2>The Cost of Complexity</h2><p>Every layer of complexity is a tax on the future — on onboarding, on debugging at 2am, on changing direction. The most senior thing you can do is resist adding complexity until you have evidence it's necessary.</p>`,
			},
			{
				ID: "public-infrastructure", Title: "Software as Public Infrastructure", Subtitle: "Systems · Society",
				Description: "What software can learn from roads and bridges.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="20" y="55" width="20" height="20" rx="2" stroke="#fafafa" stroke-width="2" fill="none"/><rect x="50" y="40" width="20" height="35" rx="2" stroke="#fafafa" stroke-width="2" fill="none"/><rect x="80" y="25" width="20" height="50" rx="2" stroke="#fafafa" stroke-width="2" fill="none"/><line x1="30" y1="55" x2="60" y2="40" stroke="#c44" stroke-width="1.5"/><line x1="60" y1="40" x2="90" y2="25" stroke="#c44" stroke-width="1.5"/></svg>`,
				Detail:      `<p>We don't think of roads as products. We think of them as infrastructure. Software should aspire to the same standard.</p><h2>What Public Infrastructure Gets Right</h2><ul><li>Designed for everyone, not just power users</li><li>Maintained incrementally over decades</li><li>Prioritises reliability over novelty</li><li>Boring on purpose — predictability is a feature</li></ul><h2>Applying This to Software</h2><p>When I approach a project, I ask: if this were a bridge, would I trust it in ten years? Would someone who didn't build it understand how it works? This changes what you optimise for.</p>`,
			},
			{
				ID: "craft-vs-speed", Title: "Craft vs Speed", Subtitle: "Engineering · Culture",
				Description: "Shipping fast vs building well — when each wins.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><line x1="20" y1="70" x2="100" y2="70" stroke="#fafafa" stroke-width="1.5"/><line x1="20" y1="70" x2="20" y2="15" stroke="#fafafa" stroke-width="1.5"/><polyline points="25,60 45,50 65,55 85,25 100,30" stroke="#c44" stroke-width="2" fill="none"/></svg>`,
				Detail:      `<p>Every engineering team lives in tension between doing it right and doing it now.</p><h2>When Speed Wins</h2><p>Early-stage products need to find their audience. Ship the minimum viable thing, learn from real users, then invest in quality where it matters.</p><h2>When Craft Wins</h2><p>Once you've found product-market fit, craft becomes essential. Technical debt compounds.</p><h2>The False Dichotomy</h2><p>The best engineers write clean code quickly because they've invested in mastering their tools. Speed and quality aren't opposites — they're both outcomes of mastery.</p>`,
			},
			{
				ID: "every-detail", Title: "Every Detail Matters", Subtitle: "Mindset · Practice",
				Description: "Caring about the small things that compound.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="30" y="20" width="60" height="50" rx="3" stroke="#fafafa" stroke-width="2" fill="none"/><line x1="30" y1="35" x2="90" y2="35" stroke="#fafafa" stroke-width="1"/><circle cx="60" cy="55" r="8" stroke="#c44" stroke-width="1.5" fill="none"/></svg>`,
				Detail:      `<p>The gap between good and great software lives in the details — consistent naming, helpful error messages, the animation that's 200ms instead of 400ms.</p><h2>Small Things That Compound</h2><ul><li>Variable names that read like documentation</li><li>Consistent spacing across the entire codebase</li><li>Error states as well-designed as success states</li><li>Commit messages that explain why, not just what</li></ul><h2>The Discipline</h2><p>Attention to detail is not perfectionism. Perfectionism prevents shipping. Attention to detail means doing each thing well as you go.</p>`,
			},
			{
				ID: "learning-in-public", Title: "Learning in Public", Subtitle: "Growth · Community",
				Description: "Sharing what you know before you feel ready.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="15" y="20" width="40" height="50" rx="3" stroke="#fafafa" stroke-width="2" fill="none"/><rect x="65" y="20" width="40" height="50" rx="3" stroke="#fafafa" stroke-width="2" fill="none"/><line x1="25" y1="35" x2="45" y2="35" stroke="#fafafa" stroke-width="1.5"/><line x1="25" y1="43" x2="40" y2="43" stroke="#fafafa" stroke-width="1"/><line x1="75" y1="35" x2="95" y2="35" stroke="#fafafa" stroke-width="1.5"/><line x1="75" y1="43" x2="90" y2="43" stroke="#fafafa" stroke-width="1"/><path d="M55,40 L60,35 L65,40" stroke="#c44" stroke-width="1.5" fill="none"/><path d="M55,50 L60,55 L65,50" stroke="#c44" stroke-width="1.5" fill="none"/></svg>`,
				Detail:      `<p>There's a difference between learning and learning in public. The first is private and safe. The second is vulnerable and far more valuable.</p><h2>Why It Works</h2><p>Writing about what you're learning forces you to understand it properly. The act of explaining is itself a form of deeper learning.</p><h2>What Holds People Back</h2><ul><li>Fear of being wrong in front of others</li><li>Feeling like you need to be an expert first</li><li>Comparing your early understanding to someone else's polished output</li><li>Overthinking the format instead of just starting</li></ul><h2>How I Approach It</h2><p>This entire portfolio is an exercise in learning in public. The notes aren't written from mastery — they're written from honest engagement with the material.</p>`,
			},
		},
	}
}

func theBullshitters() Section {
	return Section{
		ID:       "the-bullshitters",
		Title:    "The Bullshitters",
		Subtitle: "Calling it out — hype, buzzwords and things that don't hold up under scrutiny.",
		NavIcon: `<svg viewBox="0 0 64 64" fill="none"><polygon points="12,26 12,38 24,42 24,22" stroke="#1a1a1a" stroke-width="2" fill="none" stroke-linejoin="round"/><path d="M24,22 L48,10 L48,54 L24,42" stroke="#1a1a1a" stroke-width="2" fill="none" stroke-linejoin="round"/><line x1="48" y1="26" x2="56" y2="24" stroke="#c44" stroke-width="1.5"/><line x1="48" y1="32" x2="58" y2="32" stroke="#c44" stroke-width="1.5"/><line x1="48" y1="38" x2="56" y2="40" stroke="#c44" stroke-width="1.5"/></svg>`,
		Cards: []Card{
			{
				ID: "ai-everything", Title: `Slapping "AI" on Everything`, Subtitle: "Hype · Marketing",
				Description: "When the label matters more than the capability.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="25" y="20" width="70" height="50" rx="3" stroke="#fafafa" stroke-width="2" fill="none"/><text x="38" y="52" font-family="monospace" font-size="16" fill="#c44">AI</text><circle cx="40" cy="32" r="5" stroke="#fafafa" stroke-width="1.5" fill="none"/><circle cx="80" cy="32" r="5" stroke="#fafafa" stroke-width="1.5" fill="none"/><line x1="45" y1="32" x2="75" y2="32" stroke="#fafafa" stroke-width="1" stroke-dasharray="3 2"/></svg>`,
				Detail:      `<p>Somewhere around 2023, "AI" stopped being a technology descriptor and became a marketing strategy.</p><h2>The Pattern</h2><p>Take an existing product. Add a chatbot. Rebrand the landing page. Raise a round. The technology underneath often hasn't changed — only the label.</p><h2>How to Spot the Difference</h2><ul><li>Ask what the product did before AI was added</li><li>Look for specifics — "AI-powered" means nothing</li><li>Check whether removing the AI feature would make it meaningfully worse</li></ul><p>Real AI capabilities are transformative. But most of what gets labelled AI is just software with better marketing.</p>`,
			},
			{
				ID: "10x-engineer", Title: "The Myth of the 10x Engineer", Subtitle: "Culture · Industry",
				Description: "Why hero worship in teams does more harm than good.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><text x="25" y="58" font-family="monospace" font-size="36" fill="#fafafa">10</text><text x="72" y="58" font-family="monospace" font-size="24" fill="#c44">x</text><line x1="20" y1="68" x2="100" y2="68" stroke="#fafafa" stroke-width="1" stroke-dasharray="4 3"/></svg>`,
				Detail:      `<p>The "10x engineer" myth is one of the most persistent and damaging in our industry.</p><h2>What Actually Happens</h2><p>Developers labelled 10x usually write a lot of code quickly. But output isn't impact.</p><h2>The Damage</h2><ul><li>Rewards individual heroics over collaboration</li><li>Gives cover for poor communication and unreadable code</li><li>Discourages mentorship</li><li>Burns people out</li></ul><p>Great engineering is a team sport. The best engineers don't make themselves 10x. They make everyone around them 2x.</p>`,
			},
			{
				ID: "move-fast", Title: "Move Fast and Break People", Subtitle: "Ethics · Startups",
				Description: "The human cost of speed without care.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><polygon points="60,15 95,75 25,75" stroke="#fafafa" stroke-width="2" fill="none"/><text x="51" y="60" font-family="monospace" font-size="20" fill="#c44">!</text><line x1="60" y1="35" x2="60" y2="48" stroke="#fafafa" stroke-width="2"/></svg>`,
				Detail:      `<p>"Move fast and break things" was always convenient for people who didn't live with the consequences.</p><h2>The Real Cost</h2><p>Rushed deployments breaking production. Privacy as an afterthought. Accessibility as a nice-to-have.</p><h2>A Better Model</h2><ul><li>Move deliberately and build things that last</li><li>Ship frequently, but test before you ship</li><li>Treat reliability as a feature</li><li>Ask "who pays if this goes wrong?" before every shortcut</li></ul><p>Speed matters. But speed without care is just recklessness wearing a hoodie.</p>`,
			},
			{
				ID: "framework-treadmill", Title: "The Framework Treadmill", Subtitle: "JavaScript · Fatigue",
				Description: "Churn disguised as progress.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><circle cx="40" cy="45" r="15" stroke="#fafafa" stroke-width="2" fill="none"/><circle cx="80" cy="45" r="15" stroke="#fafafa" stroke-width="2" fill="none"/><path d="M55,42 C62,35 68,35 75,42" stroke="#c44" stroke-width="1.5" fill="none"/><path d="M55,48 C62,55 68,55 75,48" stroke="#c44" stroke-width="1.5" fill="none"/></svg>`,
				Detail:      `<p>Every 18 months, a new JavaScript framework promises to fix everything the last one got wrong. This isn't progress — it's a treadmill.</p><h2>The Cycle</h2><p>Framework A gains popularity. Framework B points out A's flaws. Developers migrate. Framework C arrives. Meanwhile, the product hasn't improved.</p><h2>What Stability Looks Like</h2><ul><li>Choose tools with long-term maintenance commitments</li><li>Evaluate by upgrade path, not launch features</li><li>Question whether you need a framework at all</li><li>Measure migration cost in engineer-months</li></ul>`,
			},
			{
				ID: "thought-leaders", Title: "Thought Leaders Who Don't Ship", Subtitle: "Commentary · Reality",
				Description: "The gap between posting and actually shipping.",
				IconSVG:     `<svg viewBox="0 0 120 90" fill="none"><rect x="20" y="25" width="80" height="45" rx="4" stroke="#fafafa" stroke-width="2" fill="none"/><line x1="30" y1="42" x2="90" y2="42" stroke="#fafafa" stroke-width="1"/><line x1="30" y1="50" x2="75" y2="50" stroke="#fafafa" stroke-width="1"/><circle cx="90" cy="25" r="8" stroke="#c44" stroke-width="1.5" fill="none"/><text x="87" y="29" font-family="monospace" font-size="10" fill="#c44">?</text></svg>`,
				Detail:      `<p>There's a growing class in tech who built careers talking about building things rather than actually building them.</p><h2>The Tell</h2><ul><li>Speak in absolutes about technologies they haven't used in production</li><li>Advice always high-level enough to be unfalsifiable</li><li>Reference "at scale" without specifying what scale means</li><li>Most recent hands-on work was several years ago</li></ul><h2>Why It Matters</h2><p>Junior developers look to these figures for guidance. When that guidance comes from someone who hasn't debugged a production issue in five years, it creates a gap between theory and reality.</p><p>The people worth listening to still ship code, review pull requests, and get paged at inconvenient hours.</p>`,
			},
		},
	}
}
