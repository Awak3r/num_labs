#!/usr/bin/env python3
import sys

import numpy as np

nums = list(map(float, open(sys.argv[1] if len(sys.argv) > 1 else "data/qr.txt").read().split()))
n = int(nums[0])
a = np.array(nums[1 : 1 + n * n]).reshape(n, n)

eig = np.linalg.eigvals(a)
for i, l in enumerate(eig):
    if abs(l.imag) < 1e-12:
        print(f"λ{i} = {l.real:.6f}")
    else:
        print(f"λ{i} = {l.real:.6f} {l.imag:+.6f}i")