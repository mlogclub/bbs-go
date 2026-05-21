import { useEffect } from "react"
import {
  isRouteErrorResponse,
  useLocation,
  useNavigate,
} from "react-router"

import { INSTALL_REQUIRED_STATUS } from "@/lib/api/client"

import { InstallWizard } from "./install-wizard"

export function isInstallRequiredRouteError(error: unknown) {
  return (
    isRouteErrorResponse(error) && error.status === INSTALL_REQUIRED_STATUS
  )
}

export function InstallRequiredFallback() {
  const location = useLocation()
  const navigate = useNavigate()

  useEffect(() => {
    if (location.pathname !== "/install") {
      navigate("/install", { replace: true })
    }
  }, [location.pathname, navigate])

  return <InstallWizard />
}
