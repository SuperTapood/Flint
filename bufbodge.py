# fixes for buf generation
import os
import platform
import subprocess


def force_kwargs(modules):
    for module in modules:
        with open(f"pyflint/generated/{module}/__init__.py", "r") as file:
            data = file.read().replace("@dataclass(", "@dataclass(kw_only=True, ")
        with open(f"pyflint/generated/{module}/__init__.py", "w") as file:
            file.write(data)


def fix_python(modules):
    force_kwargs(modules)


def run_buf(modules):
    for module in modules:
        cmd = (
            "protoc -I protobuf "
            f"--python_betterproto2_out=./pyflint/generated/{module} "
            f"./protobuf/{module}/*"
        )

        subprocess.run(cmd, shell=True, check=True)
    if platform.system() == "Linux":
        subprocess.run(["bunx", "buf", "generate"])
    elif platform.system() == "Windows":
        subprocess.run(["./buf.exe", "generate"])
    else:
        raise NotImplementedError()


if __name__ == "__main__":
    modules = os.listdir("./protobuf")
    run_buf(modules)
    fix_python(modules)
