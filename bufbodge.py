# fixes for buf generation
import os
from pyflakes.api import checkPath
from pyflakes.reporter import Reporter
import io
import re
import platform
import subprocess


def find_undefined_vars(path: str):
    out = io.StringIO()
    reporter = Reporter(out, out)
    checkPath(path, reporter)
    return out.getvalue()


def fix_python():
    # no imports
    files = []
    for dirpath, _, filenames in os.walk("pyflint/generated"):
        for filename in filenames:
            full_path = os.path.join(dirpath, filename)
            files.append(full_path)
    for file in files:
        result = find_undefined_vars(file)
        pattern = re.compile(r"undefined name '([^']+)'")
        undefined_names = pattern.findall(result)

        if not undefined_names:
            continue
        with open(file, "r") as f:
            data = f.read()
        imports = ""
        for undefined_name in undefined_names:
            for needed_file in files:
                if os.path.basename(needed_file) == undefined_name.lower() + ".py":
                    common = os.path.commonpath([file, needed_file])
                    needed_import = ("." if os.path.sep in needed_file.replace(common + os.path.sep, "") else "") + needed_file.replace(common + os.path.sep, "").replace(".py", "").replace(os.path.sep, ".")
                    imports += f"from .{needed_import} import {undefined_name}\n"
        data = imports + data
        with open(file, "w") as f:
            f.write(data)


if __name__ == "__main__":
    if platform.system() == "Linux":
        subprocess.run(["npx", "buf", "generate"])
    elif platform.system() == "Windows":
        subprocess.run(["./buf.exe", "generate"])
    else:
        raise NotImplementedError()
    fix_python()
