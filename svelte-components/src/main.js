import ProjectCard from './ProjectCard.svelte';

document.addEventListener('htmx:afterSwap', function(event) {
    if (event.detail.target.id === 'projects') {
        try {
            const projects = JSON.parse(event.detail.xhr.response);
            const container = document.getElementById('projects');
            container.innerHTML = '';

            if (projects.length === 0) {
                container.innerHTML = '<p>No projects found</p>';
                return;
            }

            projects.forEach(project => {
                const div = document.createElement('div');
                new ProjectCard({
                    target: div,
                    props: { project }
                });
                container.appendChild(div);
            });
        } catch (error) {
            console.error('Error processing projects:', error);
            const container = document.getElementById('projects');
            container.innerHTML = `<p>Error loading projects: ${error.message}</p>`;
        }
    }
});
