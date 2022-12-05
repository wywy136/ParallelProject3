import subprocess
import matplotlib.pyplot as plt

x = [0, 1, 2, 3, 4]
threadNums = ["2", "4", "6", "8", "12"]
numRays = ["1000000", "10000000", "100000000"]

if __name__ == "__main__":
    sTime = dict()
    wsTime = dict()
    wbTime = dict()
    
    for numRay in numRays:
        print(f"Sequential - ray: {numRay}")
        sumTime = 0.0
        
        for i in range(5):
            pipe = subprocess.Popen(
                ["go", "run", "main.go", numRay, "1000", "s"],
                stdout=subprocess.PIPE,
                stderr=subprocess.PIPE
            )
            stdout, stderr = pipe.communicate()
            res = stdout.decode().split(":")[1].strip(" ")
            sumTime += float(res)
        sTime[numRay] = sumTime / 5
        
        ws = []
        for nThread in threadNums:
            print(f"Work Stealing - ray: {numRay} - thread: {nThread}")
            
            sumTime = 0.0
            for i in range(5):
                pipe = subprocess.Popen(
                    ["go", "run", "main.go", numRay, "1000", "ws", nThread, "20"],
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE
                )
                stdout, stderr = pipe.communicate()
                # print(stdout.decode())
                res = stdout.decode().split(":")[1].strip(" ")
                sumTime += float(res)
            ws.append(sTime[numRay] / (sumTime / 5))
        wsTime[numRay] = ws
        
        wb = []
        for nThread in threadNums:
            print(f"Work Balancing - ray: {numRay} - thread: {nThread}")
            
            sumTime = 0.0
            for i in range(5):
                pipe = subprocess.Popen(
                    ["go", "run", "main.go", numRay, "1000", "wb", nThread, "20"],
                    stdout=subprocess.PIPE,
                    stderr=subprocess.PIPE
                )
                stdout, stderr = pipe.communicate()
                # print(stdout.decode())
                res = stdout.decode().split(":")[1].strip(" ")
                sumTime += float(res)
            wb.append(sTime[numRay] / (sumTime / 5))
        wbTime[numRay] = wb
        
    fig, ax = plt.subplots()
    small, = ax.plot(x, wsTime["1000000"], linestyle='--', marker='o')
    small.set_label("1000000")

    mid, = ax.plot(x, wsTime["10000000"], linestyle='--', marker='o')
    mid.set_label("10000000")

    # large, = ax.plot(x, wsTime["100000000"], linestyle='--', marker='o')
    # large.set_label("100000000")

    ax.legend(loc='upper left')
    ax.set_xticks(x)
    ax.set_xticklabels(threadNums)
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup')
    plt.savefig(f'speedup-ws-10.png')
    
    fig1, ax1 = plt.subplots()
    small1, = ax1.plot(x, wbTime["1000000"], linestyle='--', marker='o')
    small1.set_label("1000000")

    mid1, = ax1.plot(x, wbTime["10000000"], linestyle='--', marker='o')
    mid1.set_label("10000000")

    # large1, = ax1.plot(x, wbTime["100000000"], linestyle='--', marker='o')
    # large1.set_label("100000000")

    ax1.legend(loc='upper left')
    ax1.set_xticks(x)
    ax1.set_xticklabels(threadNums)
    plt.xlabel('Number of Threads')
    plt.ylabel('Speedup')
    plt.savefig(f'speedup-wb-10.png')