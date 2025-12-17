from typing import Union, List

from typeguard import typechecked
from ..generated import Service as Service_, Port, ServiceTarget


@typechecked
def Service(*, name: str, target, ports: Union[Port, List[Port]], service_type: str = "NodePort") -> Service_:
    if type(ports) != list:
        ports = [
            ports,
        ]
    class_name = target.__class__.__name__.lower()
    return Service_(type=service_type,
        name=name, target=ServiceTarget(**{class_name: target}), ports=ports
    )
