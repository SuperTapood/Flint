from .service_ import Service_, ServiceTarget
from ..base.port import Port


def Service(name: str, target, port: Port) -> Service_:
    class_name = target.__class__.__name__.lower()
    return Service_(name, ServiceTarget(**{class_name: target}), port)
