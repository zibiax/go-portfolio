<script>
    let { project } = $props();

    const languageColors = {
        'javascript': '#f1e05a',
        'typescript': '#2b7489',
        'python': '#3572A5',
        'go': '#00ADD8',
        'java': '#b07219',
        'html': '#e34c26',
        'css': '#563d7c',
        'rust': '#dea584',
        'c': '#555555',
        'c++': '#f34b7d',
        'c#': '#178600',
        'shell': '#89e051',
        'svelte': '#ff3e00',
        'lua': '#000080',
        'default': '#6e7681'
    };

    function getLanguageColor(lang) {
        if (!lang) return languageColors.default;
        return languageColors[lang.toLowerCase()] || languageColors.default;
    }

    function formatUpdated(value) {
        if (!value) return '';
        const date = new Date(value);
        if (Number.isNaN(date.getTime())) return '';

        const days = Math.floor((Date.now() - date.getTime()) / 86400000);
        if (days <= 0) return 'Updated today';
        if (days === 1) return 'Updated yesterday';
        if (days < 30) return `Updated ${days} days ago`;
        if (days < 365) {
            const months = Math.round(days / 30);
            return `Updated ${months} month${months === 1 ? '' : 's'} ago`;
        }
        const years = Math.round(days / 365);
        return `Updated ${years} year${years === 1 ? '' : 's'} ago`;
    }

    const updated = $derived(formatUpdated(project.lastUpdated));
    const topics = $derived((project.topics || []).slice(0, 3));
</script>

<article class="card">
    <div class="card-content">
        <div class="header">
            <h3 class="project-name">{project.name}</h3>
            {#if project.stars > 0}
                <span class="stars" title="{project.stars} stars">
                    <svg viewBox="0 0 16 16" width="13" height="13" aria-hidden="true" fill="currentColor">
                        <path d="M8 .25l2.06 4.17 4.6.67-3.33 3.24.79 4.58L8 10.75l-4.12 2.16.79-4.58L1.34 5.09l4.6-.67z"/>
                    </svg>
                    {project.stars}
                </span>
            {/if}
        </div>

        {#if project.description}
            <p class="description">{project.description}</p>
        {:else}
            <p class="description description--empty">No description provided.</p>
        {/if}

        {#if topics.length > 0}
            <ul class="topics">
                {#each topics as topic}
                    <li>{topic}</li>
                {/each}
            </ul>
        {/if}

        <div class="meta">
            {#if project.language}
                <span class="language">
                    <span class="language-dot" style="background-color: {getLanguageColor(project.language)}"></span>
                    {project.language}
                </span>
            {/if}
            {#if updated}
                <span class="updated">{updated}</span>
            {/if}
        </div>

        <a href={project.url} target="_blank" rel="noopener noreferrer" class="github-link">
            View on GitHub
            <span class="arrow" aria-hidden="true">→</span>
        </a>
    </div>
</article>

<style>
    .card {
        position: relative;
        display: flex;
        height: 100%;
        overflow: hidden;
        background: rgba(255, 255, 255, 0.04);
        border: 1px solid rgba(255, 255, 255, 0.09);
        border-radius: 16px;
        box-shadow: 0 20px 45px -30px rgba(0, 0, 0, 0.85);
        transition: transform 0.3s cubic-bezier(0.22, 1, 0.36, 1),
                    border-color 0.3s ease, background 0.3s ease, box-shadow 0.3s ease;
    }

    .card::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 3px;
        background: linear-gradient(120deg, #82aeff 0%, #a9a0ff 45%, #ff8d69 100%);
        opacity: 0;
        transition: opacity 0.3s ease;
    }

    .card:hover {
        transform: translateY(-6px);
        border-color: rgba(130, 174, 255, 0.35);
        background: rgba(255, 255, 255, 0.065);
        box-shadow: 0 30px 60px -32px rgba(110, 161, 255, 0.5);
    }

    .card:hover::before {
        opacity: 1;
    }

    .card-content {
        display: flex;
        flex-direction: column;
        width: 100%;
        padding: 1.6rem;
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        gap: 0.75rem;
    }

    .project-name {
        margin: 0;
        color: #f6f8ff;
        font-size: 1.12rem;
        font-weight: 700;
        letter-spacing: -0.015em;
        word-break: break-word;
    }

    .stars {
        display: inline-flex;
        align-items: center;
        gap: 0.28rem;
        flex-shrink: 0;
        color: #e3c96a;
        font-size: 0.82rem;
        font-weight: 600;
    }

    .description {
        color: rgba(233, 236, 243, 0.68);
        margin: 0.85rem 0 0;
        line-height: 1.6;
        font-size: 0.93rem;
    }

    .description--empty {
        color: rgba(233, 236, 243, 0.38);
        font-style: italic;
    }

    .topics {
        display: flex;
        flex-wrap: wrap;
        gap: 0.4rem;
        list-style: none;
        margin: 1rem 0 0;
        padding: 0;
    }

    .topics li {
        font-size: 0.74rem;
        font-weight: 600;
        letter-spacing: 0.02em;
        color: rgba(130, 174, 255, 0.9);
        background: rgba(110, 161, 255, 0.1);
        border: 1px solid rgba(110, 161, 255, 0.22);
        border-radius: 999px;
        padding: 0.15rem 0.6rem;
    }

    /* Pushes the meta row and link to the bottom so cards in a row align. */
    .meta {
        display: flex;
        align-items: center;
        flex-wrap: wrap;
        gap: 0.4rem 0.9rem;
        margin-top: auto;
        padding-top: 1.25rem;
        color: rgba(233, 236, 243, 0.5);
        font-size: 0.8rem;
    }

    .language {
        display: inline-flex;
        align-items: center;
        gap: 0.4rem;
        white-space: nowrap;
    }

    .language-dot {
        width: 9px;
        height: 9px;
        border-radius: 50%;
    }

    .updated {
        white-space: nowrap;
    }

    .github-link {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        gap: 0.45rem;
        width: 100%;
        background: rgba(255, 255, 255, 0.05);
        color: #e9ecf3;
        text-decoration: none;
        padding: 0.7rem;
        border-radius: 10px;
        border: 1px solid rgba(255, 255, 255, 0.1);
        margin-top: 1rem;
        box-sizing: border-box;
        font-weight: 600;
        font-size: 0.9rem;
        transition: background 0.2s ease, border-color 0.2s ease, color 0.2s ease;
    }

    .github-link:hover {
        background: rgba(110, 161, 255, 0.16);
        border-color: rgba(110, 161, 255, 0.5);
        color: #82aeff;
    }

    .arrow {
        transition: transform 0.2s ease;
    }

    .github-link:hover .arrow {
        transform: translateX(3px);
    }

    @media (prefers-reduced-motion: reduce) {
        .card,
        .arrow {
            transition: none;
        }

        .card:hover {
            transform: none;
        }
    }
</style>
