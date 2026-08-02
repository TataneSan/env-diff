#!/usr/bin/env python3
"""
Compare two .env files.

Shows only-left, only-right, and changed variables.
Exit codes:
  0  success
  1  I/O or CLI error
  2  gate condition not met (require-identical / no-conflicts)
"""
import argparse
import sys
import json


def parse_env(path):
    data = {}
    try:
        with open(path, "r", encoding="utf-8") as f:
            for line in f:
                line = line.strip()
                if not line or line.startswith("#"):
                    continue
                if "=" in line:
                    k, v = line.split("=", 1)
                    data[k.strip()] = v.strip()
    except OSError as e:
        raise
    return data


def main(argv=None):
    parser = argparse.ArgumentParser(description="Compare two .env files.")
    parser.add_argument("left", help="Left .env file")
    parser.add_argument("right", help="Right .env file")
    parser.add_argument("--mask-values", action="store_true",
                        help="Mask values in output")
    parser.add_argument("--shell-export", action="store_true",
                        help="Print shell export lines for changed vars")
    parser.add_argument("--require-identical", action="store_true",
                        help="Exit 2 if files differ")
    parser.add_argument("--no-conflicts", action="store_true",
                        help="Exit 2 if any variable changed")
    parser.add_argument("--json", action="store_true",
                        help="Emit JSON report")
    args = parser.parse_args(argv)

    try:
        left = parse_env(args.left)
        right = parse_env(args.right)
    except OSError as e:
        print(f"error: {e}", file=sys.stderr)
        return 1

    only_left = sorted(set(left) - set(right))
    only_right = sorted(set(right) - set(left))
    changed = []
    for k in sorted(set(left) & set(right)):
        if left[k] != right[k]:
            changed.append({
                "key": k,
                "left": left[k] if not args.mask_values else "***",
                "right": right[k] if not args.mask_values else "***",
            })

    # Gates
    if args.require_identical and (only_left or only_right or changed):
        msg = "files are not identical"
        if args.json:
            print(json.dumps({"error": msg}))
        else:
            print(msg, file=sys.stderr)
        return 2
    if args.no_conflicts and changed:
        msg = f"{len(changed)} changed variable(s) (no-conflicts violated)"
        if args.json:
            print(json.dumps({"error": msg, "changed": changed}, ensure_ascii=False))
        else:
            print(msg, file=sys.stderr)
        return 2

    if args.json:
        print(json.dumps({
            "only_left": only_left,
            "only_right": only_right,
            "changed": changed,
        }, ensure_ascii=False, indent=2))
    else:
        for k in only_left:
            print(f"- {k}")
        for k in only_right:
            print(f"+ {k}")
        for c in changed:
            print(f"! {c['key']}: {c['left']} -> {c['right']}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
