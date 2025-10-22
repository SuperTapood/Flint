from pyflint import *

pod = Pod(
    name="nginx",
    image="nginx:latest",
    ports=[
        80,
    ],
)
stack = K8SStack(
    api="https://192.168.49.2:8443",
    token="eyJhbGciOiJSUzI1NiIsImtpZCI6ImxXUGU0UEIwZWtaRVlXaHM5TEVENmFzV2FSQTJPRi1ndkVHQ2hOUlNEdWMifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzYxMTY0NDE1LCJpYXQiOjE3NjExNjA4MTUsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiNDFkMmJiM2QtN2ZmYy00YmExLTg4NzctMWVhNjgwOTg4MTY0Iiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImZsaW50IiwidWlkIjoiNTNlMDM5NzYtOWQxNi00MjgzLTlkYTAtM2QxZTIxMWQ5YmM2In19LCJuYmYiOjE3NjExNjA4MTUsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmZsaW50In0.owxBYV0SPMfSoFto0rmKUt6zP3NdrVyudzrA404N24n-nA2DWobFx-k-6OQkovWkOjvteqCFSj_qqp3TdzuBUXkBIMbMMLzVvcVuIVyFvXl9K39Ru9pY4Cn0hCt3dYgLGvi_M079g97gwLD99sCwLxH1f4TF9DWX4gvmjqcqNGGkJt0zG4W65xVKLxx0dCVgCiSh-oihgE0tP78AthF4CVJ6mudavfQugnhoBh14jKkI3aybsoZJL-GpvmrKSr1KeQxuK38XWxFCbbE9kxkOpsuDYpOWlQUG8f8MLMf1Y3UT5DN7nKXPCENPE7-8ccdlA6SPDMRv_2hWb779fgYlig"
)
stack.add_objects(pod)

# with open("bob.bin", "wb") as file:
#     file.write(pod.SerializeToString())
stack.synth()
# print(len(open("bob.bin", "rb").read()))
