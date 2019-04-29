import random
import signal
import sys
import time

index = 0
run = True
scale = 1.0
if len(sys.argv) > 1:
    scale = float(sys.argv[1])


def sig_handler(signal, frame):
    global run
    print("shutting down...")
    run = False

signal.signal(signal.SIGINT, sig_handler)
signal.signal(signal.SIGTERM, sig_handler)

while run:
    print(f'{index}')
    index = index + 1
    time.sleep(random.random()*scale)
    
