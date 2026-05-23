"use client"

import * as React from "react"

export type TopicEnvState = {
  currentCategoryId: number
  currentRootCategoryId: number
  currentChildCategoryId: number
}

type TopicEnvContextValue = TopicEnvState & {
  setCurrentCategoryId: (categoryId: number) => void
  setNodeContext: (state: Partial<TopicEnvState>) => void
}

const defaultTopicEnv: TopicEnvState = {
  currentCategoryId: 0,
  currentRootCategoryId: 0,
  currentChildCategoryId: 0,
}

const TopicEnvContext = React.createContext<TopicEnvContextValue | null>(null)

export function TopicEnvProvider({
  initialState = defaultTopicEnv,
  children,
}: {
  initialState?: Partial<TopicEnvState>
  children: React.ReactNode
}) {
  const [state, setState] = React.useState<TopicEnvState>({
    ...defaultTopicEnv,
    ...initialState,
  })

  const value = React.useMemo<TopicEnvContextValue>(
    () => ({
      ...state,
      setCurrentCategoryId: (currentCategoryId) => setState((previous) => ({ ...previous, currentCategoryId })),
      setNodeContext: (nextState) => setState((previous) => ({ ...previous, ...nextState })),
    }),
    [state],
  )

  return <TopicEnvContext.Provider value={value}>{children}</TopicEnvContext.Provider>
}

export function useTopicEnv() {
  const value = React.useContext(TopicEnvContext)

  if (!value) {
    throw new Error("useTopicEnv must be used within TopicEnvProvider")
  }

  return value
}
