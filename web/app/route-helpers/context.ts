import { createContext } from "react-router"

import type { RootLoaderData } from "./types"

export type RootDataProvider = () => Promise<RootLoaderData>

export const rootDataContext = createContext<RootDataProvider | null>(null)
