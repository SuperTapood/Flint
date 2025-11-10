from ..generated.k8s.service_ import Service_, ServiceTarget
from ..generated.base.port import Port


def Service(*, name: str, target, port: Port) -> Service_:
    class_name = target.__class__.__name__.lower()
    return Service_(name=name, target=ServiceTarget(**{class_name: target}), ports=port)
