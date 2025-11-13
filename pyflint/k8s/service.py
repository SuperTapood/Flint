from typing import Union, List
from ..generated.test import Service as Service_, Port, ServiceTarget


def Service(*, name: str, target, ports: Union[Port, List[Port]]) -> Service_:
    if type(ports) != list:
        ports = [
            ports,
        ]
    class_name = target.__class__.__name__.lower()
    return Service_(
        name=name, target=ServiceTarget(**{class_name: target}), ports=ports
    )
