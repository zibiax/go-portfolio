<script>
    export let project;

    const languageColors = {
        'javascript': '#f1e05a',
        'typescript': '#2b7489',
        'python': '#3572A5',
        'go': '#00ADD8',
        'java': '#b07219',
        'html': '#e34c26',
        'css': '#563d7c',
        'rust': '#dea584',
        'default': '#6e7681'
    };

    function getLanguageColor(lang) {
        if (!lang) return languageColors.default;
        const normalizedLang = lang.toLowerCase();
        return languageColors[normalizedLang] || languageColors.default;
    }
</script>

<div class="card">
    <div class="card-content">
        <div class="header">
            <h2 class="project-name">{project.name}</h2>
            {#if project.language}
                <div class="language">
                    <span class="language-dot" style="background-color: {getLanguageColor(project.language)}"></span>
                    {project.language}
                </div>
            {/if}
        </div>

        {#if project.description}
            <p class="description">{project.description}</p>
        {/if}

        <a href={project.url} target="_blank" rel="noopener noreferrer" class="github-link">
            View on GitHub
        </a>
    </div>
</div>

<style>
    .card {
        background: #2d2d2d;
        border-radius: 8px;
        padding: 1.5rem;
        box-shadow: 0 2px 4px rgba(0, 0, 0, 0.2);
        transition: transform 0.2s ease;
        border: 1px solid #404040;
    }

    .card:hover {
        transform: translateY(-4px);
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1rem;
    }

    .project-name {
        margin: 0;
        color: #ffffff;
        font-size: 1.25rem;
    }

    .language {
        display: flex;
        align-items: center;
        gap: 0.5rem;
        color: #999;
        font-size: 0.875rem;
    }

    .language-dot {
        width: 12px;
        height: 12px;
        border-radius: 50%;
    }

    .description {
        color: #b3b3b3;
        margin: 1rem 0;
        line-height: 1.5;
    }

    .github-link {
        display: inline-block;
        width: 100%;
        background: #404040;
        color: #ffffff;
        text-decoration: none;
        padding: 0.75rem;
        border-radius: 6px;
        text-align: center;
        margin-top: 1rem;
        box-sizing: border-box;
    }

    .github-link:hover {
        background: #505050;
    }
</style>
