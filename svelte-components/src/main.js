import { mount } from 'svelte';
import ProjectCard from './ProjectCard.svelte';

const grid = document.getElementById('projects');
const status = document.getElementById('projects-status');
const loadMore = document.getElementById('load-more');

let nextPage = 1;
let loading = false;

function setStatus(message, isError = false) {
    if (!status) return;
    status.textContent = message || '';
    status.hidden = !message;
    status.classList.toggle('projects-status--error', isError);
}

function renderProjects(projects) {
    projects.forEach((project, index) => {
        const div = document.createElement('div');
        div.className = 'project-card-wrap';
        // Stagger the entrance so a freshly loaded page fades in row by row.
        div.style.setProperty('--card-index', index);
        mount(ProjectCard, {
            target: div,
            props: { project }
        });
        grid.appendChild(div);
    });
}

async function loadPage() {
    if (loading) return;
    loading = true;

    if (loadMore) {
        loadMore.disabled = true;
        loadMore.textContent = 'Loading…';
    }
    setStatus(nextPage === 1 ? 'Loading projects…' : '');

    try {
        const response = await fetch(`/projects?page=${nextPage}`);
        if (!response.ok) {
            throw new Error(`Request failed with status ${response.status}`);
        }

        const { projects = [], hasMore = false } = await response.json();

        if (projects.length === 0 && nextPage === 1) {
            setStatus('No projects found.');
        } else {
            setStatus('');
            renderProjects(projects);
        }

        nextPage += 1;
        grid.setAttribute('aria-busy', 'false');
        if (loadMore) {
            loadMore.hidden = !hasMore;
            loadMore.disabled = false;
            loadMore.textContent = 'Load more projects';
        }
    } catch (error) {
        console.error('Error loading projects:', error);
        setStatus(`Could not load projects: ${error.message}`, true);
        grid.setAttribute('aria-busy', 'false');
        if (loadMore) {
            loadMore.hidden = false;
            loadMore.disabled = false;
            loadMore.textContent = 'Try again';
        }
    } finally {
        loading = false;
    }
}

if (grid) {
    loadMore?.addEventListener('click', loadPage);
    loadPage();
}
