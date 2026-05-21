"use client"

import * as React from "react"

export type TopicEnvState = {
  currentNodeId: number
  currentRootNodeId: number
  currentChildNodeId: number
}

type TopicEnvContextValue = TopicEnvState & {
  setCurrentNodeId: (nodeId: number) => void
  setNodeContext: (state: Partial<TopicEnvState>) => void
}

const defaultTopicEnv: TopicEnvState = {
  currentNodeId: 0,
  currentRootNodeId: 0,
  currentChildNodeId: 0,
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
      setCurrentNodeId: (currentNodeId) => setState((previous) => ({ ...previous, currentNodeId })),
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
