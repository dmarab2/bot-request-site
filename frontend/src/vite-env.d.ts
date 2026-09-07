/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly BACKEND_ROOT: string
  readonly BACKEND_GET_OPEN_REQUESTS: string
  readonly BACKEND_REQUESTS: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}