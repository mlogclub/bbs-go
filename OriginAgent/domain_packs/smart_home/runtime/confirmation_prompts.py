"""Smart-home confirmation prompt builder."""

from __future__ import annotations

from OriginAgent.agent.action_safety import ActionDecision, ActionRequest
from OriginAgent.agent.confirmation import PROMPT_MAX_CHARS, _human_action, _sanitize_text


class SmartHomeConfirmationPromptBuilder:
    def build(
        self,
        request: ActionRequest,
        decision: ActionDecision,
        *,
        kind: str,
    ) -> str:
        action_text = _human_action(request)
        if kind == "notify_only":
            return _sanitize_text(f"我会通知你：{action_text}。", PROMPT_MAX_CHARS)
        if request.requires_presence_empty and decision.presence_status != "empty":
            return _sanitize_text(
                f"这个动作需要确认家中无人，但我现在无法确认。是否仅本次继续{action_text}？",
                PROMPT_MAX_CHARS,
            )
        if decision.presence_status == "unknown":
            return _sanitize_text(
                f"我不确定家里是否还有人，因此不会自动执行{action_text}。你现在要继续吗？",
                PROMPT_MAX_CHARS,
            )
        if "fact" in decision.reason.casefold() or decision.pending_facts:
            return _sanitize_text(
                "这条动作依赖一个尚未确认的家庭规则。要仅本次继续，还是取消？",
                PROMPT_MAX_CHARS,
            )
        return _sanitize_text(f"是否仅本次继续{action_text}？", PROMPT_MAX_CHARS)
