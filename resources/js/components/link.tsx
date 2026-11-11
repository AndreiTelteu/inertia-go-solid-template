import { Link as InertiaLink } from 'inertia-adapter-solid'
import { type ComponentProps, mergeProps } from 'solid-js'

/**
 * Inertia Link that preloads the target page when the pointer hovers over it (or it is focused),
 * so the click feels instant. Use this instead of importing Link from 'inertia-adapter-solid'.
 * Pass `prefetch={false}` to opt out, or another strategy such as `prefetch="click"`.
 */
export function Link(props: ComponentProps<typeof InertiaLink>) {
  return <InertiaLink {...mergeProps({ prefetch: 'hover' as const }, props)} />
}
