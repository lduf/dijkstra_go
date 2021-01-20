import os
import matplotlib.pyplot as plt
import random
import time

def avg(lst):
    return sum(lst)/float(len(lst))

exetime={}
def time_(graph, port):
    print("Executin de 'go run Client.go {} {} > /dev/null'".format(graph, port))
    stime = time.time()
    os.system('go run Client.go {} {} > /dev/null'.format(graph, port))
    return time.time()-stime

exetime=[avg([time_("in/{}.txt".format(x),10001) for _ in range(3)]) for x in range(30,330,30)]

plt.figure()
t = [x for x in range(30,330,30)]
plt.title("Exécution en fonction du nombre de noeuds")
plt.plot(t, exetime)
plt.ylabel('temps d\'éxécution')
plt.xlabel('Nombre d\'éléments')
plt.grid(True)
plt.show()
