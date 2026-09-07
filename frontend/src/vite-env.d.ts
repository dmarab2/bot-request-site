/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_BACKEND_ROOT: string
  readonly VITE_BACKEND_GET_OPEN_REQUESTS: string
  readonly VITE_BACKEND_REQUESTS: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}