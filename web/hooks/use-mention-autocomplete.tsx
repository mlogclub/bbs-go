"use client"

import * as React from "react"
import { createPortal } from "react-dom"
import { searchUsers } from "@/lib/api/users"
import type { SearchUser } from "@/lib/api/types"
import { cn } from "@/lib/utils"

type MentionUser = {
  username: string
  nickname: string
  avatar: string
}

type MentionState = {
  active: boolean
  query: string
  index: number
  users: MentionUser[]
  position: { top: number; left: number }
}

export function useMentionAutocomplete(
  textareaRef: React.RefObject<HTMLTextAreaElement | null>,
  onInsert: (before: string, after: string) => void
) {
  const [state, setState] = React.useState<MentionState>({
    active: false,
    query: "",
    index: 0,
    users: [],
    position: { top: 0, left: 0 },
  })
  const stateRef = React.useRef(state)
  stateRef.current = state

  const close = React.useCallback(() => {
    setState((s) => ({ ...s, active: false, query: "", index: 0, users: [] }))
  }, [])

  // Search users when query changes
  React.useEffect(() => {
    if (!state.active) return
    const query = state.query
    if (!query) {
      // Show suggestions even with empty query (recent/active users would go here)
      searchUsers({ keyword: "" })
        .then((result) => {
          const results = result?.results || []
          setState((s) => {
            if (!s.active) return s
            return {
              ...s,
              users: results
                .map((u: SearchUser) => {
                  const user = (u as any).user || u
                  const uname = (user as any).username || u.username || ""
                  const nname = (user as any).nickname || u.nickname || ""
                  return {
                    username: uname.replace(/<[^>]*>/g, ""),
                    nickname: nname.replace(/<[^>]*>/g, ""),
                    avatar: (user as any).smallAvatar || (user as any).avatar || "",
                  }
                })
                .filter((item: MentionUser) => item.username)
                .slice(0, 8),
              index: 0,
            }
          })
        })
        .catch(() => close())
      return
    }

    const timer = setTimeout(() => {
      searchUsers({ keyword: query })
        .then((result) => {
          const results = result?.results || []
          setState((s) => {
            if (!s.active || s.query !== query) return s
            return {
              ...s,
              users: results
                .map((u: SearchUser) => {
                  const user = (u as any).user || u
                  const uname = (user as any).username || u.username || ""
                  const nname = (user as any).nickname || u.nickname || ""
                  return {
                    username: uname.replace(/<[^>]*>/g, ""),
                    nickname: nname.replace(/<[^>]*>/g, ""),
                    avatar: (user as any).smallAvatar || (user as any).avatar || "",
                  }
                })
                .filter((item: MentionUser) => item.username)
                .slice(0, 8),
              index: 0,
            }
          })
        })
        .catch(() => close())
    }, 150)
    return () => clearTimeout(timer)
  }, [state.active, state.query, close])

  function getCursorCoords(textarea: HTMLTextAreaElement, cursorPos: number) {
    // Create a mirror div to measure text position
    const style = window.getComputedStyle(textarea)
    const mirror = document.createElement("div")
    mirror.style.whiteSpace = "pre-wrap"
    mirror.style.wordWrap = "break-word"
    mirror.style.overflowWrap = "break-word"
    mirror.style.position = "absolute"
    mirror.style.visibility = "hidden"
    mirror.style.width = style.width
    mirror.style.padding = style.padding
    mirror.style.font = style.font
    mirror.style.lineHeight = style.lineHeight
    mirror.style.letterSpacing = style.letterSpacing
    mirror.style.border = style.border
    mirror.style.boxSizing = style.boxSizing
    mirror.textContent = textarea.value.substring(0, cursorPos) + "\u00A0"

    document.body.appendChild(mirror)
    const span = document.createElement("span")
    span.textContent = "\u00A0"
    mirror.appendChild(span)
    const rect = span.getBoundingClientRect()
    const textareaRect = textarea.getBoundingClientRect()
    document.body.removeChild(mirror)

    return {
      top: textareaRect.top + rect.top - textareaRect.top + 20,
      left: textareaRect.left + rect.left - textareaRect.left,
    }
  }

  function checkForMention(textarea: HTMLTextAreaElement) {
    const pos = textarea.selectionStart
    const value = textarea.value
    // Find the @ that starts the current mention
    let atIndex = -1
    for (let i = pos - 1; i >= 0; i--) {
      const ch = value[i]
      if (ch === "@") {
        // Check if @ is at start or preceded by whitespace/punctuation
        if (i === 0 || /[\s,.;:!?/>)\]}]/.test(value[i - 1])) {
          atIndex = i
          break
        }
      }
      if (ch === " " || ch === "\n") break
    }

    if (atIndex >= 0) {
      const query = value.substring(atIndex + 1, pos)
      // Only trigger if query is valid (starts with letter/digit or empty)
      if (query === "" || /^[\p{L}\p{N}_\-]*$/u.test(query)) {
        const position = getCursorCoords(textarea, pos)
        setState((s) => ({
          ...s,
          active: true,
          query,
          index: 0,
          users: [],
          position,
        }))
        return
      }
    }
    close()
  }

  function selectUser(user: MentionUser) {
    const textarea = textareaRef.current
    if (!textarea) return

    const pos = textarea.selectionStart
    const value = textarea.value
    // Find the @ that starts the mention
    let atIndex = pos - 1
    while (atIndex >= 0 && value[atIndex] !== "@") atIndex--
    if (atIndex < 0) return

    const before = value.substring(0, atIndex)
    const after = value.substring(pos)
    const mention = `@${user.username} `

    onInsert(before + mention, after)
    close()

    // Set cursor after the inserted mention
    setTimeout(() => {
      const newPos = (before + mention).length
      textarea.setSelectionRange(newPos, newPos)
      textarea.focus()
    }, 0)
  }

  function onKeyDown(e: React.KeyboardEvent<HTMLTextAreaElement>) {
    if (!state.active || !state.users.length) return false

    if (e.key === "ArrowDown" || (e.key === "Tab" && !e.shiftKey)) {
      e.preventDefault()
      setState((s) => ({
        ...s,
        index: (s.index + 1) % s.users.length,
      }))
      return true
    }
    if (e.key === "ArrowUp" || (e.key === "Tab" && e.shiftKey)) {
      e.preventDefault()
      setState((s) => ({
        ...s,
        index: (s.index - 1 + s.users.length) % s.users.length,
      }))
      return true
    }
    if (e.key === "Enter" || e.key === " ") {
      e.preventDefault()
      selectUser(state.users[state.index])
      return true
    }
    if (e.key === "Escape") {
      e.preventDefault()
      close()
      return true
    }
    return false
  }

  function onInput(e: React.ChangeEvent<HTMLTextAreaElement>) {
    checkForMention(e.currentTarget)
  }

  const popup = state.active ? (
    <MentionPopup
      users={state.users}
      selectedIndex={state.index}
      position={state.position}
      onSelect={selectUser}
      onClose={close}
    />
  ) : null

  return { popup, onKeyDown, onInput, close }
}

function MentionPopup({
  users,
  selectedIndex,
  position,
  onSelect,
  onClose,
}: {
  users: MentionUser[]
  selectedIndex: number
  position: { top: number; left: number }
  onSelect: (user: MentionUser) => void
  onClose: () => void
}) {
  const ref = React.useRef<HTMLDivElement>(null)

  // Adjust position to stay within viewport
  React.useEffect(() => {
    if (!ref.current) return
    const rect = ref.current.getBoundingClientRect()
    let { top, left } = position
    if (left + rect.width + 8 > window.innerWidth) {
      left = window.innerWidth - rect.width - 8
    }
    if (left < 8) left = 8
    if (top + rect.height + 8 > window.innerHeight) {
      top = position.top - rect.height - 24
    }
    if (top < 8) top = 8
    ref.current.style.top = `${top}px`
    ref.current.style.left = `${left}px`
  }, [position, users])

  if (!users.length) {
    return createPortal(
      <div ref={ref} className="mention-popup" style={{ position: "fixed", zIndex: 9999 }}>
        <div className="mention-menu">
          <div className="mention-no-results">No users found</div>
        </div>
      </div>,
      document.body
    )
  }

  return createPortal(
    <div ref={ref} className="mention-popup" style={{ position: "fixed", zIndex: 9999 }}>
      <div className="mention-menu">
        <div className="mention-items-container">
          {users.map((user, i) => (
            <button
              key={user.username}
              type="button"
              className={cn("mention-item", i === selectedIndex && "is-selected")}
              onMouseDown={(e) => {
                e.preventDefault()
                onSelect(user)
              }}
              onMouseEnter={() => {
                // Update selected index on hover
                const buttons = ref.current?.querySelectorAll(".mention-item")
                if (buttons) {
                  buttons.forEach((b, j) => b.classList.toggle("is-selected", j === i))
                }
              }}
            >
              <img
                className="mention-avatar"
                src={user.avatar || "/default-avatar.png"}
                alt=""
              />
              <span className="mention-content">
                <span className="mention-nickname">{user.nickname}</span>
                <span className="mention-username">@{user.username}</span>
              </span>
            </button>
          ))}
        </div>
      </div>
    </div>,
    document.body
  )
}