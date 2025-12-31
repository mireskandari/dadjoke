<script>
  import { createEventDispatcher, onMount, onDestroy } from 'svelte';
  import { EventsOn, EventsOff } from '../../wailsjs/runtime/runtime.js';
  import { SelectPDFFile, CompressPDF, SaveFile, OpenFile } from '../../wailsjs/go/main/App.js';
  import FileDropZone from './components/FileDropZone.svelte';
  import ProgressBar from './components/ProgressBar.svelte';

  const dispatch = createEventDispatcher();

  // State
  let file = null;
  let preset = 'printer';
  let isProcessing = false;
  let progress = 0;
  let progressMessage = '';
  let logs = [];
  let result = null;
  let error = null;
  let savedPath = null;

  const presets = [
    { id: 'screen', name: 'Screen', description: 'Smallest file, lower quality' },
    { id: 'ebook', name: 'eBook', description: 'Good for digital reading' },
    { id: 'printer', name: 'Printer', description: 'Balanced quality and size' },
    { id: 'prepress', name: 'Prepress', description: 'High quality for printing' },
    { id: 'default', name: 'Default', description: 'Minimal compression' },
  ];

  // Event handlers
  function handleProgress(data) {
    progress = data.percent;
    progressMessage = data.message;
  }

  function handleLog(line) {
    logs = [...logs, line];
  }

  onMount(() => {
    EventsOn('compress:progress', handleProgress);
    EventsOn('compress:log', handleLog);
  });

  onDestroy(() => {
    EventsOff('compress:progress');
    EventsOff('compress:log');
  });

  // Actions
  async function handleBrowse() {
    try {
      const selected = await SelectPDFFile();
      if (selected) {
        file = selected;
        error = null;
        result = null;
        savedPath = null;
      }
    } catch (e) {
      error = e.message || 'Failed to select file';
    }
  }

  async function handleCompress() {
    if (!file) return;

    isProcessing = true;
    progress = 0;
    progressMessage = 'Starting...';
    logs = [];
    error = null;
    result = null;

    try {
      result = await CompressPDF(file.path, preset);

      if (result.success) {
        // Show save dialog
        const suggestedName = file.name.replace('.pdf', '_compressed.pdf');
        savedPath = await SaveFile(result.outputPath, suggestedName);
      } else {
        error = result.error || 'Compression failed';
      }
    } catch (e) {
      error = e.message || 'Compression failed';
    } finally {
      isProcessing = false;
    }
  }

  async function handleOpen() {
    if (savedPath) {
      try {
        await OpenFile(savedPath);
      } catch (e) {
        error = 'Failed to open file';
      }
    }
  }

  function handleReset() {
    file = null;
    result = null;
    savedPath = null;
    error = null;
    progress = 0;
    logs = [];
  }

  function goBack() {
    dispatch('navigate', { view: 'home' });
  }

  function formatSize(bytes) {
    if (bytes < 1024) return bytes + ' B';
    if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB';
    return (bytes / (1024 * 1024)).toFixed(1) + ' MB';
  }
</script>

<div class="compress">
  <header class="header">
    <button class="back-btn" on:click={goBack}>
      <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
        <path d="M19 12H5M12 19l-7-7 7-7"/>
      </svg>
    </button>
    <h1>Compress PDF</h1>
  </header>

  <div class="content">
    {#if !file}
      <!-- File selection -->
      <FileDropZone on:browse={handleBrowse} on:files={(e) => console.log('dropped', e.detail)}>
        Drag & drop PDF here
      </FileDropZone>

    {:else if isProcessing}
      <!-- Processing -->
      <div class="file-info">
        <span class="file-name">{file.name}</span>
        <span class="file-size">{file.sizeText}</span>
      </div>

      <ProgressBar
        percent={progress}
        message={progressMessage}
        {logs}
        showLogs={true}
      />

    {:else if result && result.success}
      <!-- Success -->
      <div class="success-card">
        <div class="success-icon">✓</div>
        <h2>Compression Complete</h2>

        <div class="stats">
          <div class="stat">
            <span class="stat-label">Original</span>
            <span class="stat-value">{formatSize(result.originalSize)}</span>
          </div>
          <div class="stat-arrow">→</div>
          <div class="stat">
            <span class="stat-label">Compressed</span>
            <span class="stat-value">{formatSize(result.compressedSize)}</span>
          </div>
        </div>

        <div class="savings" class:negative={result.savingsPercent < 0}>
          {#if result.savingsPercent >= 0}
            Saved {result.savingsPercent}%
          {:else}
            File grew by {Math.abs(result.savingsPercent)}%
          {/if}
        </div>

        {#if savedPath}
          <div class="actions">
            <button class="btn primary" on:click={handleOpen}>Open</button>
            <button class="btn secondary" on:click={handleReset}>Compress Another</button>
          </div>
        {:else}
          <p class="save-cancelled">Save cancelled. <button class="link-btn" on:click={() => SaveFile(result.outputPath, file.name.replace('.pdf', '_compressed.pdf')).then(p => savedPath = p)}>Try again</button></p>
        {/if}
      </div>

    {:else}
      <!-- File loaded, ready to compress -->
      <div class="file-card">
        <div class="file-info">
          <span class="file-name">{file.name}</span>
          <span class="file-size">{file.sizeText}</span>
        </div>
        <button class="remove-btn" on:click={handleReset}>✕</button>
      </div>

      <div class="presets">
        <h3>Select Quality</h3>
        {#each presets as p}
          <label class="preset-option" class:selected={preset === p.id}>
            <input type="radio" bind:group={preset} value={p.id} />
            <span class="preset-name">{p.name}</span>
            <span class="preset-desc">{p.description}</span>
          </label>
        {/each}
      </div>

      <button class="btn primary compress-btn" on:click={handleCompress}>
        Compress
      </button>
    {/if}

    {#if error}
      <div class="error">{error}</div>
    {/if}
  </div>
</div>

<style>
  .compress {
    max-width: 520px;
    margin: 0 auto;
    padding: 24px;
  }

  .header {
    display: flex;
    align-items: center;
    gap: 14px;
    margin-bottom: 28px;
  }

  .header h1 {
    font-size: 1.5rem;
    font-weight: 600;
    margin: 0;
    color: var(--text-primary);
  }

  .back-btn {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 40px;
    height: 40px;
    background: var(--bg-secondary);
    border: none;
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-sm);
    color: var(--text-secondary);
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .back-btn:hover {
    background: var(--bg-hover);
    color: var(--accent-color);
    box-shadow: var(--shadow-md);
  }

  .content {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .file-card {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 18px 20px;
    background: var(--bg-secondary);
    border-radius: var(--radius-md);
    box-shadow: var(--shadow-sm);
  }

  .file-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .file-name {
    font-weight: 500;
    color: var(--text-primary);
  }

  .file-size {
    font-size: 0.85rem;
    color: var(--text-secondary);
  }

  .remove-btn {
    width: 32px;
    height: 32px;
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--bg-tertiary);
    border: none;
    border-radius: var(--radius-sm);
    color: var(--text-muted);
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .remove-btn:hover {
    background: var(--error-light);
    color: var(--error-color);
  }

  .presets {
    background: var(--bg-secondary);
    border-radius: var(--radius-md);
    padding: 20px;
    box-shadow: var(--shadow-sm);
  }

  .presets h3 {
    margin: 0 0 14px;
    font-size: 0.9rem;
    font-weight: 500;
    color: var(--text-secondary);
  }

  .preset-option {
    display: flex;
    align-items: center;
    gap: 14px;
    padding: 14px 16px;
    margin: 6px 0;
    background: var(--bg-tertiary);
    border: 2px solid transparent;
    border-radius: var(--radius-sm);
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .preset-option:hover {
    background: var(--bg-hover);
    border-color: var(--border-color);
  }

  .preset-option.selected {
    background: var(--accent-light);
    border-color: var(--accent-color);
  }

  .preset-option input {
    display: none;
  }

  .preset-name {
    font-weight: 600;
    min-width: 80px;
    color: var(--text-primary);
  }

  .preset-desc {
    font-size: 0.85rem;
    color: var(--text-secondary);
  }

  .preset-option.selected .preset-name {
    color: var(--accent-color);
  }

  .preset-option.selected .preset-desc {
    color: var(--text-secondary);
  }

  .btn {
    padding: 14px 28px;
    border: none;
    border-radius: var(--radius-sm);
    font-size: 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: all var(--transition-fast);
  }

  .btn.primary {
    background: var(--accent-color);
    color: white;
  }

  .btn.primary:hover {
    background: var(--accent-hover);
    transform: translateY(-1px);
    box-shadow: var(--shadow-md);
  }

  .btn.secondary {
    background: var(--bg-secondary);
    color: var(--text-primary);
    border: 1px solid var(--border-color);
  }

  .btn.secondary:hover {
    background: var(--bg-hover);
  }

  .compress-btn {
    width: 100%;
    margin-top: 4px;
  }

  .success-card {
    text-align: center;
    padding: 36px;
    background: var(--bg-secondary);
    border-radius: var(--radius-lg);
    box-shadow: var(--shadow-sm);
  }

  .success-icon {
    width: 56px;
    height: 56px;
    line-height: 56px;
    margin: 0 auto 20px;
    background: var(--success-light);
    border-radius: 50%;
    font-size: 1.6rem;
    color: var(--success-color);
  }

  .success-card h2 {
    margin: 0 0 24px;
    font-size: 1.4rem;
    color: var(--text-primary);
  }

  .stats {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 20px;
    margin-bottom: 20px;
    padding: 16px;
    background: var(--bg-tertiary);
    border-radius: var(--radius-sm);
  }

  .stat {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .stat-label {
    font-size: 0.8rem;
    color: var(--text-secondary);
  }

  .stat-value {
    font-size: 1.2rem;
    font-weight: 600;
    color: var(--text-primary);
  }

  .stat-arrow {
    color: var(--text-muted);
    font-size: 1.2rem;
  }

  .savings {
    font-size: 1.1rem;
    font-weight: 600;
    color: var(--success-color);
    margin-bottom: 24px;
  }

  .savings.negative {
    color: var(--error-color);
  }

  .actions {
    display: flex;
    gap: 12px;
    justify-content: center;
  }

  .save-cancelled {
    color: var(--text-secondary);
    font-size: 0.9rem;
  }

  .link-btn {
    background: none;
    border: none;
    color: var(--accent-color);
    cursor: pointer;
    text-decoration: underline;
    font-size: inherit;
  }

  .link-btn:hover {
    color: var(--accent-hover);
  }

  .error {
    padding: 14px 18px;
    background: var(--error-light);
    border: 1px solid var(--error-color);
    border-radius: var(--radius-sm);
    color: var(--error-color);
    text-align: center;
    font-size: 0.95rem;
  }
</style>
