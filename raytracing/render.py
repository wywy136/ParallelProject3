import numpy as np
import matplotlib.pyplot as plt

arr = np.loadtxt("./output", delimiter="\t")
# print(arr)
plt.imshow(arr, cmap='gray')
plt.savefig("render.jpg")