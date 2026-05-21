"use client"

import * as React from "react"

import { saveProfileAction, type UserActionState } from "@/lib/actions/user"
import { AvatarEdit } from "@/components/user/image-upload"
import { Button } from "@/components/ui/button"
import { Input } from "@/components/ui/input"
import { Label } from "@/components/ui/label"
import { Textarea } from "@/components/ui/textarea"
import { useAppState } from "@/components/app/app-provider"
import { useRequiredUser } from "@/components/auth/require-user"
import { apiFetch } from "@/lib/api/client"
import type { UserSummary } from "@/lib/api/types"
import { useI18n } from "@/lib/i18n/provider"
import { toast } from "@/lib/toast"

const initialState: UserActionState = { ok: false }

export function ProfileForm({ user: initialUser }: { user?: UserSummary }) {
  const { t } = useI18n()
  const requiredUser = useRequiredUser()
  const { setCurrentUser } = useAppState()
  const user = initialUser || requiredUser
  const [avatar, setAvatar] = React.useState(user.avatar || "")
  const [profile, setProfile] = React.useState({
    nickname: user.nickname || "",
    description: user.description || "",
    homePage: user.homePage || "",
  })
  const [state, action, pending] = React.useActionState(
    saveProfileAction,
    initialState
  )

  React.useEffect(() => {
    setAvatar(user.avatar || "")
    setProfile({
      nickname: user.nickname || "",
      description: user.description || "",
      homePage: user.homePage || "",
    })
  }, [user.avatar, user.description, user.homePage, user.id, user.nickname])

  React.useEffect(() => {
    if (state.ok) {
      if (state.profile) {
        setAvatar(state.profile.avatar)
        setProfile({
          nickname: state.profile.nickname,
          description: state.profile.description,
          homePage: state.profile.homePage,
        })
        setCurrentUser((current) =>
          current ? { ...current, ...state.profile } : current
        )
        void apiFetch<UserSummary | null>("/api/user/current")
          .then((nextUser) => {
            if (nextUser) {
              setCurrentUser(nextUser)
            }
          })
          .catch(() => {
            // The optimistic profile update above already keeps the form fresh.
          })
      }
      toast.success(t("user.profile.editSuccess"))
    } else if (state.message) {
      toast.error(`${t("user.profile.editFailed")}：${state.message}`)
    }
  }, [setCurrentUser, state, t])

  function updateProfile(next: Partial<typeof profile>) {
    setProfile((current) => ({ ...current, ...next }))
  }

  return (
    <form action={action} className="m-4 space-y-6">
      <input type="hidden" name="userId" value={user.id} />
      <input type="hidden" name="avatar" value={avatar} />
      <div className="space-y-2">
        <Label>{t("user.profile.avatar")}</Label>
        <AvatarEdit value={avatar} onChange={setAvatar} />
      </div>
      <div className="space-y-2">
        <Label htmlFor="nickname">{t("user.profile.nickname")}</Label>
        <Input
          id="nickname"
          name="nickname"
          value={profile.nickname}
          onChange={(event) =>
            updateProfile({ nickname: event.currentTarget.value })
          }
          autoComplete="off"
          placeholder={t("user.profile.nicknamePlaceholder")}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="description">{t("user.profile.signature")}</Label>
        <Textarea
          id="description"
          name="description"
          value={profile.description}
          onChange={(event) =>
            updateProfile({ description: event.currentTarget.value })
          }
          rows={3}
          placeholder={t("user.profile.signaturePlaceholder")}
        />
      </div>
      <div className="space-y-2">
        <Label htmlFor="homePage">{t("user.profile.homepage")}</Label>
        <Input
          id="homePage"
          name="homePage"
          value={profile.homePage}
          onChange={(event) =>
            updateProfile({ homePage: event.currentTarget.value })
          }
          autoComplete="off"
          placeholder={t("user.profile.homepagePlaceholder")}
        />
      </div>
      <div className="flex justify-end pt-4">
        <Button type="submit" disabled={pending}>
          {t("user.profile.save")}
        </Button>
      </div>
    </form>
  )
}
