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
    token="eyJhbGciOiJSUzI1NiIsImtpZCI6IllEWm9BVHpfNTFDdkE0QkVWcjM2X21iVVRYYkhRZFY0bUVIUklCSzVRX1kifQ.eyJhdWQiOlsiaHR0cHM6Ly9rdWJlcm5ldGVzLmRlZmF1bHQuc3ZjLmNsdXN0ZXIubG9jYWwiXSwiZXhwIjoxNzYxNDAxMDU0LCJpYXQiOjE3NjEzOTc0NTQsImlzcyI6Imh0dHBzOi8va3ViZXJuZXRlcy5kZWZhdWx0LnN2Yy5jbHVzdGVyLmxvY2FsIiwianRpIjoiMGM4NjJhN2YtOWIzNi00NzcxLThlN2YtNjQ5OWY4NDg4NTNjIiwia3ViZXJuZXRlcy5pbyI6eyJuYW1lc3BhY2UiOiJkZWZhdWx0Iiwic2VydmljZWFjY291bnQiOnsibmFtZSI6ImZsaW50IiwidWlkIjoiMWJlNGM0NGItY2UyMy00MWRiLTg0ZWEtNzc5ZTU1NTUxYWZiIn19LCJuYmYiOjE3NjEzOTc0NTQsInN1YiI6InN5c3RlbTpzZXJ2aWNlYWNjb3VudDpkZWZhdWx0OmZsaW50In0.niFDpFBYK0iVscZaIlN_QyRfKQaCUukpAgWnqcXTQfDIUl2naIat_CCCmCT1pWGjwMrSzIV5E38no8IIuA0NY-6ih8pFsaydzRgjibjtkqiGg1dLsKeqWUSj91hhSVyfqRQk1_nKTdeCbidCoGjTtQz8WSHv9TdoD8MTsWj3fKRofIUs3gkgKbsop7I3EYrkR6Z5xS4Yj9i0I4poCaGhO8Fzi1elgUzJPbWIAe4-5KxxWXwOvzAwf-NfcF74HWFzf0i0SYICTId9A02asejt29b2_djxaUo-65gpchzUcnEu4Upvs3KccKMmE8S4exD6n6E9Txh7seoCoGmzM9ioWg"
)
stack.add_objects(pod)

# with open("bob.bin", "wb") as file:
#     file.write(pod.SerializeToString())
stack.synth()
# print(len(open("bob.bin", "rb").read()))
