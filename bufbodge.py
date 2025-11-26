# fixes for buf generation
import os
import platform
import subprocess


def force_kwargs():
    with open(f"pyflint/generated/__init__.py", "r") as file:
        data = file.read().replace("@dataclass(", "@dataclass(kw_only=True, ")

    with open(f"pyflint/generated/__init__.py", "w") as file:
        file.write(data)


def fix_python():
    force_kwargs()


def run_buf():
    cmd = (
        "protoc -I protobuf "
        f"--python_betterproto2_out=./pyflint/generated "
        f"./protobuf/*/*"
    )

    subprocess.run(cmd, shell=True, check=True)
    if platform.system() == "Linux":
        subprocess.run(["bunx", "buf", "generate"])
    elif platform.system() == "Windows":
        subprocess.run(["./buf.exe", "generate"])
    else:
        raise NotImplementedError()


if __name__ == "__main__":
    run_buf()
    fix_python()
