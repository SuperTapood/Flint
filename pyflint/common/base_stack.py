from betterproto2 import Message
import traceback
import socket
import sys


class BaseStack:
    def __init__(self):
        self._prev_post_init = None
        self._prev_getitem = None

    def __enter__(self):
        self._prev_post_init = Message.__post_init__
        # self._prev_getitem = Message.__getitem__

        def post_init(obj):
            obj._unknown_fields = b""
            self.add_objects(obj)

        Message.__post_init__ = post_init

        def getitem(obj, key):
            if key == "":
                key = None
            return self.lookup(obj, key)
        
        Message.__getitem__ = getitem

    # def __exit__(self, exc_type, exc, tb: traceback):
    #     Message.__post_init__ = self._prev_post_init
    #     Message.__getitem__ = self._prev_getitem

    #     if exc:
    #         # Store the exception for later use
    #         self.error = exc
    #         print(exc)
    #         print(exc_type)

    #         # Decide if you want to suppress the exception:
    #         # return True   → suppress
    #         # return False  → re-raise
    #         return False  # ← change to True if you want to swallow errors

    def send_data(self, data):
        if len(sys.argv) < 2:
            return
        socket_path = sys.argv[1]
        # Connect to Unix socket
        sock = socket.socket(socket.AF_UNIX, socket.SOCK_STREAM)
        sock.connect(socket_path)
        sock.sendall(data)
        sock.close()
    
    def lookup(self, obj, key):
        raise NotImplementedError()

