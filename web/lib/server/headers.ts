export async function cookies() {
  return {
    get: (_name: string): { value: string } | undefined => undefined,
  }
}
