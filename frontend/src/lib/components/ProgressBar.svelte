<script>
  export let percent = 0;
  export let message = '';
  export let logs = [];
  export let showLogs = false;

  let logsExpanded = false;

  $: displayPercent = Math.min(100, Math.max(0, percent));
</script>

<div class="progress-container">
  <div class="progress-header">
    <span class="progress-message">{message}</span>
    <span class="progress-percent">{displayPercent}%</span>
  </div>

  <div class="progress-bar">
    <div class="progress-fill" style="width: {displayPercent}%"></div>
  </div>

  {#if showLogs && logs.length > 0}
    <button class="logs-toggle" on:click={() => logsExpanded = !logsExpanded}>
      <span class="toggle-icon">{logsExpanded ? '▼' : '▶'}</span>
      Details ({logs.length})
    </button>

    {#if logsExpanded}
      <div class="logs-container">
        {#each logs as log}
          <div class="log-line">{log}</div>
        {/each}
      </div>
    {/if}
  {/if}
</div>

<style>
  .progress-container {
    width: 100%;
    padding: 20px;
    background: var(--bg-secondary);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-sm);
  }

  .progress-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 12px;
    font-size: 0.9rem;
  }

  .progress-message {
    color: var(--text-primary);
    font-weight: 500;
  }

  .progress-percent {
    color: var(--accent-color);
    font-family: monospace;
    font-weight: 600;
  }

  .progress-bar {
    height: 8px;
    background: var(--bg-tertiary);
    border-radius: 4px;
    overflow: hidden;
  }

  .progress-fill {
    height: 100%;
    background: linear-gradient(90deg, var(--accent-color), var(--accent-hover));
    border-radius: 4px;
    transition: width var(--transition-smooth);
  }

  .logs-toggle {
    display: flex;
    align-items: center;
    gap: 6px;
    margin-top: 16px;
    padding: 6px 10px;
    background: var(--bg-tertiary);
    border: none;
    border-radius: var(--radius-sm);
    color: var(--text-secondary);
    font-size: 0.85rem;
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .logs-toggle:hover {
    background: var(--bg-hover);
    color: var(--text-primary);
  }

  .toggle-icon {
    font-size: 0.7rem;
    transition: transform var(--transition-fast);
  }

  .logs-container {
    margin-top: 12px;
    padding: 12px 16px;
    background: var(--bg-tertiary);
    border-radius: var(--radius-sm);
    max-height: 150px;
    overflow-y: auto;
    font-family: monospace;
    font-size: 0.8rem;
    text-align: left;
  }

  .log-line {
    color: var(--text-secondary);
    padding: 3px 0;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .log-line:last-child {
    color: var(--accent-color);
    font-weight: 500;
  }
</style>
