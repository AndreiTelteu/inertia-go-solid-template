import { createSignal, ErrorBoundary, onCleanup, onMount, Show, type JSX, type ParentProps } from 'solid-js'

/**
 * Detailed errors (message, stack trace) are shown only while developing with `./artisan dev`:
 * `import.meta.env.DEV` is false in every production build, so the stack never reaches real users.
 */
const detailed = import.meta.env.DEV

const messageOf = (error: unknown) => (typeof error === 'string' ? error : (error as Error)?.message || String(error))
const stackOf = (error: unknown) => (typeof error === 'object' && error ? ((error as Error).stack ?? '') : '')

export function ErrorPanel(props: { error: unknown; reset?: () => void }) {
  const copy = () => navigator.clipboard?.writeText(`${messageOf(props.error)}\n${stackOf(props.error)}`)
  return (
    <Show
      when={detailed}
      fallback={
        <div class="error-panel" role="alert">
          <h2>Internal Server Error</h2>
          <p>Something went wrong. Please try again.</p>
          <div class="actions">
            <Show when={props.reset}><button type="button" onClick={() => props.reset!()}>Try again</button></Show>
            <button type="button" class="secondary" onClick={() => location.reload()}>Reload page</button>
          </div>
        </div>
      }
    >
      <div class="error-panel error-panel-dev" role="alert">
        <span class="error-tag">Development only</span>
        <h2>Something went wrong on this page</h2>
        <p class="error-message">{messageOf(props.error)}</p>
        <Show when={stackOf(props.error)}><details open><summary>Stack trace</summary><pre>{stackOf(props.error)}</pre></details></Show>
        <div class="actions">
          <Show when={props.reset}><button type="button" onClick={() => props.reset!()}>Try again</button></Show>
          <button type="button" class="secondary" onClick={copy}>Copy error</button>
          <button type="button" class="secondary" onClick={() => location.reload()}>Reload page</button>
        </div>
      </div>
    </Show>
  )
}

/**
 * Catches render errors of its children. The default fallback is ErrorPanel; pass your own
 * `fallback` (same signature as Solid's ErrorBoundary) wherever a page needs a custom look.
 */
export function AppErrorBoundary(props: ParentProps<{ fallback?: (error: unknown, reset: () => void) => JSX.Element }>) {
  return <ErrorBoundary fallback={(error, reset) => props.fallback?.(error, reset) ?? <ErrorPanel error={error} reset={reset} />}>{props.children}</ErrorBoundary>
}

/** Errors outside rendering (rejected promises, event handlers) cannot reach an ErrorBoundary; show them as a banner while developing. */
export function AsyncErrorBanner() {
  const [error, setError] = createSignal<unknown>(null)
  const onError = (event: ErrorEvent) => setError(event.error ?? event.message)
  const onRejection = (event: PromiseRejectionEvent) => setError(event.reason)
  onMount(() => {
    if (!detailed) return
    window.addEventListener('error', onError)
    window.addEventListener('unhandledrejection', onRejection)
  })
  onCleanup(() => {
    window.removeEventListener('error', onError)
    window.removeEventListener('unhandledrejection', onRejection)
  })
  return (
    <Show when={error()}>
      <div class="error-banner">
        <ErrorPanel error={error()} />
        <button type="button" class="secondary" onClick={() => setError(null)}>Dismiss</button>
      </div>
    </Show>
  )
}
