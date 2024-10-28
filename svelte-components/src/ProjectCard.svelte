<script>
    export let project;
</script>

<div class="card">
    <img src={project.image_url} alt={project.name} />
    <h2>{project.name}</h2>
    <p>{project.description}</p>
    <a href={project.url} target="_blank" rel="noopener noreferrer">View on GitHub</a>
</div>

<style>
    .card {
        background-color: #fff;
        border-radius: 8px;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        padding: 20px;
        transition: transform 0.3s ease;
    }

    .card:hover {
        transform: translateY(-5px);
    }

    img {
        width: 100%;
        height: 200px;
        object-fit: cover;
        border-radius: 8px;
    }

    h2 {
        margin-top: 10px;
    }

    a {
        display: inline-block;
        margin-top: 10px;
        padding: 8px 16px;
        background-color: #333;
        color: #fff;
        text-decoration: none;
        border-radius: 4px;
    }
</style>
