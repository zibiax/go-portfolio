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
        position: relative;
        overflow: hidden;
        background: rgba(255, 255, 255, 0.045);
        border-radius: 14px;
        padding: 1.5rem;
        box-shadow: 0 20px 40px -28px rgba(0, 0, 0, 0.7);
        transition: transform 0.25s ease, border-color 0.25s ease, background 0.25s ease, box-shadow 0.25s ease;
        border: 1px solid rgba(255, 255, 255, 0.09);
    }

    .card::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        height: 3px;
        background: linear-gradient(135deg, #82aeff 0%, #ff8d69 100%);
        opacity: 0;
        transition: opacity 0.25s ease;
    }

    .card:hover {
        transform: translateY(-6px);
        border-color: rgba(130, 174, 255, 0.35);
        background: rgba(255, 255, 255, 0.065);
        box-shadow: 0 30px 60px -30px rgba(110, 161, 255, 0.45);
    }

    .card:hover::before {
        opacity: 1;
    }

    .header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        gap: 0.75rem;
        margin-bottom: 1rem;
    }

    .project-name {
        margin: 0;
        color: #f4f7ff;
        font-size: 1.15rem;
        font-weight: 700;
        letter-spacing: -0.01em;
    }

    .language {
        display: flex;
        align-items: center;
        gap: 0.4rem;
        color: rgba(233, 236, 243, 0.5);
        font-size: 0.82rem;
        white-space: nowrap;
    }

    .language-dot {
        width: 9px;
        height: 9px;
        border-radius: 50%;
    }

    .description {
        color: rgba(233, 236, 243, 0.68);
        margin: 1rem 0;
        line-height: 1.55;
        font-size: 0.95rem;
    }

    .github-link {
        display: inline-block;
        width: 100%;
        background: rgba(255, 255, 255, 0.06);
        color: #e9ecf3;
        text-decoration: none;
        padding: 0.7rem;
        border-radius: 8px;
        border: 1px solid rgba(255, 255, 255, 0.1);
        text-align: center;
        margin-top: 1rem;
        box-sizing: border-box;
        font-weight: 600;
        font-size: 0.9rem;
        transition: background 0.2s ease, border-color 0.2s ease;
    }

    .github-link:hover {
        background: rgba(110, 161, 255, 0.16);
        border-color: rgba(110, 161, 255, 0.5);
        color: #82aeff;
    }
</style>
