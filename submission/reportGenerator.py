import os
import random
import subprocess
import time
import numpy as np
import matplotlib.pyplot as plt
import math

sort_program = "./sort"
gensort = "./gensort"
#temp_dir = "C:\\TheBulk\\16_SP-25\\CSE-124\\sort_benchmark"
#os.makedirs(temp_dir, exist_ok=True)

sizes = np.linspace(1024, 10 * 1024*1024, 20, dtype=int)
runtimes = []
asymptotic_bound = []

def generate_test_file(filepath, size):
  with open(filepath, "wb") as f:
    # size in bytes
    len_size = 4
    key_size = 10
    val_size = random.randint(0, 1010)
    rec_size = key_size + val_size

    written = 0
    while written + rec_size <= size:
      len = (rec_size).to_bytes(4, byteorder="big")
      key = os.urandom(key_size)
      val = os.urandom(val_size)
      f.write(len + key + val)
      written += rec_size

for size in sizes:
  input_file = os.path.join(f"input_{size//1024}kb.dat")
  output_file = os.path.join(f"output_{size//1024}kb.dat")
  randseed = 1000
  size_str = f"\"{size}\""
  #generate_test_file(input_file, size)

  subprocess.run([gensort, "-randseed", str(randseed), size_str, input_file], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)

  start = time.time()
  result = subprocess.run([sort_program, input_file, output_file], stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
  end = time.time()
  
  runtimes.append(end - start)

  n = size / 1024
  asymptotic_bound.append(n * math.log2(n))

asymptotic_bound = np.array(asymptotic_bound)
asymptotic_bound = asymptotic_bound / max(asymptotic_bound) * max(runtimes)

plt.figure(figsize=(12, 6))
plt.plot(sizes, runtimes, label="Actual Runtime", marker='o')
plt.plot(sizes, asymptotic_bound, label="O(n log n) Bound (n = KB)", linestyle='--')
plt.xlabel("Input Size (bytes)")
plt.ylabel("Time (seconds)")
plt.title("Sort Runtime vs. Input Size")
plt.legend()
plt.grid(True)
plt.tight_layout()
plt.savefig("report.pdf")
print("✅ report.pdf generated.")