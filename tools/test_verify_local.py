#!/usr/bin/env python3
"""Exercise local validation failure paths without running application code."""
import os
import pathlib
import re
import subprocess
import tempfile
import unittest

ROOT = pathlib.Path(__file__).resolve().parents[1]
SCRIPT = ROOT / 'tools/verify-local.sh'
TOOLS = 'go gofmt cargo rustc python3 node npm pnpm helm make psql redis-cli golangci-lint luacheck busted xvfb-run'.split()


class LocalCommandFailure(unittest.TestCase):
    def test_every_group_stops_at_first_failed_tool(self):
        help_result = subprocess.run(['bash', str(SCRIPT), '--help'], capture_output=True, text=True, check=True)
        groups = re.search(r'\[([a-z|]+)\]', help_result.stdout).group(1).split('|')
        with tempfile.TemporaryDirectory(prefix='verify-command-') as directory:
            directory = pathlib.Path(directory)
            log = directory / 'calls'
            for tool in TOOLS:
                stub = directory / tool
                stub.write_text('#!/bin/sh\nprintf "%s\\n" "$0" >> "$FAILURE_LOG"\nexit 73\n')
                stub.chmod(0o755)
            env = {**os.environ, 'PATH': str(directory) + os.pathsep + '/usr/bin:/bin', 'FAILURE_LOG': str(log)}
            # A caller must explicitly attest disposable services. Never inherit
            # that permission into a test which is checking failure behavior.
            env.pop('LOCAL_TEST_SERVICES', None)
            for group in groups:
                with self.subTest(group=group):
                    if log.exists():
                        log.unlink()
                    result = subprocess.run(['/usr/bin/bash', str(SCRIPT), group], cwd=directory, env=env, capture_output=True, text=True, timeout=10)
                    self.assertNotEqual(result.returncode, 0, result.stdout)
                    calls = log.read_text().splitlines() if log.exists() else []
                    self.assertLessEqual(len(calls), 1, 'Validation continued after a tool failed: ' + repr(calls))
                    # A prerequisite refusal may happen before any tool runs.
                    if calls:
                        self.assertEqual(result.returncode, 73, result.stderr)


if __name__ == '__main__':
    unittest.main()
