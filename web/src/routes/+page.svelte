<script lang="ts">
    import PreviewInput from '$lib/components/PreviewInput.svelte';
    import GithubIcon from '$lib/assets/github.svg';
</script>

<section class="hero">
    <h2 class="hero-title">Favicon API</h2>
    <p class="hero-description">
        A simple API service to get favicons from any website
    </p>
    <div class="hero-actions">
        <a class="btn btn-primary" href="/#"> Get Started </a>
        <a class="btn" href="https://github.com/xuerzong/favicon">
            <img class="btn-icon" src={GithubIcon} alt="github" />
            View on Github
        </a>
    </div>
</section>

<section class="previwer">
    <div class="previewer-card">
        <PreviewInput />
    </div>
</section>

<style>
    .hero {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
    }
    .hero-title {
        font-size: 4rem;
        font-weight: bold;
        color: rgb(var(--color-primary));
        margin: 0;
        margin-bottom: 1rem;
    }
    .hero-actions {
        display: flex;
        align-items: center;
        gap: 1rem;
    }
    .btn {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        padding: 0.75rem 1rem;
        border-radius: var(--border-radius);
        border: 1px solid var(--border-color);
        text-decoration: none;
        color: #333;
    }

    .btn-primary {
        background-color: rgb(var(--color-primary));
        color: #fff;
        border-color: rgb(var(--color-primary));
    }

    .btn-icon {
        width: 1.5rem;
        height: 1.5rem;
    }

    .previwer {
        margin-top: 2rem;
    }

    .previewer-card {
        max-width: 60rem;
        margin: 0 auto;
    }
</style>
