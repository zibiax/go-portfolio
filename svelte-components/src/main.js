import ProjectCard from './ProjectCard.svelte';

// This will handle the HTMX response and create Svelte components
document.addEventListener('htmx:afterSwap', function(event) {
    if (event.detail.target.id === 'projects') {
        const projects = JSON.parse(event.detail.xhr.response);
        const container = document.getElementById('projects');
        container.innerHTML = '';

        projects.forEach(project => {
            const div = document.createElement('div');
            new ProjectCard({
                target: div,
                props: { project }
            });
            container.appendChild(div);
        });
    }
});
