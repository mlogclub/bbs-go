"""Central limits for built-in tools."""

from __future__ import annotations

from dataclasses import dataclass


@dataclass(frozen=True)
class ToolLimits:
    """Runtime limits shared by tool implementations."""

    read_file_max_chars: int = 128_000
    read_file_default_limit: int = 2_000
    read_file_max_pdf_pages: int = 20
    exec_max_output_chars: int = 10_000
    exec_max_timeout_seconds: int = 600
    web_fetch_max_chars: int = 50_000
    grep_max_files: int = 5_000
    grep_max_file_bytes: int = 2_000_000
    grep_max_result_chars: int = 128_000
    best_window_max_scan_lines: int = 2_000
    max_input_bytes: int = 5_000_000
    web_fetch_max_bytes: int = 5_000_000
    mcp_response_max_chars: int = 128_000
    mcp_resource_max_bytes: int = 2_000_000
