#!/usr/bin/env python3

"""
Automated smoke test for the Student API.

The script optionally starts the server, waits for /health to succeed, then
exercises every CRUD endpoint (list, create, get, update, delete) to ensure
the responses match expectations.
"""

from __future__ import annotations

import atexit
import json
import os
import shutil
import subprocess
import sys
import tempfile
import time
import uuid
from pathlib import Path
from typing import Dict, Optional, Tuple
from urllib import error, request


REPO_ROOT = Path(__file__).resolve().parents[1]
BASE_URL = os.environ.get("BASE_URL", "http://localhost:8080")
API_PREFIX = os.environ.get("API_PREFIX", "/api/v1")
START_SERVER = os.environ.get("START_SERVER", "1").lower() not in {"0", "false", "no"}


class SmokeTestError(RuntimeError):
    """Custom error for clearer failures."""


def start_server(log_path: Path) -> subprocess.Popen:
    """Start the Go API server and stream logs to the supplied path."""
    log_file = log_path.open("w")
    proc = subprocess.Popen(
        ["go", "run", "cmd/api/main.go"],
        cwd=str(REPO_ROOT),
        stdout=log_file,
        stderr=subprocess.STDOUT,
    )
    return proc


def stop_server(proc: Optional[subprocess.Popen]) -> None:
    """Shut the server down gracefully."""
    if proc is None:
        return
    if proc.poll() is not None:
        return
    proc.terminate()
    try:
        proc.wait(timeout=5)
    except subprocess.TimeoutExpired:
        proc.kill()


def call_api(method: str, path: str, payload: Optional[Dict] = None) -> Tuple[int, str, Dict]:
    """Call the API and return status code, raw body, and parsed JSON."""
    url = f"{BASE_URL}{path}"
    data_bytes = json.dumps(payload).encode("utf-8") if payload is not None else None
    headers = {"Content-Type": "application/json"} if payload is not None else {}
    req = request.Request(url, data=data_bytes, method=method, headers=headers)
    try:
        with request.urlopen(req, timeout=5) as resp:
            body_bytes = resp.read()
            status_code = resp.getcode()
    except error.HTTPError as exc:
        body_bytes = exc.read()
        status_code = exc.code
    except error.URLError as exc:
        raise SmokeTestError(f"Request to {url} failed: {exc}") from exc

    body_text = body_bytes.decode("utf-8")
    try:
        parsed = json.loads(body_text)
    except json.JSONDecodeError as exc:
        raise SmokeTestError(f"Response from {url} was not valid JSON: {exc}") from exc
    return status_code, body_text, parsed


def wait_for_health(timeout_seconds: int = 20) -> None:
    """Poll /health until the server reports healthy."""
    deadline = time.time() + timeout_seconds
    while time.time() < deadline:
        try:
            status, _, body = call_api("GET", "/health")
            if status == 200 and body.get("status") == "healthy":
                return
        except SmokeTestError:
            pass
        time.sleep(0.5)
    raise SmokeTestError("Server did not become healthy within timeout")


def assert_equal(actual, expected, message: str) -> None:
    if actual != expected:
        raise SmokeTestError(f"{message}: expected {expected!r}, got {actual!r}")


def main() -> None:
    tmp_dir = Path(tempfile.mkdtemp(prefix="student-api-smoke-"))
    log_path = tmp_dir / "server.log"
    server_proc: Optional[subprocess.Popen] = None

    atexit.register(lambda: shutil.rmtree(tmp_dir, ignore_errors=True))
    atexit.register(lambda: stop_server(server_proc))

    if START_SERVER:
        print("Starting API server ...", flush=True)
        server_proc = start_server(log_path)
        time.sleep(0.5)
    else:
        print("Skipping server startup because START_SERVER is set to false.")

    print("Waiting for server health ...", flush=True)
    wait_for_health()

    print("Checking /health ...", flush=True)
    status, _, body = call_api("GET", "/health")
    assert_equal(status, 200, "/health status code mismatch")
    assert_equal(body.get("status"), "healthy", "/health payload mismatch")

    print("Listing students ...", flush=True)
    status, _, body = call_api("GET", f"{API_PREFIX}/students")
    assert_equal(status, 200, "List status code mismatch")
    data = body.get("data")
    if not isinstance(data, list):
        raise SmokeTestError("List response did not return an array in data")

    unique_email = f"smoke-{uuid.uuid4()}@example.com"
    student_payload = {
        "name": "Smoke Test Student",
        "email": unique_email,
        "age": 23,
        "major": "Computer Science",
        "gpa": 3.75,
    }
    print("Creating student ...", flush=True)
    status, _, body = call_api("POST", f"{API_PREFIX}/students", student_payload)
    assert_equal(status, 201, "Create status code mismatch")
    student = body.get("data", {})
    student_id = student.get("id")
    if not student_id:
        raise SmokeTestError("Create response did not include an id")
    assert_equal(student.get("email"), unique_email, "Create payload email mismatch")

    print(f"Retrieving student {student_id} ...", flush=True)
    status, _, body = call_api("GET", f"{API_PREFIX}/students/{student_id}")
    assert_equal(status, 200, "Get status code mismatch")
    fetched = body.get("data", {})
    assert_equal(fetched.get("id"), student_id, "Get payload mismatch")

    update_payload = {
        "major": "Applied Mathematics",
        "gpa": 3.9,
    }
    print("Updating student ...", flush=True)
    status, _, body = call_api("PUT", f"{API_PREFIX}/students/{student_id}", update_payload)
    assert_equal(status, 200, "Update status code mismatch")
    updated = body.get("data", {})
    assert_equal(updated.get("major"), "Applied Mathematics", "Update major mismatch")
    assert_equal(float(updated.get("gpa", 0)), 3.9, "Update GPA mismatch")

    print("Deleting student ...", flush=True)
    status, _, body = call_api("DELETE", f"{API_PREFIX}/students/{student_id}")
    assert_equal(status, 200, "Delete status code mismatch")
    assert_equal(body.get("message"), "Student deleted successfully", "Delete payload mismatch")

    print("Confirming deletion ...", flush=True)
    status, _, body = call_api("GET", f"{API_PREFIX}/students")
    assert_equal(status, 200, "List status code mismatch after delete")
    data = body.get("data", [])
    if any(item.get("id") == student_id for item in data):
        raise SmokeTestError("Deleted student still present in list response")

    print("All API smoke tests passed!", flush=True)
    if server_proc:
        print(f"Server logs are available at: {log_path}")


if __name__ == "__main__":
    try:
        main()
    except SmokeTestError as exc:
        print(f"[FAIL] {exc}", file=sys.stderr)
        sys.exit(1)
