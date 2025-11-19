export {}

declare global {
  interface ImportMetaEnv {
    readonly VITE_API_BASE_URL: string
    readonly VITE_MAX_FILE_SIZE: string
    readonly VITE_SUPPORTED_FORMATS: string
  }

  interface ImportMeta {
    readonly env: ImportMetaEnv
  }
}
