#!/usr/bin/env python3
import numpy as np

nums = list(map(float, open("data/qr.txt").read().split()))
n = int(nums[0])
a = np.array(nums[1 : 1 + n * n]).reshape(n, n)

eig = np.linalg.eigvals(a)
for i, l in enumerate(eig):
    if abs(l.imag) < 1e-12:
        print(f"λ{i} = {l.real:.6f}")
    else:
        print(f"λ{i} = {l.real:.6f} {l.imag:+.6f}i")
print(f"сумма СЗ = {eig.sum():.6f}, trace(A) = {np.trace(a):.6f}")
print(f"произведение СЗ = {eig.prod():.6f}, det(A) = {np.linalg.det(a):.6f}")