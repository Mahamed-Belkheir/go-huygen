import sklearn.svm as svm
import sklearn.preprocessing as prep
import pickle
import util

model: svm.OneClassSVM = pickle.load(open("./probe.model", "r"))
scaler: prep.RobustScaler = pickle.load(open("./probe.model", "r"))

