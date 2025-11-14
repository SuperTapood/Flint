from typeguard import install_import_hook

from typeguard import install_import_hook

install_import_hook("betterproto2")
# Add runtime type checking to a package (and its submodules)
with install_import_hook("k8s"):
    from .k8s.output import Output as K8SOutput
    from .k8s.secret import Secret
    from .k8s.service import Service
    from .k8s.K8SStack import K8SStack
with install_import_hook("generated"):
    from .generated import Deployment, Pod, Port
