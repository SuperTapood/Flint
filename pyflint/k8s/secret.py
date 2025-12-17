from ..generated import Secret as _secret, SecretData


def Secret(name, data, secret_type = "opaque"):
    secret_data = [SecretData(key=k, value=v) for k, v in data.items()]
    return _secret(name=name, data=secret_data, type=secret_type)
