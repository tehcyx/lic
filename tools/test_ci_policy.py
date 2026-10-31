#!/usr/bin/env python3
"""Offline regression checks for local-first validation. Requires PyYAML 6.0.3."""
import pathlib
import subprocess
import unittest

import yaml

ROOT = pathlib.Path(__file__).resolve().parents[1]
POLICY = {'ci.yml': []}


class LocalValidationPolicy(unittest.TestCase):
    def test_no_automatic_validation(self):
        for filename, manual_jobs in POLICY.items():
            with self.subTest(workflow=filename):
                data = yaml.load((ROOT / '.github/workflows' / filename).read_text(), Loader=yaml.BaseLoader)
                events = data['on']
                self.assertIn('workflow_dispatch', events)
                self.assertNotIn('pull_request', events)
                if not manual_jobs:
                    self.assertEqual(set(events), {'workflow_dispatch'})
                else:
                    self.assertIn('push', events)
                    for job in manual_jobs:
                        self.assertEqual(data['jobs'][job]['if'], "github.event_name == 'workflow_dispatch'")
                    # Publication must not silently inherit the skipped test gate.
                    for job, spec in data['jobs'].items():
                        if job in manual_jobs:
                            continue
                        needs = spec.get('needs', [])
                        needs = [needs] if isinstance(needs, str) else needs
                        blocked = set(needs) & set(manual_jobs)
                        if blocked:
                            # Journal publication deliberately preserves its manual
                            # test gate and explicitly accepts skipped tests on push.
                            self.assertEqual(filename, 'journal-server.yml')
                            self.assertEqual(job, 'build-and-push')
                            self.assertEqual(spec['if'], "${{ !cancelled() && (github.event_name == 'push' || (github.event_name == 'workflow_dispatch' && needs.test.result == 'success')) }}")

    def test_workflow_dependency_graphs(self):
        for file in (ROOT / '.github/workflows').glob('*'):
            data = yaml.load(file.read_text(), Loader=yaml.BaseLoader)
            jobs = data['jobs']
            def visit(job, ancestors):
                self.assertNotIn(job, ancestors, f'Cycle in {file.name}')
                needs = jobs[job].get('needs', [])
                for dependency in [needs] if isinstance(needs, str) else needs:
                    self.assertIn(dependency, jobs, f'Missing job in {file.name}')
                    visit(dependency, ancestors | {job})
            for job in jobs:
                visit(job, set())

    def test_local_command_contract(self):
        script = ROOT / ('bin/verify' if ROOT.name == 'agent-workspace' else 'tools/verify-local.sh')
        # Worktree directory names need not equal repository names.
        if 'verify.yml' in POLICY:
            script = ROOT / 'bin/verify'
        self.assertTrue(script.is_file())
        subprocess.run(['bash', '-n', str(script)], check=True)
        if 'verify.yml' not in POLICY:
            result = subprocess.run(['bash', str(script), '--help'], cwd='/tmp', capture_output=True, text=True)
            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertIn('Usage:', result.stdout)
            result = subprocess.run(['bash', str(script), 'not-a-validation-group'], cwd='/tmp', capture_output=True, text=True)
            self.assertEqual(result.returncode, 2, result.stderr)

    def test_bot_merge_has_not_lost_its_test_gate(self):
        if 'master.yml' in POLICY or ('ci.yml' in POLICY and (ROOT / '.github/workflows/codeql-analysis.yml').exists()):
            self.assertFalse((ROOT / '.github/workflows/dependabot-automerge.yml').exists())

    def test_portal_build_inputs_and_chart_exception(self):
        if 'portal.yml' not in POLICY:
            return
        data = yaml.load((ROOT / '.github/workflows/portal.yml').read_text(), Loader=yaml.BaseLoader)
        self.assertEqual(data['jobs']['build-binary']['needs'], ['frontend'])
        self.assertEqual(data['jobs']['build-and-push']['needs'], ['build-binary'])
        frontend = data['jobs']['frontend']['steps']
        tests = next(s for s in frontend if s.get('name') == 'Frontend tests')
        self.assertEqual(tests['if'], "github.event_name == 'workflow_dispatch'")
        self.assertTrue(any(s.get('run') == 'pnpm run build' and 'if' not in s for s in frontend))
        charts = yaml.load((ROOT / '.github/workflows/charts.yml').read_text(), Loader=yaml.BaseLoader)
        self.assertIn('renovate/**', charts['on']['push']['branches'])
        self.assertIn('pull_request', charts['on'])
        self.assertEqual(set(charts['jobs']), {'lint-template', 'storage-quantities'})


if __name__ == '__main__':
    unittest.main()
