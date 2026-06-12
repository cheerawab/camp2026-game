const MATCH_CODE_PATTERN = /[^A-Z0-9]/g

export function normalizeMatchCode(value: string) {
  return value.trim().toUpperCase().replace(MATCH_CODE_PATTERN, "")
}
