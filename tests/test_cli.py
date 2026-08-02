import os
import subprocess
import sys
import tempfile
import unittest


def run(args):
    cmd = [sys.executable, "-m", "env_diff"] + args
    return subprocess.run(cmd, capture_output=True, text=True)


class TestEnvDiff(unittest.TestCase):
    def make(self, text):
        fd, path = tempfile.mkstemp()
        with os.fdopen(fd, "w") as f:
            f.write(text)
        return path

    def test_identical(self):
        a = self.make("A=1\nB=2\n")
        b = self.make("A=1\nB=2\n")
        r = run([a, b])
        self.assertEqual(r.returncode, 0)
        os.unlink(a); os.unlink(b)

    def test_changed(self):
        a = self.make("A=1\n")
        b = self.make("A=2\n")
        r = run([a, b])
        self.assertIn("! A", r.stdout)
        os.unlink(a); os.unlink(b)

    def test_require_identical_fail(self):
        a = self.make("A=1\n")
        b = self.make("A=2\n")
        r = run([a, b, "--require-identical"])
        self.assertEqual(r.returncode, 2)
        os.unlink(a); os.unlink(b)

    def test_only_left(self):
        a = self.make("A=1\n")
        b = self.make("\n")
        r = run([a, b])
        self.assertIn("- A", r.stdout)
        os.unlink(a); os.unlink(b)


if __name__ == "__main__":
    unittest.main()
