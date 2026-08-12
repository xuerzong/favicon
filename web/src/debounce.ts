export const debounce = <T extends (...args: any[]) => any>(
  fn: T,
  delay: number,
  immediate = false
) => {
  let timer: ReturnType<typeof setTimeout> | null = null

  return (...args: Parameters<T>): ReturnType<T> | undefined => {
    if (timer) clearTimeout(timer)

    if (immediate && !timer) {
      return fn(...args)
    }

    timer = setTimeout(() => {
      timer = null
      if (!immediate) {
        fn(...args)
      }
    }, delay)
  }
}
