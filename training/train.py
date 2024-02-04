import numpy as np
import sklearn.svm as svm
import sklearn.preprocessing as prep
import sklearn.datasets as datasets

f1 = np.genfromtxt("../data/127.0.0.1_8081-1706975167.csv", delimiter=",")
f2 = np.genfromtxt("../data/127.0.0.1_8082-1706975167.csv", delimiter=",")
f3 = np.genfromtxt("../data/127.0.0.1_8083-1706975167.csv", delimiter=",")

all = np.concatenate([f1, f2])

def biggestDifference(arr):
    biggestDiffs = []
    for row in arr:
        maxDiff = 0
        row = row[:5]
        for i in range(0, len(row)-4, 4):
            diff = abs(row[i+1] - row[i+3])
            if diff > maxDiff:
                maxDiff = diff
        biggestDiffs.append(maxDiff)
    return biggestDiffs

all = np.delete(all, [0,1,2], 1)
all = np.c_[all, biggestDifference(all)]

scaler = prep.RobustScaler().fit(all)

all = scaler.transform(all)


f3 = np.delete(f3, [0,1,2], 1)
f3 = np.c_[f3, biggestDifference(f3)]
f3 = scaler.transform(f3)

m = svm.OneClassSVM(gamma="scale", kernel="rbf").fit(all[:,[0, -1]])

import pickle

pickle.dump(m, open("probe.model", 'wb'))
pickle.dump(scaler, open("probe.scaler", 'wb'))