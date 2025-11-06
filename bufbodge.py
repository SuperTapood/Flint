# fixes for buf generation
import os
from pyflakes.api import checkPath
from pyflakes.reporter import Reporter
import io
import re
import platform
import subprocess
import tokenize
from io import StringIO
import ast


def find_undefined_vars(path: str):
    out = io.StringIO()
    reporter = Reporter(out, out)
    checkPath(path, reporter)
    return out.getvalue()


def fix_python_imports():
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
                if needed_file.split(".")[-1] == "py":
                    with open(needed_file, "r") as maybe_file:
                        if f"class {undefined_name}" in maybe_file.read():
                            common = os.path.commonpath([file, needed_file])
                            needed_import = (
                                "."
                                if os.path.sep
                                in needed_file.replace(common + os.path.sep, "")
                                else ""
                            ) + needed_file.replace(common + os.path.sep, "").replace(
                                ".py", ""
                            ).replace(os.path.sep, ".")
                            imports += (
                                f"from .{needed_import} import {undefined_name}\n"
                            )
        data = imports + data
        with open(file, "w") as f:
            f.write(data)


def get_field_comments(source: str):
    comments = {}
    for tok_type, tok_str, start, _, _ in tokenize.generate_tokens(
        StringIO(source).readline
    ):
        if tok_type == tokenize.COMMENT:
            lineno = start[0]
            comments[lineno] = tok_str.lstrip("# ").strip()
    tree = ast.parse(source)
    result = {}

    for node in tree.body:
        if isinstance(node, ast.ClassDef):
            if any(
                isinstance(dec, ast.Name) and dec.id == "dataclass"
                for dec in node.decorator_list
            ):
                fields = {}
                for stmt in node.body:
                    if isinstance(stmt, ast.AnnAssign) and isinstance(
                        stmt.target, ast.Name
                    ):
                        name = stmt.target.id
                        lineno = stmt.lineno
                        type_str = (
                            ast.unparse(stmt.annotation) if stmt.annotation else None
                        )
                        fields[name] = {
                            "comment": comments.get(lineno - 1, None),
                            "type": type_str,
                        }
                result[node.name] = fields
    return result


def fix_python_pyi():
    for dirpath, _, filenames in os.walk("pyflint/generated"):
        for filename in filenames:
            full_path = os.path.join(dirpath, filename)
            if os.path.basename(full_path).split(os.path.extsep)[1] == "py":
                with open(full_path, "r") as file:
                    data = file.read()
                    comments = get_field_comments(data)
                    imports = data.split("import betterproto")[0]
                if not comments:
                    continue
                with open(f"{full_path}i", "w"):
                    pass

                for class_name, params in comments.items():
                    doc = '        """\n'
                    init = "def __init__(self, \n"
                    for name, data in params.items():
                        doc += f"        :param {name}: {data['comment']}\n"
                        init += f"         {name}: {data['type']},\n"
                    doc += '        """'
                    init += "        ):\n"
                    init += doc
                    pyi = f"""{imports}
class {class_name}:
    {init}
    ...
    """
                    with open(f"{full_path}i", "a") as file:
                        file.write(pyi)


def fix_python():
    fix_python_imports()
    fix_python_pyi()


def run_buf():
    if platform.system() == "Linux":
        subprocess.run(["bunx", "buf", "generate"])
    elif platform.system() == "Windows":
        subprocess.run(["./buf.exe", "generate"])
    else:
        raise NotImplementedError()


if __name__ == "__main__":
    run_buf()
    fix_python()
