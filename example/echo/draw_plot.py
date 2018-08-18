##
# usage:   python3 xxx.py default_latency.txt improved_latency.txt  start_from_line  totaly_read_lines 
#python3 throughput_latency.py 1000 300000 start end increment data_path_for_folders
##

import numpy as np
import sys
import statsmodels.api as sm
import matplotlib
#matplotlib.use('Agg')
import matplotlib.pyplot as plt
import os

plt.switch_backend('agg')

path = sys.argv[6]


files = os.listdir(path)

start = int(sys.argv[3])
end = int(sys.argv[4])
increament = int(sys.argv[5])

dataset_quic = []
percentile_set_quic = []
throughputs_quic = []

dataset_tcp = []
percentile_set_tcp = []
throughputs_tcp = []

print(files)
files.sort()
#a = files[1]
#files.remove("latency_tcp_5000_10.txt")
#files.append(a)
print(files)


mylist = []

#for i in range(1, 63, 5):
for i in range(start, end + 1, increament):
	#print(i)
	mylist.append(str(i))
	throughputs_quic.append(i/1000)# * 10000)

for i in mylist:
	try:
		dataset_quic.append(np.genfromtxt(path + '/' + i + '.log', delimiter=',', skip_header=int(sys.argv[1]), max_rows=int(sys.argv[2])))	
	except:
		print("Error", i, "not found")
		sys.exit()
		continue

for i in dataset_quic:
	percentile_quic = np.percentile(i, 99,axis=0)
	percentile_set_quic.append(percentile_quic)
	#print(percentile)


print(mylist)
print(len(throughputs_quic))
print(len(percentile_set_quic))
'''
path = "/home/alireza/Documents/quic_results/tcp"
files = os.listdir(path)

files.sort()
a = files[1]
files.remove("latency_tcp_5000_10.txt")
files.append(a)
print(files)

for i in files:
	dataset_tcp.append(np.genfromtxt(path + '/' + i, delimiter=',', skip_header=int(sys.argv[1]), max_rows=int(sys.argv[2])))	

for i in dataset_tcp:
	percentile_tcp = np.percentile(i, 99,axis=0)
	percentile_set_tcp.append(percentile_tcp)
	#print(percentile)

for i in range(5000, 50001, 5000):
	#print(i)
	throughputs_tcp.append(i)

'''



#p = plt.plot(throughputs_quic, percentile_set_quic, 'r', throughputs_tcp, percentile_set_tcp, 'g')
#p = plt.plot(throughputs_tcp, percentile_set_tcp, 'g')
#plt.axis([0, 0.3, -10, 2000])
p = plt.plot(throughputs_quic, percentile_set_quic, 'o')
plt.xlabel('Throughput(MRPS)')
plt.ylabel('Latency at $99^{th}$(us)')
plt.setp(p, linewidth=2.0)
plt.grid(True, which='both')
plt.savefig('plots/tcp_zoom.png', format='png', dpi=300)
#plt.show()

'''
intput_file_count = sys.argv[1]

dataset = []
for i in range(intput_file_count):
	dataset.append(np.genfromtxt(sys.argv[1], delimiter=',', skip_header=int(sys.argv[3]), max_rows=int(sys.argv[4])))

a = np.genfromtxt(sys.argv[1], delimiter=',', skip_header=int(sys.argv[3]), max_rows=int(sys.argv[4]))
a2 = np.genfromtxt(sys.argv[2], delimiter=',', skip_header=int(sys.argv[3]), max_rows=int(sys.argv[4]))

print('removing these from default file:')
while a.max() > 1000000:
    t = np.argmax(a, axis=0)
    print(a.max(),t)
    a = np.delete(a, t)

print('removing these from improved file:')
while a2.max() > 1000000:
    t = np.argmax(a2, axis=0)
    print(a.max(), t)
    a2 = np.delete(a2, t)

print('\npercentile\tdefault\timproved')
p = np.percentile(a, 99,axis=0)
p2 = np.percentile(a2, 99,axis=0)
print('99\t\t', p, '\t\t', p2)
p = np.percentile(a, 99.9,axis=0)
p2 = np.percentile(a2, 99.9,axis=0)
print('99.9\t',p, '\t', p2)
p = np.percentile(a, 99.99,axis=0)
p2 = np.percentile(a2, 99.99,axis=0)
print('99.99\t',p, '\t', p2)

ecdf = sm.distributions.ECDF(a)
x1 = np.linspace(min(a), max(a), 300000)
y1 = ecdf(x1)

ecdf = sm.distributions.ECDF(a2)
x2 = np.linspace(min(a2), max(a2), 300000)
y2 = ecdf(x2)

#plt.axis([200, 5000, 0.99, 1.001])
p = plt.plot(x1, y1, 'r', x2, y2, 'g')
plt.xlabel('uSeconds')
plt.ylabel('Percentile')
plt.setp(p, linewidth=2.0)
plt.grid(True, which='both')

plt.show()
'''
