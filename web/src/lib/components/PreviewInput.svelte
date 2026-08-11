<script lang="ts">
    import { debounce } from "../utils/debounce";

    const API_BASE_URL = "http://localhost:8080/";

    let url = $state("");
    let debounceUrl = $state("");
    let urlInputting = $state(false);
    let hasFocused = $state(false);

    const handleUrlChange = debounce((nextUrl) => {
        debounceUrl = nextUrl;
        urlInputting = false;
    }, 1000);

    function handleFocus(e: FocusEvent) {
        hasFocused = true;
    }

    function handleBlur() {
        hasFocused = false;
    }

    function handleInput() {
        urlInputting = true;
        handleUrlChange(url);
    }
</script>

<label class="input-wrapper" data-focused={hasFocused}>
    <img
        class="input-image"
        src={`${API_BASE_URL}${debounceUrl || "example.com"}`}
        alt="favicon"
    />
    <span class="input-prefix">{API_BASE_URL}</span>
    <input
        oninput={handleInput}
        bind:value={url}
        onfocus={handleFocus}
        onblur={handleBlur}
        class="input"
        type="text"
        placeholder="example.com"
    />
</label>

<div>
    {#if urlInputting}
        <div>加载中</div>
    {/if}
</div>

<style>
    .input-wrapper {
        display: flex;
        align-items: center;
        border: 1px solid var(--border-color);
        border-radius: var(--border-radius);
        font-size: 1rem;
        padding: 1rem;
    }

    .input-wrapper[data-focused="true"] {
        border-color: rgb(var(--color-primary))
    }

    .input-image {
        width: 1.5rem;
        height: 1.5rem;
        margin-right: 0.5rem;
    }

    .input-prefix {
        flex-shrink: 0;
        white-space: nowrap;
    }

    .input {
        flex: 1;
        border: none;
        outline: none;
        font-size: 1rem;
    }
</style>
