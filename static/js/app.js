import ProjectCard from './ProjectCard.svelte';

document.addEventListener('htmx:afterSwap', function(event) {
    if (event.detail.target.id === 'projects') {
        console.log('HTMX afterSwap event triggered');
        try {
            const response = event.detail.xhr.response;
            console.log('Raw response:', response);
            
            const projects = JSON.parse(response);
            console.log('Parsed projects:', projects);
            
            const container = document.getElementById('projects');
            container.innerHTML = '';

            if (projects.length === 0) {
                console.log('No projects found');
                container.innerHTML = '<p>No projects found</p>';
                return;
            }

            projects.forEach(project => {
                console.log('Creating component for project:', project);
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
