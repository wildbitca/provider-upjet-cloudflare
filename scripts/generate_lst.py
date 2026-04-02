#!/usr/bin/env python3
"""Generate config/generated.lst (JSON array of resource names) from provider-metadata.yaml.

Used by make schema-version-diff to compare Terraform schema versions across provider upgrades.
Reads the 'resources' keys from provider-metadata.yaml and writes a JSON array to generated.lst.
"""
import json
import re
import sys


def main():
    if len(sys.argv) != 3:
        print("usage: generate_lst.py <provider-metadata.yaml> <generated.lst>", file=sys.stderr)
        sys.exit(1)
    meta_path = sys.argv[1]
    out_path = sys.argv[2]

    resources = []
    in_resources = False
    # Resource keys under top-level "resources:" are indented with exactly 4 spaces.
    resource_line = re.compile(r"^    ([a-z0-9_]+):\s*$")

    with open(meta_path) as f:
        for line in f:
            if line.rstrip() == "resources:":
                in_resources = True
                continue
            if in_resources:
                # Top-level key (no leading space) ends the resources block
                if line and not line[0].isspace():
                    in_resources = False
                    continue
                m = resource_line.match(line)
                if m:
                    resources.append(m.group(1))

    resources.sort()
    with open(out_path, "w") as f:
        json.dump(resources, f, indent=2)
    print(f"Wrote {len(resources)} resource names to {out_path}")


if __name__ == "__main__":
    main()
