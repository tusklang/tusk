from ctypes import *

dll = CDLL("./mydll.dll")
print(dll.test())