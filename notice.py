# This script adds the required copyright header to all .go files
# so i dont have to myself

# i am very lazy

import os
import re

header = """/*
Copyright © 2025 this guy Labs <thisguy@thisguylabs.com>

This file is part of GVT (Guy's Versioning Tool).

GVT is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

GVT is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with GVT. If not, see <https://www.gnu.org/licenses/>.

Do not remove or modify this notice.
*/
"""

skip_dirs = {".gvt", "vendor"}

header_pattern = re.compile(r"(?s)^/\*.*?\*/\s*")

def add_or_replace_header(file_path):
    with open(file_path, "r", encoding="utf-8") as f:
        content = f.read()

    if "This file is part of GVT (Guy's Versioning Tool)." in content:
        return

    match = header_pattern.match(content)
    if match:
        content = content[match.end():]

    with open(file_path, "w", encoding="utf-8") as f:
        f.write(header + "\n" + content)
    print(f"Header updated for {file_path}")

def walk_and_update_headers(root="."):
    for dirpath, dirnames, filenames in os.walk(root):
        dirnames[:] = [d for d in dirnames if d not in skip_dirs]

        for filename in filenames:
            if filename.endswith(".go"):
                file_path = os.path.join(dirpath, filename)
                add_or_replace_header(file_path)

if __name__ == "__main__":
    walk_and_update_headers()
