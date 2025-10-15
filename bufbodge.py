# fixes for buf generation
import os
from pyflakes.api import checkPath
from pyflakes.reporter import Reporter
import io
import re


def find_undefined_vars(path: str):
    out = io.StringIO()
    reporter = Reporter(out, out)
    checkPath(path, reporter)
    return out.getvalue()


def fix_python():
    # no imports
    for dirpath, dirnames, filenames in os.walk("pyflint/generated"):
        for filename in filenames:
            full_path = os.path.join(dirpath, filename)
            print(full_path)
            result = find_undefined_vars(full_path)
            pattern = re.compile(r"undefined name '([^']+)'")
            undefined_names = pattern.findall(result)

            print(undefined_names)
            if not undefined_names:
                continue
            with open(full_path, "r") as file:
                data = file.read()
            imports = ""
            for undefined_name in undefined_names:
                imports += f"from .{undefined_name.lower()} import {undefined_name}\n"
            data = imports + data
            with open(full_path, "w") as file:
                file.write(data)


if __name__ == "__main__":
    fix_python()
