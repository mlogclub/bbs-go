import { defineConfig, globalIgnores } from "eslint/config"
const eslintConfig = defineConfig([
  globalIgnores([
    "out/**",
    "dist/**",
    "build/**",
  ]),
])

export default eslintConfig
