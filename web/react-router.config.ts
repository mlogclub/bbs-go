import type { Config } from "@react-router/dev/config"

const spa = process.env.BBSGO_WEB_SPA === "true"

export default {
  ssr: !spa,
  future: {
    v8_middleware: true,
  },
} satisfies Config
